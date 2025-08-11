// controller.go

package fan

import (
	"machine"
)

// 二重反転ファンコントローラ
type FanController struct {
	Fans *DualFan
	adc  machine.ADC
}

func NewFanController() (*FanController, error) {
	adc := machine.ADC{Pin: machine.ADC0}
	adc.Configure(machine.ADCConfig{})

	pwm := machine.PWM1
	err := pwm.Configure(machine.PWMConfig{Period: 40000})
	if err != nil {
		return nil, err
	}

	return &FanController{
		Fans: NewDualFan("台風", machine.GPIO4, machine.GPIO5),
		adc:  adc,
	}, nil
}

// UpdatePWM ポテンショメータの値を読んでPWMを更新する
func (fc *FanController) UpdatePWM() {
	// このメソッドの中だけでPWMを扱うから、ローカル変数で十分
	pwm := machine.PWM1
	potValue := fc.adc.Get()
	pwm.Set(0, uint32(potValue))
	pwm.Set(1, uint32(potValue))
}

// 計算したPWM値を返却
func (fc *FanController) GetRPMs() (uint32, uint32) {
	return fc.Fans.CalculateRPMs()
}
