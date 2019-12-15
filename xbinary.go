package xbinary

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"sync"
)

type XBinary struct {
	l     sync.RWMutex
	ready bool
	sum   []byte
	err   error

	Header      string
	ErrorHeader string
}

func (x *XBinary) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if sum, err := x.getSum(); err != nil {
		if x.ErrorHeader != "" {
			rw.Header().Set(x.ErrorHeader, err.Error())
		}
	} else if sum != nil {
		if x.Header != "" {
			rw.Header().Set(x.Header, hex.EncodeToString(sum))
		}
	}

	next(rw, r)
}

func (x *XBinary) getSum() ([]byte, error) {
	x.l.RLock()
	if x.ready {
		x.l.RUnlock()
		return x.sum, x.err
	}
	x.l.RUnlock()

	x.l.Lock()
	if x.ready {
		x.l.Unlock()
		return x.sum, x.err
	}
	defer x.l.Unlock()

	exe, err := os.Executable()
	if err != nil {
		x.ready = true
		x.err = err
		x.sum = nil
		return x.sum, x.err
	}

	fd, err := os.Open(exe)
	if err != nil {
		x.ready = true
		x.err = err
		x.sum = nil
		return x.sum, x.err
	}
	defer fd.Close()

	exeHash := md5.New()

	if _, err := io.Copy(exeHash, fd); err != nil {
		x.ready = true
		x.err = err
		x.sum = nil
		return x.sum, x.err
	}

	x.ready = true
	x.err = nil
	x.sum = exeHash.Sum(nil)

	return x.sum, x.err
}
