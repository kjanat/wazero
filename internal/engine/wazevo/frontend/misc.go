package frontend

import (
	"github.com/kjanat/wazero/internal/engine/wazevo/ssa"
	"github.com/kjanat/wazero/internal/wasm"
)

func FunctionIndexToFuncRef(idx wasm.Index) ssa.FuncRef {
	return ssa.FuncRef(idx)
}
