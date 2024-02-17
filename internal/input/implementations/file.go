package implementations

import (
	"context"
	"io"
	"os"
	"syscall"

	"github.com/landru29/adsb1090/internal/processor"
)

// File is the data file source.
type File struct {
	filename string
	loop     bool
}

// FileConfigurator is the Source configurator.
type FileConfigurator func(*File)

// NewFile creates a new data source process.
func NewFile(filename string, opts ...FileConfigurator) *File {
	output := &File{
		filename: filename,
	}

	for _, opt := range opts {
		opt(output)
	}

	return output
}

// WithLoop is the data loop configurator.
func WithLoop() FileConfigurator {
	return func(s *File) {
		s.loop = true
	}
}

// Start implements the input.Starter interface.
func (s *File) Start(ctx context.Context, processors ...processor.Processer) error {
	if s.loop {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				if err := s.start(ctx, processors...); err != nil {
					return err
				}
			}
		}
	}

	err := s.start(ctx, processors...)

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	return err
}

func (s *File) start(ctx context.Context, processors ...processor.Processer) error {
	fileDescriptor, err := os.Open(s.filename)
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(fileDescriptor)

	return NewReader(fileDescriptor).Start(ctx, processors...)
}
