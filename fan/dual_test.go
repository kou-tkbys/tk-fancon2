package fan

import "testing"

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
