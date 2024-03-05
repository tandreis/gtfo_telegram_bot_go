package logger

import (
	"testing"
)

func TestMustInit(t *testing.T) {
	_ = MustInit("info")
}
