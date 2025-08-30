package fan

// DualFan represents a dual contra-rotating fan unit.
//
// DualFanは、二重反転ファンユニット
type DualFan struct {
	Name string
	// The front fan component
	// 前側のファン
	Front *Fan
	// The rear fan component
	// 後ろ側のファン
	Rear *Fan
}

// NewDualFan creates a new DualFan instance.
//
// NewDualFanは、新しいDualFanインスタンスを作る
func NewDualFan(name string, counterFront, counterRear PulseCounter) *DualFan {
	// It internally holds two Fan instances.
	// 内部的に2つのFanインスタンスを保持
	df := &DualFan{
		Name:  name,
		Front: NewFan(name+"-F", counterFront),
		Rear:  NewFan(name+"-R", counterRear),
	}
	return df
}

// CalculateRPMs calculates the RPM for both fan components at once and returns the results.
//
// CalculateRPMsは、両方のファン部品のRPMを一度に計算して結果を返す
func (df *DualFan) CalculateRPMs() (uint32, uint32) {
	frontRpm := df.Front.CalculateRPM()
	rearRpm := df.Rear.CalculateRPM()
	return frontRpm, rearRpm
}
