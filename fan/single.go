// fan/single.go

package fan

// PulseCounter は、パルスのカウント機能を提供するインターフェース
type PulseCounter interface {
	// ReadAndReset には、現在までのカウント数を取得し、内部カウンタを0にリセットする処理を実装する
	ReadAndReset() uint32
}

// Fan represents a single fan unit.
type Fan struct {
	Name    string
	counter PulseCounter
	rpm     uint32
}

// NewFan creates a new Fan instance. The name can be any string.
func NewFan(name string, counter PulseCounter) *Fan {
	return &Fan{
		Name:    name,
		counter: counter,
	}
}

func (f *Fan) CalculateRPM() uint32 {
	// 内部カウンタの値を取得し、計算したうえで返却する
	count := f.counter.ReadAndReset()
	f.rpm = (count / 2) * 60
	return f.rpm
}
