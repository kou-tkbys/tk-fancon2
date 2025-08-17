// fan/dual.go

package fan

// DualFan represents a dual contra-rotating fan unit.
type DualFan struct {
	Name string
	// The front fan component
	Front *Fan
	// The rear fan component
	Rear *Fan
}

// NewDualFan creates a new DualFan instance.
func NewDualFan(name string, counterFront, counterRear PulseCounter) *DualFan {
	// It internally holds two Fan instances.
	df := &DualFan{
		Name:  name,
		Front: NewFan(name+"-F", counterFront),
		Rear:  NewFan(name+"-R", counterRear),
	}
	return df
}

// CalculateRPMs calculates the RPM for both fan components at once and returns the results.
func (df *DualFan) CalculateRPMs() (uint32, uint32) {
	frontRpm := df.Front.CalculateRPM()
	rearRpm := df.Rear.CalculateRPM()
	return frontRpm, rearRpm
}
