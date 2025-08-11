package main

import (
	"machine"
	"strconv"
	"time"

	"github.com/kou-tkbys/tk-fancon2/fan"
	"github.com/kou-tkbys/tk-fancon2/ht16k33"
)

// メモ：tinygo build -target=pico -o test.uf2
func main() {
	time.Sleep(1 * time.Second)

	// --- 部品の準備 ---

	// ファン制御の専門家を呼ぶ
	fanController, err := fan.NewFanController()
	if err != nil {
		println("Failed to create fan controller:", err.Error())
		return
	}

	// 表示の専門家を呼ぶ
	machine.I2C0.Configure(machine.I2CConfig{
		SDA: machine.I2C0_SDA_PIN, // GP4
		SCL: machine.I2C0_SCL_PIN, // GP5
	})
	display1 := ht16k33.New(*machine.I2C0, 0x70)
	display1.Configure()
	display2 := ht16k33.New(*machine.I2C0, 0x71)
	display2.Configure()

	// --- 主処理 ---
	rpmTicker := time.NewTicker(1 * time.Second)
	println(fanController.Fans.Name, " system, online. Starting application.")

	for {
		select {
		// 1秒ごとにタイマーが鳴ったら…
		case <-rpmTicker.C:
			// 1. ファン制御の専門家に、RPMの値を聞く
			rpm1, rpm2 := fanController.GetRPMs()
			// 1台目のディスプレイ
			display1.WriteString(strconv.Itoa(int(rpm1))) //数値を文字列に変換
			display1.Display()

			// 2台目のディスプレイ
			display2.WriteString(strconv.Itoa(int(rpm2)))
			display2.Display()

		// それ以外の時は…
		default:
			// ファン制御の専門家に、速度を更新するよう命令する
			fanController.UpdatePWM()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
