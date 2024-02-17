package model_test

import (
	"bufio"
	"encoding/hex"
	"os"
	"testing"

	"github.com/landru29/adsb1090/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecksum(t *testing.T) {
	t.Parallel()

	t.Run("ok extended squitter", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D40621D58C382D690C8AC2863A7")
		require.NoError(t, err)

		require.NoError(t, model.ModeS(dataByte).CheckSum())
	})

	t.Run("ko extended squitter", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("8D40621D59C382D690C8AC2863A7")
		require.NoError(t, err)

		require.Error(t, model.ModeS(dataByte).CheckSum())
	})

	t.Run("ok short squitter", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("5d4ca92bf0802f")
		require.NoError(t, err)

		require.NoError(t, model.ModeS(dataByte).CheckSum())
	})

	t.Run("ko short squitter", func(t *testing.T) {
		t.Parallel()
		dataByte, err := hex.DecodeString("5d4ac9c46451eb")
		require.NoError(t, err)

		require.Error(t, model.ModeS(dataByte).CheckSum())
	})
}

func TestSquitter(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/dump1090.txt")
	require.NoError(t, err)

	scanner := bufio.NewScanner(file)

	lineIdx := 1

	for scanner.Scan() {
		line := scanner.Text()

		message := line[1 : len(line)-1]

		dataByte, err := hex.DecodeString(message)
		require.NoError(t, err, "line #%d: %s", lineIdx, message)

		squitter, err := model.ModeS(dataByte).QualifiedMessage()
		require.NoError(t, err, "line #%d: %s", lineIdx, message)

		switch {
		case len(dataByte) == 7:
			assert.Equal(t, "short message", squitter.Name(), "line #%d: %s", lineIdx, message)
		case len(dataByte) == 14:
			assert.Contains(t, []string{"extended squitter", "long message"}, squitter.Name(), "line #%d: %s", lineIdx, message)
			// assert.Equal(t, "extended squitter", squitter.Name(), "line #%d: %s", lineIdx, message)
		default:
			t.Fatalf("wrong size of message, line #%d: %s", lineIdx, message)
		}

		require.NoError(t, model.ModeS(dataByte).CheckSum(), "line #%d: %s", lineIdx, message)

		lineIdx++
	}
}

func TestIcaoAddrChecksum(t *testing.T) {
	t.Parallel()

	t.Run("Short Air-Air Surveillance (0)", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("02e61838fb04f6")
		require.NoError(t, err)

		assert.Equal(t, model.ICAOAddr(0x346204), model.ModeS(dataByte).IcaoAddrChecksum())
	})

	t.Run("Altitude Reply (4)", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("2000191052962c")
		require.NoError(t, err)

		assert.Equal(t, model.ICAOAddr(0x4ca92b), model.ModeS(dataByte).IcaoAddrChecksum())
	})

	t.Run("Identity Reply (5)", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("28000426550278")
		require.NoError(t, err)

		assert.Equal(t, model.ICAOAddr(0x4ca92b), model.ModeS(dataByte).IcaoAddrChecksum())
	})

	t.Run("CommB With Altitude Reply (20)", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("A0001839CA380030AA0000C8B28A")
		require.NoError(t, err)

		assert.Equal(t, model.ICAOAddr(0x346204), model.ModeS(dataByte).IcaoAddrChecksum())
	})

	t.Run("CommB With Identity Reply (21)", func(t *testing.T) {
		t.Parallel()

		dataByte, err := hex.DecodeString("A8000CAC80105938BFF4DC616BB7")
		require.NoError(t, err)

		assert.Equal(t, model.ICAOAddr(0x346204), model.ModeS(dataByte).IcaoAddrChecksum())
	})
}
