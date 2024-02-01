package middlewares

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"github.com/labstack/echo/v4"
)

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *bodyDumpResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func StructuredLogger(cfg *config.Config, lg logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if strings.Contains(c.Path(), "swagger") {
				next(c)
			} else {
				start := time.Now() // start
				path := c.Path()

				// Request
				reqBody := []byte{}
				if c.Request().Body != nil { // Read
					reqBody, _ = io.ReadAll(c.Request().Body)
				}
				c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

				// Response
				resBody := new(bytes.Buffer)
				mw := io.MultiWriter(c.Response().Writer, resBody)
				writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
				c.Response().Writer = writer

				if err = next(c); err != nil {
					c.Error(err)
				}

				keys := map[logging.ExtraKey]interface{}{}
				keys[logging.Path] = path
				keys[logging.ClientIp] = c.RealIP()
				keys[logging.Method] = c.Request().Method
				keys[logging.Latency] = time.Since(start)
				keys[logging.StatusCode] = c.Response().Status
				keys[logging.RequestBody] = string(reqBody)
				keys[logging.ResponseBody] = resBody

				lg.Info(logging.RequestResponse, logging.Api, "", keys)

				return
			}
			return
		}
	}
}
