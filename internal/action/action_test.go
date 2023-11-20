package action

import (
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestByteSplitInt16(t *testing.T) {
	testCases := []struct {
		input         int
		expectedByte1 byte
		expectedByte2 byte
	}{
		{0, 0x00, 0x00},       // Zero case
		{1, 0x00, 0x01},       // Smallest positive integer
		{-1, 0xFF, 0xFF},      // Smallest negative integer (all bits set)
		{2, 0x00, 0x02},       // Small positive power of two
		{-2, 0xFF, 0xFE},      // Small negative power of two
		{32767, 0x7F, 0xFF},   // Largest positive integer for 16-bit
		{-32768, 0x80, 0x00},  // Largest negative integer for 16-bit (boundary)
		{128, 0x00, 0x80},     // Affects only the lower byte
		{-129, 0xFF, 0x7F},    // Affects both bytes, negative value
		{256, 0x01, 0x00},     // Affects only the upper byte
		{-32767, 0x80, 0x01},  // One more than the largest negative integer
		{100, 0x00, 0x64},     // Arbitrary positive number
		{-100, 0xFF, 0x9C},    // Arbitrary negative number
		{0x1234, 0x12, 0x34},  // Hexadecimal representation in the middle of the range
		{-0x1234, 0xED, 0xCC}, // Negative hexadecimal in the middle of the range
	}

	for _, tc := range testCases {
		byte1, byte2 := ByteSplitInt16(tc.input)
		if byte1 != tc.expectedByte1 || byte2 != tc.expectedByte2 {
			t.Errorf("byteSplitInt(%d) = %08b, %08b; want %08b, %08b",
				tc.input, byte1, byte2, tc.expectedByte1, tc.expectedByte2)
		}
	}
}

// TestFloat64To4Bytes tests the float64To4Bytes function with various test cases.
func TestFloat64To4Bytes(t *testing.T) {
	testCases := []struct {
		input    float64
		expected []byte
	}{
		{0.0, []byte{0x00, 0x00, 0x00, 0x00}},                         // Zero
		{1.0, []byte{0x3F, 0x80, 0x00, 0x00}},                         // Positive 1
		{-1.0, []byte{0xBF, 0x80, 0x00, 0x00}},                        // Negative 1
		{math.MaxFloat32, []byte{0x7F, 0x7F, 0xFF, 0xFF}},             // Max float32
		{math.SmallestNonzeroFloat32, []byte{0x00, 0x00, 0x00, 0x01}}, // Smallest nonzero float32
		{math.Pi, []byte{0x40, 0x49, 0x0F, 0xDB}},                     // Pi
		{math.Inf(1), []byte{0x7F, 0x80, 0x00, 0x00}},                 // Positive infinity
		{math.Inf(-1), []byte{0xFF, 0x80, 0x00, 0x00}},                // Negative infinity
		{math.NaN(), []byte{0x7F, 0xC0, 0x00, 0x00}},                  // NaN
		// ... additional test cases as needed
	}

	for _, tc := range testCases {
		bytes := make([]byte, 4)
		float64To4Bytes(tc.input, bytes)
		for i, v := range tc.expected {
			if bytes[i] != v {
				t.Errorf("float64To4Bytes(%f) = %02X; want %02X at index %d", tc.input, bytes[i], v, i)
			}
		}
	}
}

// TestAppendPositionData tests whether appendPositionData correctly appends int16 and float32 values.
func TestAppendPositionData(t *testing.T) {
	testCases := []struct {
		pos               *mat.VecDense
		startPos          int
		expectedBytes     []byte
		expectedNextIndex int
	}{
		{
			mat.NewVecDense(3, []float64{300, 400, 5.5}),
			2,
			[]byte{0, 0, 0x01, 0x2C, 0x01, 0x90, 0x40, 0xB0, 0x00, 0x00, 0, 0},
			10,
		},
		{
			mat.NewVecDense(3, []float64{32767, -32768, 1.0}),
			0,
			[]byte{0x7F, 0xFF, 0x80, 0x00, 0x3F, 0x80, 0x00, 0x00},
			8,
		},
		{
			mat.NewVecDense(3, []float64{12345, -12345, math.Pi}),
			3,
			[]byte{0, 0, 0, 0x30, 0x39, 0xCF, 0xC7, 0x40, 0x49, 0x0F, 0xDB},
			11,
		},
	}

	for _, tc := range testCases {
		// Create a message slice with some initial data
		message := make([]byte, len(tc.expectedBytes))
		copy(message, tc.expectedBytes) // Copy initial data if needed

		// Call appendPositionData
		nextIndex := appendPositionData(message, tc.startPos, tc.pos)

		// Check if the data in message is as expected
		for i, v := range tc.expectedBytes {
			if message[i] != v {
				t.Errorf("Test case failed for VecDense %+v: expected byte at index %d to be %02X, got %02X", tc.pos, i, v, message[i])
			}
		}

		// Check if the next index is as expected
		if nextIndex != tc.expectedNextIndex {
			t.Errorf("Test case failed for VecDense %+v: expected next index to be %d, got %d", tc.pos, tc.expectedNextIndex, nextIndex)
		}
	}
}

func TestMoveTranslateReal(t *testing.T) {
	testCases := []struct {
		move          Move
		expectedBytes []byte
	}{
		{
			Move{
				Id:   1,
				Pos:  mat.NewVecDense(3, []float64{100, 200, math.Pi}),
				Goal: mat.NewVecDense(3, []float64{300, 400, -math.Pi}),
			},
			[]byte{19, byte(RealMove), 1,
				0x00, 0x64, 0x00, 0xC8, 0x40, 0x49, 0x0F, 0xDB,
				0x01, 0x2C, 0x01, 0x90, 0xC0, 0x49, 0x0F, 0xDB,
			},
		},
		{
			Move{
				Id:   2,
				Pos:  mat.NewVecDense(3, []float64{150, 250, 0}),
				Goal: mat.NewVecDense(3, []float64{350, 450, math.Pi / 2}),
			},
			[]byte{19, byte(RealMove), 2,
				0x00, 0x96, 0x00, 0xFA, 0x00, 0x00, 0x00, 0x00,
				0x01, 0x5E, 0x01, 0xC2, 0x3F, 0xC9, 0x0F, 0xDB,
			},
		},
	}

	for _, tc := range testCases {
		gotBytes := tc.move.TranslateReal()
		if !reflect.DeepEqual(gotBytes, tc.expectedBytes) {
			posValues := make([]float64, tc.move.Pos.Len())
			goalValues := make([]float64, tc.move.Goal.Len())
			for i := 0; i < tc.move.Pos.Len(); i++ {
				posValues[i] = tc.move.Pos.AtVec(i)
			}
			for i := 0; i < tc.move.Goal.Len(); i++ {
				goalValues[i] = tc.move.Goal.AtVec(i)
			}
			t.Errorf("Move.TranslateReal() for Move {id: %d, pos: %v, goal: %v} = %v, want %v", tc.move.Id, posValues, goalValues, gotBytes, tc.expectedBytes)
		}
	}
}

func TestStopTranslateReal(t *testing.T) {
	testCases := []struct {
		stop          Stop
		expectedBytes []byte
	}{
		{
			Stop{Id: 1},
			[]byte{3, byte(RealStop), 1},
		},
		{
			Stop{Id: 2},
			[]byte{3, byte(RealStop), 2},
		},
		{
			Stop{Id: 255}, // Testing the edge case with the maximum byte value
			[]byte{3, byte(RealStop), 255},
		},
	}

	for _, tc := range testCases {
		gotBytes := tc.stop.TranslateReal()
		if !reflect.DeepEqual(gotBytes, tc.expectedBytes) {
			t.Errorf("Stop.TranslateReal() for Stop {id: %d} = %v, want %v", tc.stop.Id, gotBytes, tc.expectedBytes)
		}
	}
}

func TestKickTranslateReal(t *testing.T) {
	testCases := []struct {
		kick          Kick
		expectedBytes []byte
	}{
		{
			Kick{Id: 1, Speed: 255},
			[]byte{4, byte(RealKick), 1, 255},
		},
		{
			Kick{Id: 2, Speed: 0},
			[]byte{4, byte(RealKick), 2, 0},
		},
		{
			Kick{Id: 255, Speed: 2}, // Testing the edge case with the maximum byte value
			[]byte{4, byte(RealKick), 255, 2},
		},
	}

	for _, tc := range testCases {
		gotBytes := tc.kick.TranslateReal()
		if !reflect.DeepEqual(gotBytes, tc.expectedBytes) {
			t.Errorf("Kick.TranslateReal() for Kick {id: %d} = %v, want %v", tc.kick.Id, gotBytes, tc.expectedBytes)
		}
	}
}

func TestInitTranslateReal(t *testing.T) {
	testCases := []struct {
		init          Init
		expectedBytes []byte
	}{
		{
			Init{Id: 1},
			[]byte{3, byte(RealInit), 1},
		},
		{
			Init{Id: 2},
			[]byte{3, byte(RealInit), 2},
		},
		{
			Init{Id: 255}, // Testing the edge case with the maximum byte value
			[]byte{3, byte(RealInit), 255},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		gotBytes := tc.init.TranslateReal()
		if !reflect.DeepEqual(gotBytes, tc.expectedBytes) {
			t.Errorf("Init.TranslateReal() for Init {id: %d} = %v, want %v", tc.init.Id, gotBytes, tc.expectedBytes)
		}
	}
}
