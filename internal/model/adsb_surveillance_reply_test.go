package model_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/adsb1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSurveillanceReplyWithAltitude(t *testing.T) {
	t.Parallel()

	t.Run("altitude", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("2000171806A983")
		require.NoError(t, err)

		qualifiedMessage, err := model.ModeS(dataByte).QualifiedMessage()
		require.NoError(t, err)

		require.Equal(t, "short message", qualifiedMessage.Name())

		shortMessage, ok := qualifiedMessage.(model.ShortMessage)
		assert.True(t, ok)

		surveillanceReplyWithAltitude := model.SurveillanceReplyWithAltitude{shortMessage}

		assert.EqualValues(t, 36000.0, surveillanceReplyWithAltitude.Altitude())
	})
}

func TestSurveillanceReplyWithIdentification(t *testing.T) {
	t.Parallel()

	t.Run("identity", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("2A00516D492B80")
		require.NoError(t, err)

		qualifiedMessage, err := model.ModeS(dataByte).QualifiedMessage()
		require.NoError(t, err)

		require.Equal(t, "short message", qualifiedMessage.Name())

		shortMessage, ok := qualifiedMessage.(model.ShortMessage)
		assert.True(t, ok)

		surveillanceReplyWithIdentification := model.SurveillanceReplyWithIdentification{shortMessage}

		assert.EqualValues(t, 356, surveillanceReplyWithIdentification.Identity())
	})
}
