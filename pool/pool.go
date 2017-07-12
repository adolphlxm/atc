package pool

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var nowFunc = time.Now // for testing

// ErrPoolExhausted is returned from a pool connection method (Do, Send,
// Receive, Flush, Err) when the maximum number of database connections in the
// pool has been reached.
var ErrPoolExhausted = errors.New("atc: connection pool exhausted")

type Pool struct {

	// Dial is an application supplied function for creating and configuring a
	// connection.
	//
	// The connection returned from Dial must not be in a special state
	// (subscribed to pubsub channel, transaction started, ...).
	Dial func() (interface{}, error)

	// Close is an application supplied functoin for closeing connections.
	Close func(c interface{}) error

	// TestOnBorrow is an optional application supplied function for checking
	// the health of an idle connection before the connection is used again by
	// the application. Argument t is the time that the connection was returned
	// to the pool. If the function returns an error, then the connection is
	// closed.
	TestOnBorrow func(c interface{}, t time.Time) error

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool

	// mu protects fields defined below.
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
	active int

	// Stack of idleConn with most recently used at the front.
	idle list.List
}

type idleConn struct {
	c interface{}
	t time.Time
}

// New creates a new pool. This function is deprecated. Applications should
// initialize the Pool fields directly as shown in example.
func New(dialFn func() (interface{}, error), closeFn func(c interface{}) error, maxIdle int) *Pool {
	return &Pool{Dial: dialFn, Close: closeFn, MaxIdle: maxIdle}
}

// Get gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection. If there is an error
// getting an underlying connection, then the connection Err,
func (p *Pool) Get() (interface{}, error) {
	return p.get()
}

// Put adds conn back to the pool, use forceClose to close the connection forcely
func (p *Pool) Put(c interface{}, forceClose bool) error {
	p.mu.Lock()

	if !forceClose {
		if !p.closed {
			p.idle.PushFront(idleConn{t: nowFunc(), c: c})
			if p.idle.Len() > p.MaxIdle {
				// remove exceed conn
				c = p.idle.Remove(p.idle.Back()).(idleConn).c
			} else {
				c = nil
			}
		}
	}
	// close exceed conn
	if c != nil {
		p.release()
		p.mu.Unlock()
		return p.Close(c) // 关闭该连接
	}

	if p.cond != nil {
		p.cond.Signal() // 成功放回空闲连接通知其他阻塞的进程
	}
	p.mu.Unlock()
	return nil
}

// get prunes stale connections and returns a connection from the idle list or
// creates a new connection.
func (p *Pool) get() (interface{}, error) {
	p.mu.Lock()

	// Prune stale connections.
	// 最大空闲连接等待时间, 超过此时间后关闭连接
	if timeout := p.IdleTimeout; timeout > 0 {
		for i, n := 0, p.idle.Len(); i < n; i++ {
			// 返回链表最后一个元素(空闲时间最长)
			e := p.idle.Back()
			if e == nil {
				break
			}

			// 获取连接内容
			ic := e.Value.(idleConn)

			// ic.t + timeout
			// 如果空闲时间最长的连接都没有超时，则不修剪
			if ic.t.Add(timeout).After(nowFunc()) {
				break
			}

			// 从空闲连接链表中删除该元素
			p.idle.Remove(e)
			// 减少p.active, 发消息给阻塞的请求
			p.release()
			p.mu.Unlock()

			// 关闭该连接
			p.Close(ic.c)

			p.mu.Lock()
		}
	}

	for {

		// Get idle connection.
		for i, n := 0, p.idle.Len(); i < n; i++ {
			// 返回链表第一个元素，刚刚使用过的连接
			e := p.idle.Front()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			p.idle.Remove(e)
			test := p.TestOnBorrow
			p.mu.Unlock()

			// 空闲连接再次使用
			if test == nil || test(ic.c, ic.t) == nil {
				return ic.c, nil
			}
			// 如果返回错误则关闭连接
			// ic.c.Close()
			p.Close(ic.c)
			p.mu.Lock()
			p.release()
		}

		// 检查pool本身有没有关闭
		if p.closed {
			p.mu.Unlock()
			return nil, errors.New("atc: get on closed pool")
		}

		// Dial new connection if under limit.

		if p.MaxActive == 0 || p.active < p.MaxActive {
			dial := p.Dial
			p.active += 1
			p.mu.Unlock()
			c, err := dial()
			if err != nil {
				p.mu.Lock()
				p.release()
				p.mu.Unlock()
				c = nil
			}
			return c, err
		}

		if !p.Wait {
			p.mu.Unlock()
			return nil, ErrPoolExhausted
		}

		if p.cond == nil {
			p.cond = sync.NewCond(&p.mu)
		}

		// 等待通知,暂时阻塞
		p.cond.Wait()
	}
}

// release decrements the active count and signals waiters. The caller must
// hold p.mu during the call.
func (p *Pool) release() {
	p.active -= 1
	if p.cond != nil {
		// 下发一个通知给已获取锁的goroutine
		p.cond.Signal()
	}
}
