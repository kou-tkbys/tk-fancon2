package fan

import (
	"fmt"
	"testing"
)

// パルスカウンターモック
type dualMockPulseCounter struct {
	mockCount uint32
}

func (m *dualMockPulseCounter) ReadAndReset() uint32 {
	return m.mockCount
}

func TestDualFan_CalculateRPMs(t *testing.T) {
	// 異なるモックを設定
	mockCounterF := &dualMockPulseCounter{mockCount: 120} // 3600 RPM になるはず
	mockCounterR := &dualMockPulseCounter{mockCount: 60}  // 1800 RPM になるはず

	// 二重反転ファン作成
	dualFan := NewDualFan("Test DualFan", mockCounterF, mockCounterR)

	// RPM計算
	frontRpm, rearRpm := dualFan.CalculateRPMs()

	// それぞれの値が正しいことをチェック
	expectedFrontRpm := uint32(3600)
	if frontRpm != expectedFrontRpm {
		t.Errorf("Frontの期待RPMは %d 、実際は %d で異なる", expectedFrontRpm, frontRpm)
	}

	expectedRearRpm := uint32(1800)
	if rearRpm != expectedRearRpm {
		t.Errorf("Rearの期待RPMは %d 、実際は %d で異なる", expectedRearRpm, rearRpm)
	}
}

// ExampleDualFan_CalculateRPMs shows how to use the DualFan type.
//
// ExampleDualFan_CalculateRPMsは、DualFan型の使い方を示す。
func ExampleDualFan_CalculateRPMs() {
	// Create mock counters for the front and rear fans.
	// 前後ファンのためのモックカウンターを作る。
	mockCounterF := &dualMockPulseCounter{mockCount: 120} // 3600 RPM
	mockCounterR := &dualMockPulseCounter{mockCount: 60}  // 1800 RPM

	// Create a new dual fan unit.
	// 新しい二重反転ファンユニットを作る。
	dualFan := NewDualFan("Typhoon", mockCounterF, mockCounterR)

	// Calculate the RPMs for both fans.
	// 両方のファンのRPMを計算する。
	frontRpm, rearRpm := dualFan.CalculateRPMs()

	fmt.Printf("Front: %d RPM, Rear: %d RPM\n", frontRpm, rearRpm)
	// Output: Front: 3600 RPM, Rear: 1800 RPM
}
