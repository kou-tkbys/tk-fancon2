package ht16k33

import (
	"bytes"
	"testing"
)

// NOTE: tinygo test ./ht16k33
// ```-target Select the target you want to use. Leave it empty to compile for the host.````

// machine.I2Cのテスト用のモック
// This is a mock (a fake object) for testing that pretends to be machine.I2C
type mockI2C struct {
	addr uint16
	data []byte
}

// Txメソッドを偽装する。本物のI2C通信はせず、送られてきたデータを記録するだけ
// This fakes the Tx method. It doesn't perform real I2C communication,
// it just records the data that was supposed to be sent.
func (m *mockI2C) Tx(addr uint16, w, r []byte) error {
	m.addr = addr
	// w (write buffer)の内容をコピーして保存するだけ
	// Copy and save the contents of w (the write buffer)
	m.data = make([]byte, len(w))
	copy(m.data, w)
	return nil
}

// SetDigit test
func TestSetDigit(t *testing.T) {
	// table of test cases
	testCases := []struct {
		name         string // テストケースの名前
		position     int    // 入力：桁
		num          byte   // 入力：数字
		dot          bool   // 入力：ドットの有無
		expectedByte byte   // 期待するバッファの値
	}{
		{
			name:         "Digit 0, Number 8 with dot",
			position:     0,
			num:          8,
			dot:          true,
			expectedByte: 0b11111111, // "8" + dot
		},
		{
			name:         "Digit 1, Number 2 without dot",
			position:     1,
			num:          2,
			dot:          false,
			expectedByte: 0b01011011, // "2"
		},
		{
			name:         "Digit 7, Number 9 with dot",
			position:     7,
			num:          9,
			dot:          true,
			expectedByte: 0b11101111, // "9" + dot
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockBus := &mockI2C{}
			device := New(mockBus, 0x70)

			device.SetDigit(tc.position, tc.num, tc.dot)
			device.Display()

			// Create the expected full buffer
			expectedData := make([]byte, 17)
			expectedData[0] = 0x00 // アドレスバイト
			// 該当の桁（2バイトで1桁）に期待値を入れる
			// Set the expected value at the correct digit position (2 bytes per digit)
			expectedData[tc.position*2+1] = tc.expectedByte

			if !bytes.Equal(mockBus.data, expectedData) {
				t.Errorf("FAIL: 送信されたデータが違う\n期待した値: %08b\n実際の値:   %08b", expectedData, mockBus.data)
				t.Errorf("FAIL: The transmitted data is wrong!\nExpected: %08b\nGot:      %08b", expectedData, mockBus.data)
			}
		})
	}
}

// WriteString test
func TestWriteString(t *testing.T) {
	mockBus := &mockI2C{}
	device := New(mockBus, 0x70)

	// "1.2"と表示させてみる
	// Try to display "1.2"
	device.WriteString("1.2")
	device.Display()

	// 期待値
	expectedData := []byte{
		0x00,       // 先頭のアドレスバイト
		0b10000110, // "1" + dot
		0x00,
		0b01011011, // "2"
		0x00,       // ↓ここから後ろ13バイトは空のはず
		0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	if !bytes.Equal(mockBus.data, expectedData) {
		t.Errorf("FAIL: WriteStringで送信されたデータが違う\n期待した値: %08b\n実際の値:   %08b", expectedData, mockBus.data)
		t.Errorf("FAIL: Data sent by WriteString is wrong!\nExpected: %08b\nGot:      %08b", expectedData, mockBus.data)
	}
}
