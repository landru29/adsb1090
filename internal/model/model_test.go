package model_test

import (
	"testing"

	"github.com/landru29/adsb1090/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGrayToBinary(t *testing.T) {
	t.Parallel()

	assert.EqualValues(t, 0, model.GrayToBinary(0b00000))
	assert.EqualValues(t, 1, model.GrayToBinary(0b00001))
	assert.EqualValues(t, 2, model.GrayToBinary(0b00011))
	assert.EqualValues(t, 3, model.GrayToBinary(0b00010))
	assert.EqualValues(t, 4, model.GrayToBinary(0b00110))
	assert.EqualValues(t, 5, model.GrayToBinary(0b00111))
	assert.EqualValues(t, 6, model.GrayToBinary(0b00101))
	assert.EqualValues(t, 7, model.GrayToBinary(0b00100))
}

// func TestGillhamToBinary(t *testing.T) {
// 	assert.EqualValues(t, 0, model.GillhamToBinary(0b00000))
// 	assert.EqualValues(t, 1, model.GillhamToBinary(0b00001))
// 	assert.EqualValues(t, 2, model.GillhamToBinary(0b00011))
// 	assert.EqualValues(t, 3, model.GillhamToBinary(0b00010))
// 	assert.EqualValues(t, 4, model.GillhamToBinary(0b00110))
// 	assert.EqualValues(t, 7, model.GillhamToBinary(0b00100))
// 	assert.EqualValues(t, 8, model.GillhamToBinary(0b01100))
// 	assert.EqualValues(t, 11, model.GillhamToBinary(0b01110))
// 	assert.EqualValues(t, 12, model.GillhamToBinary(0b01010))
// 	assert.EqualValues(t, 13, model.GillhamToBinary(0b01011))
// 	assert.EqualValues(t, 14, model.GillhamToBinary(0b01001))
// 	assert.EqualValues(t, 17, model.GillhamToBinary(0b11001))
// 	assert.EqualValues(t, 18, model.GillhamToBinary(0b11011))
// 	assert.EqualValues(t, 19, model.GillhamToBinary(0b11010))
// }
