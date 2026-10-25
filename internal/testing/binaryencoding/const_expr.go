package binaryencoding

import (
	"github.com/kjanat/wazero/internal/wasm"
)

func encodeConstantExpression(expr wasm.ConstantExpression) (ret []byte) {
	ret = append(ret, expr.Data...)
	return
}
