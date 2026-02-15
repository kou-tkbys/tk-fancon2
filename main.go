package main

import (
	// "strconv"
	"time"
	// "github.com/kou-tkbys/tk-fancon2/ht16k33"
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
	time.Sleep(1 * time.Second)

	// --- Prepare components ---
	// --- 準備 ---

	// Call in the fan control specialist
	// ファンの制御インスタンス生成
	fanController, err := NewFanController()
	if err != nil {
		println("Failed to create fan controller:", err.Error())
		return
	}

	// Call in the display specialist
	// ディスプレイ制御インスタンス生成
	// i2c := SetupI2C()

	// Initialize the dual display controlled by a single HT16K33 IC at address 0x70.
	// 2つのディスプレイを1つのIC(0x70)で制御
	// dualDisplay := ht16k33.New(i2c, 0x70)
	// dualDisplay.Configure()

	// --- Main processing loop ---
	// --- メインの処理ループ ---
	rpmTicker := time.NewTicker(rpmUpdateInterval)
	pwmTicker := time.NewTicker(pwmUpdateInterval)
	println(fanController.Fans.Name, " system, online. Starting application.")

	for {
		select {
		// Process triggered by a 1-second timer.
		// 1秒ごとのタイマー処理
		case <-rpmTicker.C: // This case is for slow updates (display).
			// Get fan RPMs.
			// ファンのRPMを取得
			rpm1, rpm2 := fanController.GetRPMs()
			println("Fan1:", rpm1, " Fan2:", rpm2)

			// Write the RPMs to displays 0 and 1 on the single device.
			// 1つのデバイスに、ディスプレイ0と1を指定して書き込む
			// dualDisplay.WriteString(0, strconv.Itoa(int(rpm1)))
			// dualDisplay.WriteString(1, strconv.Itoa(int(rpm2)))
			// Transfer the buffer to the display driver all at once.
			// 最後にまとめて転送！
			// dualDisplay.Display()

		case <-pwmTicker.C: // This case is for fast updates (potentiometer).
			// Control the fan rotation speed.
			// ファンの回転速度を制御
			fanController.UpdatePWM()
		}
	}
}
