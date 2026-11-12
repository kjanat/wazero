module github.com/kjanat/wazero

// Floor Go version of wazero
go 1.27.0

toolchain go1.27.4

// All the beta tags are retracted and replaced with "pre" to prevent users
// from accidentally upgrading into the broken beta 1.
retract (
	v1.0.0-beta1
	v1.0.0-beta.3
	v1.0.0-beta.2
	v1.0.0-beta.1
)

require golang.org/x/sys v0.44.0
