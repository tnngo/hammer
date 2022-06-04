package hammer

import (
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		vrs []hammer
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				vrs: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.vrs...)

			t.Log(got)
		})
	}
}
