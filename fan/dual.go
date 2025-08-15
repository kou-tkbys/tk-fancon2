// fan/dual.go

package fan

// DualFan 二重反転ファン構造体
// DualFan represents a dual contra-rotating fan unit.
type DualFan struct {
	Name string
	// 前側のファン部品
	// The front fan component
	Front *Fan
	// 後ろ側のファン部品
	// The rear fan component
	Rear *Fan
}

// NewDualFan 二重反転ファン
// NewDualFan creates a new DualFan instance.
func NewDualFan(name string, counterFront, counterRear PulseCounter) *DualFan {

	// Fan２つを内部にもつ
	// It internally holds two Fan instances.
	df := &DualFan{
		Name:  name,
		Front: NewFan(name+"-F", counterFront),
		Rear:  NewFan(name+"-R", counterRear),
	}
	return df
}

// CalculateRPMs 部品たちのRPMを一斉に計算して結果を取得する
// CalculateRPMs calculates the RPM for both fan components at once and returns the results.
func (df *DualFan) CalculateRPMs() (uint32, uint32) {
	frontRpm := df.Front.CalculateRPM()
	rearRpm := df.Rear.CalculateRPM()
	return frontRpm, rearRpm
}
