package fan

import "testing"

// Note: tinygo test ./fan

// PulseCounterモック
type singleMockPulseCounter struct {
	mockCount uint32
}

// 値をそのまま返却
func (m *singleMockPulseCounter) ReadAndReset() uint32 {
	return m.mockCount
}

// CalculateRPM test
func TestFan_CalculateRPM(t *testing.T) {
	testCases := []struct {
		name        string
		pulseCount  uint32 // mockに仕込むパルス値
		expectedRPM uint32 // 期待するRPM値
	}{
		{
			name:        "通常回転：1秒120パルス",
			pulseCount:  120,
			expectedRPM: 3600, // (120 / 2) * 60
		},
		{
			name:        "停止：0パルス",
			pulseCount:  0,
			expectedRPM: 0,
		},
		{
			name:        "低速回転：1秒30パルス",
			pulseCount:  30,
			expectedRPM: 900, // (30 / 2) * 60
		},
		{
			name:        "奇数パルス：1秒121パルス",
			pulseCount:  121,
			expectedRPM: 3600, // (121 / 2) * 60 = 60 * 60
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 準備
			mockCounter := &singleMockPulseCounter{mockCount: tc.pulseCount}
			testFan := NewFan("Test Fan", mockCounter)

			// 実行
			rpm := testFan.CalculateRPM()

			// 検証
			if rpm != tc.expectedRPM {
				t.Errorf("期待するRPMは %d 、実際は %d で異なる", tc.expectedRPM, rpm)
			}
		})
	}
}
