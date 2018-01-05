package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sync/atomic"
	"time"
)



var (
	pid          = os.Getpid()
	machineId    = readMachineId()
	msgIdCounter = readRandomUint32()
)

func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		io.ReadFull(rand.Reader, id)
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

// readRandomUint32 returns a random msgIdCounter.
func readRandomUint32() uint32 {
	var b [4]byte
	io.ReadFull(rand.Reader, b[:])
	return uint32((uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24))
}

func GenMsgID() string {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&msgIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return string(b[:])
}

func TimestampFromMessageID(id string) (int64, error) {
	if len(id) != 12 {
		return 0, errors.New("invalid MessageID")
	}
	return int64(binary.BigEndian.Uint32([]byte(id[:4]))), nil
}
