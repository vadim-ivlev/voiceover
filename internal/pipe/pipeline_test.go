package pipe

import "testing"

func TestStartTestPipeline(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test StartPipeline",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartTestPipeline()
		})
	}
}
