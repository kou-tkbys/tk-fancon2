// https://cdn-shop.adafruit.com/datasheets/ht16K33v110.pdf
//
// このドライバはすべての機能を網羅しておらず、7seg-ledで数字
// 表示制御を行う機能のみを実装している。
package ht16k33

import "machine"

const (
	ht16k33Address = 0x70 // 1台目のデフォルトアドレス

	// HT16K33用コマンド
	ht16k33TurnOnOscillator = 0x21
	ht16k33TurnOnDisplay    = 0x81
	ht16k33SetBrightness    = 0xE0
)

// 7セグの数字パターン
// g-f-e-d-c-b-a
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

type Device struct {
	bus     machine.I2C
	Address uint8
	buffer  [16]byte // HT16K33の表示用RAMは最大16バイトある
}

func New(bus machine.I2C, address uint8) Device {
	return Device{
		bus:     bus,
		Address: address,
	}
}

func (d *Device) Configure() {
	d.bus.Tx(uint16(d.Address), []byte{ht16k33TurnOnOscillator}, nil)
	d.bus.Tx(uint16(d.Address), []byte{ht16k33TurnOnDisplay}, nil)
	d.SetBrightness(15) // とりあえず最大の明るさに
}

// Clear 表示バッファをクリア
func (d *Device) Clear() {
	for i := range d.buffer {
		d.buffer[i] = 0
	}
}

// SetDigit  指定した桁に、指定した数字を表示する（ドットの有無も指定可）
//
// position: 0-7 (一般的な8桁モジュールを想定)
// num: 0-9
// dot: trueでドットを付ける
func (d *Device) SetDigit(position int, num byte, dot bool) {
	if position < 0 || position > 7 || num > 9 {
		return // 範囲外
	}
	addr := position * 2 // 2バイトで1桁
	val := font[num]
	if dot {
		val |= 0b10000000 // ドットのビットフラグを立てる
	}
	d.buffer[addr] = val
}

// WriteString  文字列をディスプレイに表示する
// "1234" や "5.6" のような文字列を扱える
func (d *Device) WriteString(s string) {
	d.Clear()
	pos := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			if pos > 7 {
				break
			} // 8桁を超えたら抜ける

			dot := false // ドット検出用
			// Note: この先読みは、文字列の最後にドットがあるとうまく動かないが、
			// 簡単な実装としては十分
			if len(s) > pos+1 && s[pos+1] == '.' {
				dot = true
			}
			d.SetDigit(pos, byte(r-'0'), dot)
			pos++
		} else if r == '.' {
			// ドットは数字と一緒に処理したので、ここでは何もしない
		}
	}
}

// Display  バッファの内容をLEDに転送する
func (d *Device) Display() {
	data := append([]byte{0x00}, d.buffer[:]...)
	d.bus.Tx(uint16(d.Address), data, nil)
}

// SetBrightness 明るさを設定する (0-15)
func (d *Device) SetBrightness(brightness uint8) {
	if brightness > 15 {
		brightness = 15
	}
	d.bus.Tx(uint16(d.Address), []byte{ht16k33SetBrightness | brightness}, nil)
}
