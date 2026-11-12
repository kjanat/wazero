module github.com/kjanat/wazero/internal/integration_test/fuzz

go 1.27.0

toolchain go1.27.4

require github.com/kjanat/wazero v0.0.0

require golang.org/x/sys v0.44.0 // indirect

replace github.com/kjanat/wazero => ../../../
