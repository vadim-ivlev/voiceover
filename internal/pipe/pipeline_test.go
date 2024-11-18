package pipe

import "testing"

func TestStartTestPipeline(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test StartTestPipeline",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartTestPipeline()
		})
	}
}
