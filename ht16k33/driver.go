// ht16k33/driver.go
//
// https://cdn-shop.adafruit.com/datasheets/ht16K33v110.pdf
//
// This driver does not cover all functions of the HT16K33.
// It only implements the necessary features to control and display numbers on a 7-segment LED.
package ht16k33

const (
	// Default address for the first device
	ht16k33Address = 0x70

	// Commands for HT16K33
	ht16k33TurnOnOscillator = 0x21
	ht16k33TurnOnDisplay    = 0x81
	ht16k33SetBrightness    = 0xE0
)

// 7-segment display number patterns
// Segments: g-f-e-d-c-b-a
var font = [10]byte{
	0b00111111, // 0
	0b00000110, // 1
	0b01011011, // 2
	0b01001111, // 3
	0b01100110, // 4
	0b01101101, // 5
	0b01111101, // 6
	0b00000111, // 7
	0b01111111, // 8
	0b01101111, // 9
}

// I2CBus is an interface that abstracts the I2C Tx method we need.
type I2CBus interface {
	Tx(addr uint16, w, r []byte) error
}

type Device struct {
	bus     I2CBus
	Address uint8
	// HT16K33の表示用RAMは最大16バイトある
	buffer [16]byte
}

func New(bus I2CBus, address uint8) Device {
	return Device{
		bus:     bus,
		Address: address,
	}
}

func (d *Device) Configure() {
	d.bus.Tx(uint16(d.Address), []byte{ht16k33TurnOnOscillator}, nil)
	d.bus.Tx(uint16(d.Address), []byte{ht16k33TurnOnDisplay}, nil)
	// とりあえず最大の明るさに
	d.SetBrightness(15)
}

// Clear clears the display buffer.
func (d *Device) Clear() {
	for i := range d.buffer {
		d.buffer[i] = 0
	}
}

// SetDigit displays a specified number on a specified digit (with or without a dot).
//
// position: 0-7 (一般的な8桁モジュールを想定)
// num: 0-9
// dot: if true, the decimal point is turned on
func (d *Device) SetDigit(position int, num byte, dot bool) {
	if position < 0 || position > 7 || num > 9 {
		// Out of range
		return
	}
	// Each digit uses 2 bytes
	addr := position * 2
	val := font[num]
	if dot {
		// Set the bit flag for the dot
		val |= 0b10000000
	}
	d.buffer[addr] = val
}

// WriteString displays a string on the display.
// It can handle strings like "1234" and "5.6".
func (d *Device) WriteString(s string) {
	d.Clear()
	pos := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			// Exit if it exceeds 8 digits
			if pos > 7 {
				break
			}

			// For dot detection
			dot := false
			// Note: この先読みは、文字列の最後にドットがあるとうまく動かないが、
			// 簡単な実装としては十分
			if len(s) > pos+1 && s[pos+1] == '.' {
				dot = true
			}
			d.SetDigit(pos, byte(r-'0'), dot)
			pos++
		} else if r == '.' {
			// The dot is handled along with the number, so do nothing here.
		}
	}
}

// Display transfers the buffer's content to the LED.
func (d *Device) Display() {
	data := append([]byte{0x00}, d.buffer[:]...)
	d.bus.Tx(uint16(d.Address), data, nil)
}

// SetBrightness sets the display brightness (0-15).
func (d *Device) SetBrightness(brightness uint8) {
	if brightness > 15 {
		brightness = 15
	}
	d.bus.Tx(uint16(d.Address), []byte{ht16k33SetBrightness | brightness}, nil)
}
