// Package input is the data source.
package input

import (
	"context"

	"github.com/landru29/adsb1090/internal/processor"
)

// Starter is a process starter.
type Starter interface {
	Start(ctx context.Context, processors ...processor.Processer) error
}
