// main.go!!!
package main

import (
	"machine"
	"strconv"
	"time"

	"github.com/kou-tkbys/tk-fancon2/fan"
	"github.com/kou-tkbys/tk-fancon2/ht16k33"
)

// Note: tinygo build -target=pico -o test.uf2
func main() {
	time.Sleep(1 * time.Second)

	// --- 部品の準備 ---
	// --- Prepare components ---

	// ファン制御の専門家を呼ぶ
	// Call in the fan control specialist
	fanController, err := fan.NewFanController()
	if err != nil {
		println("Failed to create fan controller:", err.Error())
		return
	}

	// 表示の専門家を呼ぶ
	// Call in the display specialist
	machine.I2C0.Configure(machine.I2CConfig{
		SDA: machine.I2C0_SDA_PIN, // GP2
		SCL: machine.I2C0_SCL_PIN, // GP3
	})
	display1 := ht16k33.New(machine.I2C0, 0x70)
	display1.Configure()
	display2 := ht16k33.New(machine.I2C0, 0x71)
	display2.Configure()

	// --- 主処理 ---
	// --- Main processing loop ---
	rpmTicker := time.NewTicker(1 * time.Second)
	println(fanController.Fans.Name, " system, online. Starting application.")

	for {
		select {
		// 1秒ごとのタイマー処理
		// When the ticker fires every second...
		case <-rpmTicker.C:
			// 1. ファン制御の専門家に、RPMの値を聞く
			// 1. Ask the fan control specialist for the RPM values
			rpm1, rpm2 := fanController.GetRPMs()

			// 1台目のディスプレイ
			// First display
			//
			// 数値を文字列に変換
			// Convert the number to a string
			display1.WriteString(strconv.Itoa(int(rpm1)))
			display1.Display()

			// 2台目のディスプレイ
			// Second display
			display2.WriteString(strconv.Itoa(int(rpm2)))
			display2.Display()

		// それ以外の時
		// At all other times...
		default:
			// ファン制御の専門家に、速度を更新するよう命令する
			// Command the fan control specialist to update the speed
			fanController.UpdatePWM()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
