package hammer

import (
	"testing"
)

func Test_yml_Load(t *testing.T) {
	t.Run("yml_Load", func(t *testing.T) {
		y := &yml{}
		got, err := y.load()
		if err != nil {
			t.Error(err)
			return
		}

		t.Log(got)
	})
}
