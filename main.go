package main

import (
	"machine"
	"time"
)

// メモ：tinygo build -target=pico -o test.uf2
func main() {
	// 開始待機
	time.Sleep(1 * time.Second)

	// ポテンショメーターをつなぐピン
	// GPIO26:PicoだとADC0というアナログ入力ピン
	adc := machine.ADC{Pin: machine.ADC0}
	adc.Configure(machine.ADCConfig{})

	// PWMファン制御用のピン設定
	// GPIO2とGPIO3はPicoだとPWMスライス1グループ
	// つまり、machine.PWM1を選ぶと、この2つのピンを一緒に扱える
	pwm := machine.PWM1
	pwm.Configure(machine.PWMConfig{}) // デフォルトの設定をそのまま利用

	// 主処理
	for {
		// ポテンショメーターの値の読み取り
		// TinyGoのADC.Get()は、0から65535までの16ビット
		potValue := adc.Get()

		// 読み取った値をPWMのデューティー比として設定する
		// TinyGoのPWMの最大値は65535なので読み取った値をそのまま使える？
		pwm.Set(0, uint32(potValue)) // チャンネルA (GPIO2)
		pwm.Set(1, uint32(potValue)) // チャンネルB (GPIO3)

		// 待ち
		time.Sleep(10 * time.Millisecond)
	}
}
