package implementations

import (
	"context"
	"errors"
	"io"

	localcontext "github.com/landru29/adsb1090/internal/input/context"
	"github.com/landru29/adsb1090/internal/processor"
)

// Reader is device reader.
type Reader struct {
	reader io.Reader
}

// NewReader creates a new device reader.
func NewReader(rd io.Reader) *Reader {
	return &Reader{
		reader: rd,
	}
}

// Start implements the input.Starter interface.
func (r *Reader) Start(ctx context.Context, processors ...processor.Processer) error {
	cContext := localcontext.New(ctx, processors)

	defer func() {
		localcontext.DisposeContext(cContext.Key)
	}()

	for {
		data := make([]byte, 1024) //nolint: gomnd

		cnt, err := r.reader.Read(data)
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		processRaw(data[:cnt], cContext.Ccontext)
	}
}
