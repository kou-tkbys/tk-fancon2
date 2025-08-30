package fan

// PulseCounter is an interface that provides pulse counting
// functionality.
//
// PulseCounterは、パルスのカウント機能を提供するインターフェース
type PulseCounter interface {
	// ReadAndReset should be implemented to retrieve the current count
	// and reset the internal counter to 0.
	//
	// ReadAndResetには、現在までのカウント数を取得し、内部カウンタを0にリ
	// セットする処理を実装する
	ReadAndReset() uint32
}

// Fan represents a single fan unit.
type Fan struct {
	Name    string
	counter PulseCounter
	rpm     uint32
}

// NewFan creates a new Fan instance. The name can be any string.
//
// NewFanは、新しいFanインスタンスを作る。名前は好きな文字列で良い。
func NewFan(name string, counter PulseCounter) *Fan {
	return &Fan{
		Name:    name,
		counter: counter,
	}
}

// CalculateRPM retrieves the value from the internal counter, calculates
// the RPM, and returns it.
//
// CalculateRPMは、内部カウンタの値を取得し、RPMを計算して返却する。
func (f *Fan) CalculateRPM() uint32 {
	// Retrieve the value from the internal counter, calculate, and
	// return.
	//
	// 内部カウンタの値を取得し、計算したうえで返却する。
	count := f.counter.ReadAndReset()
	f.rpm = (count / 2) * 60
	return f.rpm
}
