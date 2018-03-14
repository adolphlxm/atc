package atc

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/adolphlxm/atc/context"
	"io"
	"os"
	"strconv"
	"sync"
	"github.com/adolphlxm/atc/logs"
	"errors"
)

var (
	ErrorCode *ErrorMap

	defaultMsg  = "Undefined error"
	defaultCode = map[int]string{
		500: "Internal error, Loading handler to fail.",
		404: "Not found, Not find the handler or action.",
		406: "Not Acceptable, JSON Marshal fail.",
		407: "Not Acceptable, Websocket JSON reveice fail.",
		408: "Not Acceptable, Base64 decryption fail.",
		409: "Not Acceptable, Body argument is invalid.",
		410: "Not Acceptable, Invalid websocket json request method.",
		411: "Not Acceptable, Invalid websocket json request module.",
		400: "Unauthorized, Token expires or invalid.",
	}
)

/************************************/
/*********    Error code   **********/
/************************************/
// A ErrorMap is a Error code.
//
// Store global error codes defined
// Add the map element, there are the mutex (lock)
type ErrorMap struct {
	lock *sync.Mutex
	msg  map[int]string
}

// NewErrorMap returns a new ErrorMap
// It initialize the default error code (parseDefault())
func NewErrorMap() *ErrorMap {
	errorMap := &ErrorMap{
		lock: new(sync.Mutex),
		msg:  make(map[int]string),
	}

	errorMap.parseDefault()

	return errorMap
}

// Get returns the error code description
//
// If ErrorCode map does not exist the code, return the default msg
func (m *ErrorMap) Get(code int) string {
	if val, ok := m.msg[code]; ok {
		return val
	}
	return defaultMsg
}

// Set returns true and false
//
// If ErrorCode does not exist the code, then appended and return true.
// else return false
func (m *ErrorMap) Set(code int, msg string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.msg[code]; !ok {
		m.msg[code] = msg
	} else {
		return false
	}

	return true
}

// Delete the element with the specified key
func (m *ErrorMap) Delete(code int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.msg, code)
}

func (m *ErrorMap) parseDefault() {
	for code, msg := range defaultCode {
		m.Set(code, msg)
	}
}


func checkFileIsExist(filename string) bool {
	f, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if f != nil && f.IsDir() {
		return false
	}
	return true
}

// Parse error code file, then append to the ErrorCode map
func (m *ErrorMap) parse(ename string) error {
	if !checkFileIsExist(ename) {
		return errors.New("file does not exist.")
	}

	f, err := os.Open(ename)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	var lineNum int

	for {
		lineNum++
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if bytes.Equal(line, []byte{}) {
			continue
		} else if err != nil {
			return err
		}

		line = bytes.TrimSpace(line)
		switch {
		case bytes.HasPrefix(line, []byte{'#'}):
			continue
		case bytes.HasPrefix(line, []byte{';'}):
			continue
		default:
			optionVal := bytes.SplitN(line, []byte{'='}, 2)
			if len(optionVal) != 2 {
				return fmt.Errorf("parse %s the content error : line %d , %s = ? ", ename, lineNum, optionVal[0])
			}
			code := bytes.TrimSpace(optionVal[0])
			value := bytes.TrimSpace(optionVal[1])
			if c, err := strconv.Atoi(string(code)); err == nil {
				m.Set(c, string(value))
			}
		}
	}

	return nil
}

/************************************/
/*********  Error Response **********/
/************************************/

type Error struct {
	ctx *context.Context

	error *ErrorResponse
}

// A ResponseError is a error response content
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Request string `json:"request"`
}

func NewError(ctx *context.Context) *Error {
	return &Error{
		ctx:   ctx,
		error: &ErrorResponse{Request: ctx.Path()},
	}
}

func (res *Error) Code(status, code int) *Error {
	res.ctx.SetStatus(status)
	res.error.Code = code
	res.error.Error = ErrorCode.Get(code)
	return res
}

func (res *Error) Message(msg string) *Error {
	res.error.Error = msg
	return res
}

func (res *Error) JSON() {
	switch res.ctx.ReqType {
	case context.RPC_HTTP:
		content, _ := json.Marshal(res.error)
		res.ctx.Write(content)
	case context.RPC_WEBSOCKET:
		websocket.JSON.Send(res.ctx.WS, res.error)
	default:

	}
}


// Initialize error file.
func initError(){
	ErrorCode = NewErrorMap()
	// In the conf/error. Ini file parsing error code
	err := ErrorCode.parse(AppConfig.DefaultString("error.file", "../conf/error.ini"))
	if err != nil {
		logs.Warnf("Error file loading err:%v", err.Error())
	}
}