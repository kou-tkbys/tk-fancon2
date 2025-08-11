package fan

import (
	"machine"
	"sync/atomic"
)

// Fan ファン単体の構造体
type Fan struct {
	Name       string
	tachoPin   machine.Pin
	pulseCount atomic.Uint32
	rpm        uint32
}

// NewFan nameは自由
func NewFan(name string, pin machine.Pin) *Fan {
	f := &Fan{
		Name:     name,
		tachoPin: pin,
	}

	// 回転数取得用のタコメータピン（タコピン）のプルアップ設定
	f.tachoPin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// 回転数取得用のタコメータピンに対する割り込みコールバック設定
	// pinFallingを条件にパルスカウントアップ
	f.tachoPin.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		f.pulseCount.Add(1)
	})
	return f
}

// CalculateRPM  RPMを計算して返却する
func (f *Fan) CalculateRPM() uint32 {
	count := f.pulseCount.Swap(0) // 計算・取得と同時にカウンタをリフレッシュ
	f.rpm = (uint32(count) / 2) * 60
	return f.rpm
}
