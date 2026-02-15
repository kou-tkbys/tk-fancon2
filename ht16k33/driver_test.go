// To run test `tinygo test ./ht16k33/`

package ht16k33

import (
	"bytes"
	"fmt"
	"testing"
)

// mockI2C is a mock for testing that pretends to be machine.I2C
type mockI2C struct {
	addr uint16
	data []byte
}

// Tx fakes the I2C transaction, recording the data that was supposed to be sent.
func (m *mockI2C) Tx(addr uint16, w, r []byte) error {
	m.addr = addr
	m.data = make([]byte, len(w))
	copy(m.data, w)
	return nil
}

// TestSetDigit verifies that setting a single digit correctly modifies the buffer.
func TestSetDigit(t *testing.T) {
	testCases := []struct {
		name           string
		display        int
		position       int
		num            byte
		dot            bool
		expectedBuffer [16]byte
	}{
		{
			name:     "Display 0, Position 0, Number 8, with dot",
			display:  0,
			position: 0,
			num:      8,
			dot:      true,
			// For number 8 (all segments on) at position 0, bit 0 should be set for rows 0-6.
			// For dot, bit 0 should be set for row 7.
			expectedBuffer: [16]byte{
				1 << 0, 1 << 0, 1 << 0, 1 << 0, 1 << 0, 1 << 0, 1 << 0, 1 << 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
		{
			name:     "Display 1, Position 7, Number 1, no dot",
			display:  1,
			position: 7,
			num:      1,
			dot:      false,
			// For number 1 (segments b, c) at position 7, bit 7 should be set for rows 8+1 and 8+2.
			expectedBuffer: [16]byte{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 1 << 7, 1 << 7, 0, 0, 0, 0, 0,
			},
		},
		{
			name:     "Set blank on Display 0, Position 3",
			display:  0,
			position: 3,
			num:      blankPatternIndex, // blank
			dot:      false,
			// Should result in an all-zero buffer as it clears the position.
			expectedBuffer: [16]byte{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockBus := &mockI2C{}
			device := New(mockBus, 0x70) // Creates a device with a zeroed buffer

			device.SetDigit(tc.display, tc.position, tc.num, tc.dot)

			if !bytes.Equal(device.buffer[:], tc.expectedBuffer[:]) {
				t.Errorf("FAIL: Buffer content is wrong!\nExpected: %08b\nGot:      %08b", tc.expectedBuffer, device.buffer)
			}
		})
	}
}

// TestWriteString verifies that writing a string correctly populates the buffer.
func TestWriteString(t *testing.T) {
	mockBus := &mockI2C{}
	device := New(mockBus, 0x70)

	// Write "1." to display 0 and "2" to display 1
	device.WriteString(0, "1.")
	device.WriteString(1, "2")

	expectedBuffer := [16]byte{
		0,      // D0 P0, seg a
		1 << 0, // D0 P0, seg b (from "1")
		1 << 0, // D0 P0, seg c (from "1")
		0,      // D0 P0, seg d
		0,      // D0 P0, seg e
		0,      // D0 P0, seg f
		0,      // D0 P0, seg g
		1 << 0, // D0 P0, dot (from "1.")
		1 << 0, // D1 P0, seg a (from "2")
		1 << 0, // D1 P0, seg b (from "2")
		0,      // D1, seg c
		1 << 0, // D1, seg d (from "2")
		1 << 0, // D1, seg e (from "2")
		0,      // D1, seg f
		1 << 0, // D1, seg g (from "2")
		0,      // D1, dot
	}

	if !bytes.Equal(device.buffer[:], expectedBuffer[:]) {
		t.Errorf("FAIL: Buffer content after WriteString is wrong!\nExpected: %08b\nGot:      %08b", expectedBuffer, device.buffer)
	}

	// Now, test that Display() sends the correct data over I2C
	device.Display()

	expectedI2CData := append([]byte{0x00}, expectedBuffer[:]...)
	if !bytes.Equal(mockBus.data, expectedI2CData) {
		t.Errorf("FAIL: Data sent by Display() is wrong!\nExpected: %08b\nGot:      %08b", expectedI2CData, mockBus.data)
	}
}

// TestClearDisplay verifies that a single display can be cleared.
func TestClearDisplay(t *testing.T) {
	mockBus := &mockI2C{}
	device := New(mockBus, 0x70)

	// Write something to both displays first
	device.WriteString(0, "88")
	device.WriteString(1, "99")

	// Now clear display 0
	device.ClearDisplay(0)

	// Expected buffer state: Display 0 is clear, Display 1 still shows "99".
	var expectedBuffer [16]byte
	// Calculate "99" on display 1 (font[9] = 0b01101111)
	pattern9 := font[9]
	for seg := 0; seg < 7; seg++ {
		if (pattern9>>seg)&1 == 1 {
			// Set bit 0 (for first digit "9") and bit 1 (for second digit "9")
			expectedBuffer[8+seg] = (1 << 0) | (1 << 1)
		}
	}

	if !bytes.Equal(device.buffer[:], expectedBuffer[:]) {
		t.Errorf("FAIL: Buffer content after ClearDisplay is wrong!\nExpected: %08b\nGot:      %08b", expectedBuffer, device.buffer)
	}
}

// ExampleDevice_WriteString shows how to use the Device to write strings
// to both displays.
//
// ExampleDevice_WriteStringは、Deviceを使って両方のディスプレイに文字列を
// 書き込む方法を示す。
func ExampleDevice_WriteString() {
	// Create a mock I2C bus for demonstration.
	// In a real application, this would be machine.I2C0.
	// デモ用にモックのI2Cバスを作る。
	// 実際のアプリケーションでは、これはmachine.I2C0になる。
	mockBus := &mockI2C{}

	// Create and configure the display driver.
	// ディスプレイドライバを作成して設定する。
	display := New(mockBus, 0x70)
	display.Configure()

	// Write different strings to each 8-digit display.
	// それぞれの8桁ディスプレイに、違う文字列を書き込む。
	display.WriteString(0, "3600")
	display.WriteString(1, "1800")
	display.Display()

	fmt.Println("Wrote '3600' to display 0 and '1800' to display 1.")
	// Output: Wrote '3600' to display 0 and '1800' to display 1.
}
