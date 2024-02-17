package model_test

import (
	"encoding/hex"
	"testing"

	"github.com/landru29/adsb1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func airborneVelocity(t *testing.T, input string) model.AirborneVelocity {
	t.Helper()

	dataByte, err := hex.DecodeString(input)
	require.NoError(t, err)

	require.NoError(t, model.ModeS(dataByte).CheckSum())

	squitter, err := model.ModeS(dataByte).QualifiedMessage()
	require.NoError(t, err)

	require.Equal(t, "extended squitter", squitter.Name())

	extendedSquitter, ok := squitter.(model.ExtendedSquitter)
	assert.True(t, ok)

	msg, err := extendedSquitter.Decode()
	require.NoError(t, err)

	assert.Equal(t, "airborne velocity", msg.Name())

	velocity, ok := msg.(model.AirborneVelocity)
	assert.True(t, ok)

	return velocity
}

func TestAirborneVelocity(t *testing.T) {
	t.Parallel()

	t.Run("basics", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D34620499083E383008054D8CB4")
		require.NoError(t, err)

		require.NoError(t, model.ModeS(dataByte).CheckSum())

		squitter, err := model.ModeS(dataByte).QualifiedMessage()
		require.NoError(t, err)

		require.Equal(t, "extended squitter", squitter.Name())

		extendedSquitter, ok := squitter.(model.ExtendedSquitter)
		assert.True(t, ok)

		msg, err := extendedSquitter.Decode()
		require.NoError(t, err)

		assert.Equal(t, "airborne velocity", msg.Name())

		_, ok = msg.(model.AirborneVelocity)
		assert.True(t, ok)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		messageA := airborneVelocity(t, "8D485020994409940838175B284F")
		messageB := airborneVelocity(t, "8DA05F219B06B6AF189400CBC33F")

		assert.EqualValues(t, -832, messageA.VerticalRate())
		assert.EqualValues(t, -2304, messageB.VerticalRate())

		assert.EqualValues(t, 550, messageA.DeltaBarometric())
		assert.EqualValues(t, 0, messageB.DeltaBarometric())

		// Message A
		assert.True(t, messageA.IsGroundSpeed())
		assert.False(t, messageA.IsTrueAirSpeed())

		speedA, headingA := messageA.Speed()
		assert.InDelta(t, 159.20, speedA, 0.01)
		assert.InDelta(t, 182.88, headingA, 0.01)

		// Message B
		assert.False(t, messageB.IsGroundSpeed())
		assert.True(t, messageB.IsTrueAirSpeed())

		speedB, headingB := messageB.Speed()
		assert.InDelta(t, 375, speedB, 0.01)
		assert.InDelta(t, 243.98, headingB, 0.01)
	})
}
