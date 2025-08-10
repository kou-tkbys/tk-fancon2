package main

import (
	"machine"
	"time"
)

// PWMクロック周波数設定について
// デフォルトでもモーター制御自体は特に問題なく可能だが、PWM周波数25kHz近辺に変更する
// これはDCモーターのPWM制御における一般的な推奨周波数に合わせるためである

// 参考：picoのデフォルト設定の場合の周波数
//
// クロック分周（clkdiv）:1（これはつまりシステムクロックをそのまま利用すること）
// ラップ値（wrap）:65535（PWMカウンターが利用できる最大値）
// 周波数計算：
//   125,000,000 Hz（システムクロック） / 1（clkdiv） * 65536（ラップ値＋１）≒ 1907Hz(約1.9kHz)
//
// つまりデフォルトのままだと約1.9kHzの周波数となる

// picoのPWM周波数を25kHzに変更する（tinygoで）
// PWMConfigに対して周期をナノ秒で指定する。周期のナノ秒変換式は以下
//
// 周期（秒）= 1 / 周波数（Hz） = 1 / 25,000 = 0.00004秒
// 0.00004秒 = 40,000ナノ秒
// これを設定することでC++で指定したのと同様な周波数設定が可能となる。

// 利用するぽテンションメーターについて
// 抵抗値が低すぎる：無駄な電流が流れる＝ぽテンションメーターが熱くなる、Picoに無駄な負担をかける
// 抵抗値が高すぎる：高すぎる抵抗（1Mなど）でPicoの内部抵抗（入力インピーダンス）に近づくと、Picoが電圧を正確に読み取れなくなる。
//
// これらの理由から、 5kΩ 〜 50kΩ の範囲のポテンショメータが、こういう使い方にはとってもバランスが良くて理想的。
// その中でも、10kΩ のポテンショメータが最も一般的で、どんな場面でも安心して使える「黄金の抵抗値」として扱われるのでだいたいみんな使ってる。

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

	err := pwm.Configure(machine.PWMConfig{
		// 周期をナノ秒で指定する (1 / 25,000 Hz = 40,000 ns)
		Period: 40000,
	})
	if err != nil {
		println("could not configure PWM")
		return
	}

	// 主処理
	for {
		// ポテンショメーターの値の読み取り
		// TinyGoのADC.Get()は、0から65535までの16ビット
		potValue := adc.Get()

		// 読み取った値をPWMのデューティー比として設定する
		// TinyGoのPWMの最大値は65535なので読み取った値をそのまま利用可能
		pwm.Set(0, uint32(potValue)) // チャンネルA (GPIO2)
		pwm.Set(1, uint32(potValue)) // チャンネルB (GPIO3)

		// 待ち
		time.Sleep(10 * time.Millisecond)
	}
}
