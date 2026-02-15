package main

import (
	"machine"
	"strconv"
	"time"

	"github.com/kou-tkbys/tk-fancon2/ht16k33"
)

// Note: tinygo test ./...
// Note: tinygo build -target=pico -o test.uf2

const (
	rpmUpdateInterval = 1 * time.Second
	pwmUpdateInterval = 50 * time.Millisecond
)

// main is the entry point of the application.
// It initializes the fan controller and display, then enters an infinite
// loop to update fan speed and display RPMs.
//
// mainは、このアプリケーションのエントリーポイントファンコントローラーと
// ディスプレイを初期化し、無限ループに入ってファンの速度を更新、RPMの表示
// などを行う
func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// 1. 起動確認：ゆっくり3回点滅
	// これで「プログラムが走り出した」ことはわかるぞ。
	for i := 0; i < 3; i++ {
		led.High()
		time.Sleep(300 * time.Millisecond)
		led.Low()
		time.Sleep(300 * time.Millisecond)
	}

	// 2. ファンコントローラーの初期化
	// ここで死ぬなら、配線（特にGPIO周り）か初期化コードに問題があるぞ。
	fanController, err := NewFanController()
	if err != nil {
		// 初期化失敗なら高速点滅（SOS）じゃ！
		for {
			led.Set(!led.Get())
			time.Sleep(50 * time.Millisecond)
		}
	}

	// 3. 初期化成功：点灯しっぱなしで1秒待機
	led.High()
	time.Sleep(1 * time.Second)
	led.Low()

	// シリアルモニタの準備ができたら、高らかに宣言するのじゃ！
	println("Typhoon system, online. Starting application.")

	// Call in the display specialist
	// ディスプレイ制御インスタンス生成
	i2c := SetupI2C()

	// Initialize the dual display controlled by a single HT16K33 IC at address 0x70.
	// 2つのディスプレイを1つのIC(0x70)で制御
	dualDisplay := ht16k33.New(i2c, 0x70)
	dualDisplay.Configure()

	// --- Main processing loop ---
	rpmTicker := time.NewTicker(rpmUpdateInterval)
	pwmTicker := time.NewTicker(pwmUpdateInterval)

	for {
		select {
		case <-rpmTicker.C:
			rpm1, rpm2 := fanController.GetRPMs()
			println("Fan1:", rpm1, " Fan2:", rpm2)

			// Write the RPMs to displays 0 and 1 on the single device.
			// 1つのデバイスに、ディスプレイ0と1を指定して書き込む
			dualDisplay.WriteString(0, strconv.Itoa(int(rpm1)))
			dualDisplay.WriteString(1, strconv.Itoa(int(rpm2)))
			// Transfer the buffer to the display driver all at once.
			// 最後にまとめて転送！
			dualDisplay.Display()

		case <-pwmTicker.C:
			fanController.UpdatePWM()
			led.Set(!led.Get())
		}
	}
}
