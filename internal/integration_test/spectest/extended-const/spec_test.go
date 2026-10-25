package spectest

import (
	"context"
	"embed"
	"testing"

	"github.com/kjanat/wazero"
	"github.com/kjanat/wazero/api"
	"github.com/kjanat/wazero/experimental"
	"github.com/kjanat/wazero/internal/integration_test/spectest"
	"github.com/kjanat/wazero/internal/platform"
)

//go:embed testdata/*.wasm
//go:embed testdata/*.json
var testcases embed.FS

const enabledFeatures = api.CoreFeaturesV2 | experimental.CoreFeaturesExtendedConst

func TestCompiler(t *testing.T) {
	if !platform.CompilerSupported() {
		t.Skip()
	}
	spectest.Run(t, testcases, context.Background(), wazero.NewRuntimeConfigCompiler().WithCoreFeatures(enabledFeatures))
}

func TestInterpreter(t *testing.T) {
	spectest.Run(t, testcases, context.Background(), wazero.NewRuntimeConfigInterpreter().WithCoreFeatures(enabledFeatures))
}
