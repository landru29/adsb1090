package implementations_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/landru29/adsb1090/internal/input/implementations"
	"github.com/landru29/adsb1090/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestReader(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockProcessor := mocks.NewMockProcesser(ctrl)

	mockProcessor.EXPECT().Process(gomock.Any()).AnyTimes()

	file, err := os.Open("../../../testdata/modes1.bin")
	require.NoError(t, err)

	defer func(closer io.Closer) {
		require.NoError(t, closer.Close())
	}(file)

	reader := implementations.NewReader(file)
	require.NotNil(t, reader)

	require.NoError(t, reader.Start(context.Background(), mockProcessor))
}
