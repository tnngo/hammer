package hammer

import (
	"testing"

	"github.com/tnngo/lad"
)

func TestNewUnitTest(t *testing.T) {
	UnitTest()
	lad.L().Debug("test")
}
