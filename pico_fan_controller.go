//go:build pico

// pico_fan_controller.go

package main

import (
	"machine"
	"sync/atomic"

	"github.com/kou-tkbys/tk-fancon2/fan"
)

// pico専用のパルスカウント実装
type picoTachoCounter struct {
	pulseCount atomic.Uint32
}

func (p *picoTachoCounter) ReadAndReset() uint32 {
	return p.pulseCount.Swap(0)
}

func newPicoTachoCounter(pin machine.Pin) fan.PulseCounter {
	p := &picoTachoCounter{}
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	_ = pin.SetInterrupt(machine.PinFalling, func(m machine.Pin) {
		p.pulseCount.Add(1)
	})
	return p
}

// PicoFanController manages the dual contra-rotating fan.
type PicoFanController struct {
	Fans *fan.DualFan
	adc  machine.ADC
}

func NewPicoFanController() (*PicoFanController, error) {
	adc := machine.ADC{Pin: machine.ADC0}
	adc.Configure(machine.ADCConfig{})

	pwm := machine.PWM1
	// For 25kHz
	err := pwm.Configure(machine.PWMConfig{Period: 40000})
	if err != nil {
		return nil, err
	}

	// picoのピン情報を渡しつつカウンタを設定する
	counterF := newPicoTachoCounter(machine.GPIO4)
	counterR := newPicoTachoCounter(machine.GPIO5)

	fans := fan.NewDualFan("Typhoon", counterF, counterR)

	return &PicoFanController{
		Fans: fans,
		adc:  adc,
	}, nil
}

// UpdatePWM reads the value from the potentiometer and updates the PWM duty cycle.
func (fc *PicoFanController) UpdatePWM() {
	// Since PWM is only handled within this method, a local variable is sufficient.
	pwm := machine.PWM1
	potValue := fc.adc.Get()
	pwm.Set(0, uint32(potValue))
	pwm.Set(1, uint32(potValue))
}

// GetRPMs returns the calculated RPM values for both fans.
func (fc *PicoFanController) GetRPMs() (uint32, uint32) {
	return fc.Fans.CalculateRPMs()
}
