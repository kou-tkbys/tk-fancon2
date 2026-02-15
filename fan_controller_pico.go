//go:build pico

package main

import (
	"machine"
	"sync/atomic"

	"github.com/kou-tkbys/tk-fancon2/fan"
)

// picoTachoCounter is a Pico-specific implementation for counting pulses.
// It uses atomic operations to safely increment the count from an
// interrupt.
//
// picoTachoCounterは、Pico専用のパルスカウント実装。
// 割り込みから安全にカウントを増やすために、アトミック操作を使う。
type picoTachoCounter struct {
	pulseCount atomic.Uint32
}

// ReadAndReset reads the current pulse count and resets it to zero
// atomically.
//
// ReadAndResetは、現在のパルス数を読み取って、アトミックにゼロにリセット
// する。
func (p *picoTachoCounter) ReadAndReset() uint32 {
	return p.pulseCount.Swap(0)
}

// newPicoTachoCounter creates a new pulse counter for a given pin.
// It configures the pin as an input with a pull-up and sets up a
// falling-edge interrupt.
//
// newPicoTachoCounterは、指定されたピンのための新しいパルスカウンターを作
// る。ピンをプルアップ付きの入力として設定し、立ち下がりエッジの割り込み
// を設定する。
func newPicoTachoCounter(pin machine.Pin) fan.PulseCounter {
	p := &picoTachoCounter{}
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	_ = pin.SetInterrupt(machine.PinFalling, func(m machine.Pin) {
		p.pulseCount.Add(1)
	})
	return p
}

// PicoFanController manages the dual contra-rotating fan.
//
// PicoFanControllerは、二重反転ファンを管理する。
type PicoFanController struct {
	Fans *fan.DualFan
	adc  machine.ADC
}

// NewFanController creates and configures a new fan controller.
// It sets up the ADC for the potentiometer, configures PWM for a 25kHz
// frequency, and initializes the pulse counters for both fans.
//
// NewFanControllerは、新しいファンコントローラーを作成して設定する。
// ポテンショメータ用のADCを設定し、25kHzの周波数でPWMを設定し、両方のファ
// ンのパルスカウンターを初期化する。
func NewFanController() (*PicoFanController, error) {
	adc := machine.ADC{Pin: machine.ADC0}
	adc.Configure(machine.ADCConfig{})

	pwm := machine.PWM1
	// For 25kHz
	err := pwm.Configure(machine.PWMConfig{Period: 40000})
	if err != nil {
		return nil, err
	}

	// Set up counters by passing Pico's pin information.
	// picoのピン情報を渡しつつカウンタを設定
	counterF := newPicoTachoCounter(machine.GPIO4)
	counterR := newPicoTachoCounter(machine.GPIO5)

	fans := fan.NewDualFan("Typhoon", counterF, counterR)

	return &PicoFanController{
		Fans: fans,
		adc:  adc,
	}, nil
}

// UpdatePWM reads the value from the potentiometer and updates the PWM
// duty cycle.
//
// UpdatePWMは、ポテンショメータから値を読み取り、PWMのデューティサイクル
// を更新する。
func (fc *PicoFanController) UpdatePWM() {
	// Since PWM is only handled within this method, a local variable is
	// sufficient.
	//
	// PWMはこのメソッド内でしか扱わないので、ローカル変数で十分。
	pwm := machine.PWM1
	potValue := fc.adc.Get()
	pwm.Set(0, uint32(potValue))
	pwm.Set(1, uint32(potValue))
}

// GetRPMs returns the calculated RPM values for both fans.
//
// GetRPMsは、計算された両方のファンのRPM値を返す。
func (fc *PicoFanController) GetRPMs() (uint32, uint32) {
	return fc.Fans.CalculateRPMs()
}

// SetupI2C configures the I2C bus for Pico.
//
// SetupI2Cは、Pico用のI2Cバスを設定する。
func SetupI2C() *machine.I2C {
	machine.I2C0.Configure(machine.I2CConfig{
		SDA: machine.GPIO0, // GP0 (I2C0 SDA)
		SCL: machine.GPIO1, // GP1 (I2C0 SCL)
	})
	return machine.I2C0
}
