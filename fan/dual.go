// fan/dual.go

package fan

import "machine"

// DualFan 二重反転ファン構造体
type DualFan struct {
	Name  string
	Front *Fan // 前側のファン部品
	Rear  *Fan // 後ろ側のファン部品
}

// NewDualFan 二重反転ファン
func NewDualFan(name string, pinFront, pinRear machine.Pin) *DualFan {
	// Fan２つを内部にもつ
	df := &DualFan{
		Name:  name,
		Front: NewFan(name+"-F", pinFront),
		Rear:  NewFan(name+"-R", pinRear),
	}
	return df
}

// CalculateRPMs 部品たちのRPMを一斉に計算して結果を取得する
func (df *DualFan) CalculateRPMs() (uint32, uint32) {
	frontRpm := df.Front.CalculateRPM()
	rearRpm := df.Rear.CalculateRPM()
	return frontRpm, rearRpm
}
