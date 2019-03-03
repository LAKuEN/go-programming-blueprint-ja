package trace

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("retval is nil")
	}
	msg := "Hello, trace package"
	want := fmt.Sprintf("%v\n", msg)

	tracer.Trace(msg)
	if buf.String() != want {
		t.Errorf("Outputted wrong mesage: %v", buf.String())
	}
}

func TestOff(t *testing.T) {
	silentTracer := Off()
	silentTracer.Trace("Data")
}
