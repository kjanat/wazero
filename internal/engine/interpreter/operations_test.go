package interpreter

import (
	"testing"

	"github.com/kjanat/wazero/internal/testing/require"
)

// TestInstructionName ensures that all the operation Kind's stringer is well-defined.
func TestOperationKind_String(t *testing.T) {
	for k := range operationKindEnd {
		require.NotEqual(t, "", k.String())
	}
}

// Test_unionOperation_String ensures that UnionOperation's stringer is well-defined for all supported OpKinds.
func Test_unionOperation_String(t *testing.T) {
	op := unionOperation{}
	for k := range operationKindEnd {
		op.Kind = k
		require.NotEqual(t, "", op.String())
	}
}

func TestLabel(t *testing.T) {
	for k := range labelKindNum {
		label := newLabel(k, 12345)
		require.Equal(t, k, label.Kind())
		require.Equal(t, 12345, label.FrameID())
	}
}
