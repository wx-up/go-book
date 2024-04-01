package logger

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessLog struct {
	Method   string `json:"method"`
	Url      string `json:"url"`
	Duration string `json:"duration"`
	ReqBody  string `json:"req_body"`
	RespBody string `json:"resp_body"`
	Status   int    `json:"status"`
}
type Builder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx *gin.Context, al *AccessLog)
}

func NewBuilder(loggerFunc func(ctx *gin.Context, al *AccessLog)) *Builder {
	return &Builder{
		allowReqBody:  false,
		allowRespBody: false,
		loggerFunc:    loggerFunc,
	}
}

func (b *Builder) AllowReqBody() *Builder {
	b.allowReqBody = true
	return b
}

func (b *Builder) AllowRespBody() *Builder {
	b.allowRespBody = true
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()

		// 限制 url 的最大长度
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}

		// 限制请求体的最大长度
		if b.allowReqBody && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			reader := io.NopCloser(bytes.NewReader(body))
			ctx.Request.Body = reader
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}

		if b.allowRespBody {
			ctx.Writer = &responseWriter{al, ctx.Writer}
		}

		defer func() {
			al.Duration = time.Now().Sub(start).String()
			b.loggerFunc(ctx, al)
		}()

		ctx.Next()
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if len(b) > 1024 {
		w.al.RespBody = string(b[:1024])
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteHeader(code int) {
	w.al.Status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	if len(s) > 1024 {
		w.al.RespBody = s[:1024]
	}
	return w.ResponseWriter.WriteString(s)
}
