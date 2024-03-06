package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustInit(t *testing.T) {
	tests := []string{"debug", "info", "warn", "error"}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			assert.NotPanics(t, func() { MustInit(tt) })
		})
	}
}
