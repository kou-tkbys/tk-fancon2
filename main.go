// main.go!!!
package main

import (
	"machine"
	"strconv"
	"time"

	"github.com/kou-tkbys/tk-fancon2/ht16k33"
)

// Note: tinygo test ./...
// Note: tinygo build -target=pico -o test.uf2
func main() {
	time.Sleep(1 * time.Second)

	// --- Prepare components ---

	// Call in the fan control specialist
	fanController, err := NewPicoFanController()
	if err != nil {
		println("Failed to create fan controller:", err.Error())
		return
	}

	// Call in the display specialist
	machine.I2C0.Configure(machine.I2CConfig{
		SDA: machine.I2C0_SDA_PIN, // GP2
		SCL: machine.I2C0_SCL_PIN, // GP3
	})
	display1 := ht16k33.New(machine.I2C0, 0x70)
	display1.Configure()
	display2 := ht16k33.New(machine.I2C0, 0x71)
	display2.Configure()

	// --- Main processing loop ---
	rpmTicker := time.NewTicker(1 * time.Second)
	println(fanController.Fans.Name, " system, online. Starting application.")

	for {
		select {
		// 1秒ごとのタイマー処理
		case <-rpmTicker.C:
			// ファンrpms取得
			rpm1, rpm2 := fanController.GetRPMs()

			// First display
			display1.WriteString(strconv.Itoa(int(rpm1)))
			display1.Display()

			// Second display
			display2.WriteString(strconv.Itoa(int(rpm2)))
			display2.Display()

		default:
			// ファンの回転速度制御
			fanController.UpdatePWM()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
