package trace

import (
	"fmt"
	"io"
)

// インタフェースを定義
type Tracer interface {
	Trace(...interface{})
}

// tracerを生成
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprintf(t.out, fmt.Sprintf("%v\n", a...))
}

// 何もしないtracer
type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

func Off(a ...interface{}) Tracer {
	return &nilTracer{}
}
