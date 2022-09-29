package contract

import (
	"context"
	"net/http"
)

const TraceKey = "web:trace"

const (
	TraceKeyTraceID  = "trace_id"
	TraceKeySpanID   = "span_id"
	TraceKeyCspanID  = "cspan_id"
	TraceKeyParentID = "parent_id"
	TraceKeyMethod   = "method"
	TraceKeyCaller   = "caller"
	TraceKeyTime     = "time"
)

type Trace interface {
	// WithTrace register new trace to context
	WithTrace(context.Context, *TraceContext) context.Context
	// GetTrace From trace context
	GetTrace(context.Context) *TraceContext
	// NewTrace generate a new trace
	NewTrace() *TraceContext
	// StartSpan generate cspan for child call
	StartSpan(*TraceContext) *TraceContext
	// ToMap 将logger转换为map
	ToMap(*TraceContext) map[string]string
	// 通过http
	ExtractHTTP(*http.Request) *TraceContext
	// 给http添加trace
	InjectHTTP(*http.Request, *TraceContext) *http.Request
}

type TraceContext struct {
	TraceID    string // traceID global unique
	ParentID   string // 父节点SpanID
	SpanID     string // 当前节点SpanID
	CspanId    string // 子节点调用的SpanID，由调用方指定
	Annotation map[string]string
}
