//go:build esp32

package main

import (
	"machine"
	"sync/atomic"
	"time"

	"github.com/kou-tkbys/tk-fancon2/fan"
)

// esp32TachoCounter is an ESP32-specific implementation for counting pulses.
//
// esp32TachoCounterは、ESP32専用のパルスカウント実装じゃ。
type esp32TachoCounter struct {
	pulseCount atomic.Uint32
}

// ReadAndReset reads the current pulse count and resets it to zero atomically.
func (p *esp32TachoCounter) ReadAndReset() uint32 {
	return p.pulseCount.Swap(0)
}

// newESP32TachoCounter creates a new pulse counter for a given pin.
func newESP32TachoCounter(pin machine.Pin) fan.PulseCounter {
	p := &esp32TachoCounter{}
	pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// ESP32 in TinyGo does not support SetInterrupt yet.
	// Use polling with a goroutine instead.
	go func() {
		lastState := pin.Get()
		for {
			currentState := pin.Get()
			if lastState && !currentState { // Falling edge (High -> Low)
				p.pulseCount.Add(1)
			}
			lastState = currentState
			// Check every 100 microseconds (10kHz sampling is enough for fans)
			time.Sleep(100 * time.Microsecond)
		}
	}()

	return p
}

// ESPFanController manages the dual contra-rotating fan on ESP32.
//
// ESPFanControllerは、ESP32上で二重反転ファンを管理するのじゃ。
type ESPFanController struct {
	Fans *fan.DualFan
	pinF machine.Pin
	pinR machine.Pin
}

// NewFanController creates and configures a new fan controller for ESP32.
//
// NewFanControllerは、ESP32用のファンコントローラーを作成して設定するぞ。
func NewFanController() (*ESPFanController, error) {
	pinF := machine.GPIO18
	pinR := machine.GPIO19
	pinF.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pinR.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Set up counters (Using GPIO16 and GPIO17 for tacho)
	counterF := newESP32TachoCounter(machine.GPIO16)
	counterR := newESP32TachoCounter(machine.GPIO17)

	fans := fan.NewDualFan("Typhoon-ESP", counterF, counterR)

	return &ESPFanController{
		Fans: fans,
		pinF: pinF,
		pinR: pinR,
	}, nil
}

// UpdatePWM reads the value from the potentiometer and updates the PWM duty cycle.
func (fc *ESPFanController) UpdatePWM() {
	fc.pinF.High()
	fc.pinR.High()
}

// GetRPMs returns the calculated RPM values for both fans.
func (fc *ESPFanController) GetRPMs() (uint32, uint32) {
	return fc.Fans.CalculateRPMs()
}

// SetupI2C configures the I2C bus for ESP32.
func SetupI2C() *machine.I2C {
	machine.I2C0.Configure(machine.I2CConfig{
		SDA: machine.GPIO21,
		SCL: machine.GPIO22,
	})
	return machine.I2C0
}
