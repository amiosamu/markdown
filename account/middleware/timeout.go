package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/amiosamu/markdown/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

func Timeout(timeout time.Duration, errTimeout *apperrors.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		tw := &timeoutWriter{ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			e := apperrors.NewInternalServerError()
			tw.ResponseWriter.WriteHeader(e.StatusCode())
			eResp, _ := json.Marshal(gin.H{
				"error": e,
			})
			tw.ResponseWriter.Write(eResp)
		case <-finished:
			tw.mu.Lock()
			defer tw.mu.Unlock()

			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			tw.ResponseWriter.Write(tw.wbuf.Bytes())
		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(errTimeout.StatusCode())
			eResp, _ := json.Marshal(gin.H{
				"error": errTimeout,
			})
			tw.ResponseWriter.Write(eResp)
			c.Abort()
			tw.SeTimeOut()
		}
	}
}

type timeoutWriter struct {
	gin.ResponseWriter
	h    http.Header
	wbuf bytes.Buffer

	mu          sync.Mutex
	timeOut     bool
	wroteHeader bool
	code        int
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timeOut {
		return 0, nil
	}
	return tw.wbuf.Write(b)
}

func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)

	tw.mu.Lock()

	defer tw.mu.Unlock()

	if tw.timeOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid write header code %v", code))
	}
}

func (tw *timeoutWriter) SeTimeOut() {
	tw.timeOut = true
}
