//go:build pico

// fan/controller.go

package fan

import (
	"machine"
)

// 二重反転ファンコントローラ
// FanController manages the dual contra-rotating fan.
type FanController struct {
	Fans *DualFan
	adc  machine.ADC
}

func NewFanController() (*FanController, error) {
	adc := machine.ADC{Pin: machine.ADC0}
	adc.Configure(machine.ADCConfig{})

	pwm := machine.PWM1
	// For 25kHz
	err := pwm.Configure(machine.PWMConfig{Period: 40000})
	if err != nil {
		return nil, err
	}

	return &FanController{
		// The name "Typhoon" is so cool!
		Fans: NewDualFan("Typhoon", machine.GPIO4, machine.GPIO5),
		adc:  adc,
	}, nil
}

// UpdatePWM ポテンショメータの値を読んでPWMを更新する
// UpdatePWM reads the value from the potentiometer and updates the PWM duty cycle.
func (fc *FanController) UpdatePWM() {
	// このメソッドの中だけでPWMを扱うから、ローカル変数で十分
	// Since PWM is only handled within this method, a local variable is sufficient.
	pwm := machine.PWM1
	potValue := fc.adc.Get()
	pwm.Set(0, uint32(potValue))
	pwm.Set(1, uint32(potValue))
}

// GetRPMs 計算したRPM値を返却
// GetRPMs returns the calculated RPM values for both fans.
func (fc *FanController) GetRPMs() (uint32, uint32) {
	return fc.Fans.CalculateRPMs()
}
