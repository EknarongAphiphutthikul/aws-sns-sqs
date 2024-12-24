package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type customStruct struct {
	Str         string
	Flag        bool
	StrPointer  *string
	FlagPointer *bool
	Int         int
	IntPointer  *int
}

func TestIsZeroValue(t *testing.T) {
	var valueInt8 int8
	var valueInt16 int16
	var valueInt32 int32
	var valueInt64 int64
	var valueUInt8 uint8
	var valueUInt16 uint16
	var valueUInt32 uint32
	var valueUInt64 uint64
	var valueString string
	var valueBool bool
	var valueFloat32 float32
	var valueFloat64 float64
	var valueComplex64 complex64
	var valueComplex128 complex128
	var valueArray [1]int
	var valueSlice []int
	var valueMap map[int]int
	var valueStruct struct{}
	var valueInterface interface{}
	var valueFunc func()
	var valueCustomStruct customStruct

	assert.True(t, IsZeroValue(valueInt8))
	assert.True(t, IsZeroValue(valueInt16))
	assert.True(t, IsZeroValue(valueInt32))
	assert.True(t, IsZeroValue(valueInt64))
	assert.True(t, IsZeroValue(valueUInt8))
	assert.True(t, IsZeroValue(valueUInt16))
	assert.True(t, IsZeroValue(valueUInt32))
	assert.True(t, IsZeroValue(valueUInt64))
	assert.True(t, IsZeroValue(valueString))
	assert.True(t, IsZeroValue(valueBool))
	assert.True(t, IsZeroValue(valueFloat32))
	assert.True(t, IsZeroValue(valueFloat64))
	assert.True(t, IsZeroValue(valueComplex64))
	assert.True(t, IsZeroValue(valueComplex128))
	assert.True(t, IsZeroValue(valueArray))
	assert.True(t, IsZeroValue(valueSlice))
	assert.True(t, IsZeroValue(valueMap))
	assert.True(t, IsZeroValue(valueStruct))
	assert.True(t, IsZeroValue(valueInterface))
	assert.True(t, IsZeroValue(valueFunc))
	assert.True(t, IsZeroValue(valueCustomStruct))
	assert.True(t, IsZeroValue(nil))
	assert.True(t, IsZeroValue(&valueInt8))
	assert.True(t, IsZeroValue(&valueInt16))
	assert.True(t, IsZeroValue(&valueInt32))
	assert.True(t, IsZeroValue(&valueInt64))
	assert.True(t, IsZeroValue(&valueUInt8))
	assert.True(t, IsZeroValue(&valueUInt16))
	assert.True(t, IsZeroValue(&valueUInt32))
	assert.True(t, IsZeroValue(&valueUInt64))
	assert.True(t, IsZeroValue(&valueString))
	assert.True(t, IsZeroValue(&valueBool))
	assert.True(t, IsZeroValue(&valueFloat32))
	assert.True(t, IsZeroValue(&valueFloat64))
	assert.True(t, IsZeroValue(&valueComplex64))
	assert.True(t, IsZeroValue(&valueComplex128))
	assert.True(t, IsZeroValue(&valueArray))
	assert.True(t, IsZeroValue(&valueSlice))
	assert.True(t, IsZeroValue(&valueMap))
	assert.True(t, IsZeroValue(&valueStruct))
	assert.True(t, IsZeroValue(&valueInterface))
	assert.True(t, IsZeroValue(&valueFunc))
	assert.True(t, IsZeroValue(&valueCustomStruct))

	valueInt8 = 1
	valueInt16 = 1
	valueInt32 = 1
	valueInt64 = 1
	valueUInt8 = 1
	valueUInt16 = 1
	valueUInt32 = 1
	valueUInt64 = 1
	valueString = "1"
	valueBool = true
	valueFloat32 = 1
	valueFloat64 = 1
	valueComplex64 = 1
	valueComplex128 = 1
	valueArray = [1]int{1}
	valueSlice = []int{1}
	valueMap = map[int]int{1: 1}
	valueInterface = 1
	valueFunc = func() {}
	tempStr := "2"
	tempInt := 2
	tempFlag := true
	valueCustomStruct = customStruct{Str: "1", Flag: true, Int: 10, StrPointer: &tempStr, FlagPointer: &tempFlag, IntPointer: &tempInt}

	assert.False(t, IsZeroValue(valueInt8))
	assert.False(t, IsZeroValue(valueInt16))
	assert.False(t, IsZeroValue(valueInt32))
	assert.False(t, IsZeroValue(valueInt64))
	assert.False(t, IsZeroValue(valueUInt8))
	assert.False(t, IsZeroValue(valueUInt16))
	assert.False(t, IsZeroValue(valueUInt32))
	assert.False(t, IsZeroValue(valueUInt64))
	assert.False(t, IsZeroValue(valueString))
	assert.False(t, IsZeroValue(valueBool))
	assert.False(t, IsZeroValue(valueFloat32))
	assert.False(t, IsZeroValue(valueFloat64))
	assert.False(t, IsZeroValue(valueComplex64))
	assert.False(t, IsZeroValue(valueComplex128))
	assert.False(t, IsZeroValue(valueArray))
	assert.False(t, IsZeroValue(valueSlice))
	assert.False(t, IsZeroValue(valueMap))
	assert.False(t, IsZeroValue(valueInterface))
	assert.False(t, IsZeroValue(valueFunc))
	assert.False(t, IsZeroValue(valueCustomStruct))
	assert.False(t, IsZeroValue(&valueInt8))
	assert.False(t, IsZeroValue(&valueInt16))
	assert.False(t, IsZeroValue(&valueInt32))
	assert.False(t, IsZeroValue(&valueInt64))
	assert.False(t, IsZeroValue(&valueUInt8))
	assert.False(t, IsZeroValue(&valueUInt16))
	assert.False(t, IsZeroValue(&valueUInt32))
	assert.False(t, IsZeroValue(&valueUInt64))
	assert.False(t, IsZeroValue(&valueString))
	assert.False(t, IsZeroValue(&valueBool))
	assert.False(t, IsZeroValue(&valueFloat32))
	assert.False(t, IsZeroValue(&valueFloat64))
	assert.False(t, IsZeroValue(&valueComplex64))
	assert.False(t, IsZeroValue(&valueComplex128))
	assert.False(t, IsZeroValue(&valueArray))
	assert.False(t, IsZeroValue(&valueSlice))
	assert.False(t, IsZeroValue(&valueMap))
	assert.False(t, IsZeroValue(&valueInterface))
	assert.False(t, IsZeroValue(&valueFunc))
	assert.False(t, IsZeroValue(&valueCustomStruct))
}
