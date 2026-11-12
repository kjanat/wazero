package experimental_test

import (
	"context"
	"math"
	"testing"

	"github.com/kjanat/wazero"
	"github.com/kjanat/wazero/experimental"
	"github.com/kjanat/wazero/internal/testing/binaryencoding"
	"github.com/kjanat/wazero/internal/testing/require"
)

func TestTableLimiterRejectsHostSizeBeforeReservation(t *testing.T) {
	ctx := context.Background()
	budget := newTableBudget(math.MaxUint32)
	func() {
		r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
		defer r.Close(ctx)
		mod, err := r.Instantiate(experimental.WithTableLimiter(ctx, budget), binaryencoding.EncodeModule(tableGrowModule(0, nil, false)))
		require.NoError(t, err)
		// A billion four-byte references cannot fit in a 32-bit allocation.
		result, err := mod.ExportedFunction("grow").Call(ctx, 1<<30)
		require.NoError(t, err)
		require.Equal(t, []uint64{math.MaxUint32}, result)
		budget.check(t, 0, [3]uint32{0, 0, math.MaxUint32})
		result, err = mod.ExportedFunction("grow").Call(ctx, 1)
		require.NoError(t, err)
		require.Equal(t, []uint64{0}, result)
		budget.check(t, 1, [3]uint32{0, 0, math.MaxUint32}, [3]uint32{0, 1, math.MaxUint32})
	}()
	awaitTableRelease(t, budget, 1)
}

func TestTableGrowHostSizeWithoutLimiter(t *testing.T) {
	ctx := context.Background()
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
	defer r.Close(ctx)
	mod, err := r.Instantiate(ctx, binaryencoding.EncodeModule(tableGrowModule(0, nil, false)))
	require.NoError(t, err)
	result, err := mod.ExportedFunction("grow").Call(ctx, 1<<30)
	require.NoError(t, err)
	require.Equal(t, []uint64{math.MaxUint32}, result)
	result, err = mod.ExportedFunction("grow").Call(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, []uint64{0}, result)
}
