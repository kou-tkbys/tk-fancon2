// fan/single.go

package fan

import (
	"machine"
	"sync/atomic"
)

// タコピンに必要な機能だけインターフェースで分離
// TachoPin is an interface that extracts only the functions needed for a tacho pin.
type TachoPin interface {
	Configure(config machine.PinConfig)
	SetInterrupt(trigger machine.PinChange, callback func(machine.Pin)) error
}

// Fan ファン単体の構造体
// Fan represents a single fan unit.
type Fan struct {
	Name       string
	tachoPin   TachoPin
	pulseCount atomic.Uint32
	rpm        uint32
}

// NewFan nameは自由
// NewFan creates a new Fan instance. The name can be any string.
func NewFan(name string, pin TachoPin) *Fan {
	f := &Fan{
		Name:     name,
		tachoPin: pin,
	}

	// 回転数取得用のタコメータピン（タコピン）のプルアップ設定
	// Configure the tacho pin (for RPM measurement) as an input with a pull-up resistor.
	f.tachoPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// 回転数取得用のタコメータピンに対する割り込みコールバック設定
	// Set up an interrupt callback for the tacho pin to count pulses.
	//
	// pinFallingを条件にパルスカウントアップ
	// Increment the pulse count on each falling edge.
	_ = f.tachoPin.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		f.pulseCount.Add(1)
	})
	return f
}

// CalculateRPM RPMを計算して返却する
// CalculateRPM calculates and returns the fan's speed in RPM.
func (f *Fan) CalculateRPM() uint32 {
	// 計算・取得と同時にカウンタをリフレッシュ
	// Atomically read the count and reset the counter to zero.
	count := f.pulseCount.Swap(0)
	// 2 pulses per revolution, 60 seconds per minute.
	f.rpm = (uint32(count) / 2) * 60
	return f.rpm
}
