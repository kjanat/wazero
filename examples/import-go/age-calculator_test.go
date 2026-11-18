package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kjanat/wazero/internal/testing/maintester"
	"github.com/kjanat/wazero/internal/testing/require"
)

// Test_main ensures the following will work:
//
//	go run age-calculator.go 2000
func Test_main(t *testing.T) {
	// Set ENV to ensure this test doesn't need maintenance every year.
	t.Setenv("CURRENT_YEAR", "2021")

	stdout, _ := maintester.TestMain(t, main, "age-calculator", "2000")
	require.Equal(t, `println >> 21
log_i32 >> 21
`, stdout)
}

func TestOverflowingYearOverrideUsesClock(t *testing.T) {
	t.Setenv("CURRENT_YEAR", "4294967296")
	before := time.Now().Year()
	stdout, _ := maintester.TestMain(t, main, "age-calculator", "0")
	after := time.Now().Year()
	// The two Wasm calls read the clock separately, so New Year can fall between them.
	same := func(year int) string {
		return fmt.Sprintf("println >> %d\nlog_i32 >> %d\n", year, year)
	}
	if stdout == same(before) || stdout == same(after) {
		return
	}
	if before != after && stdout == fmt.Sprintf("println >> %d\nlog_i32 >> %d\n", before, after) {
		return
	}
	if strings.Contains(stdout, "println >> 0\n") || strings.Contains(stdout, "log_i32 >> 0\n") {
		t.Fatalf("wrapped year 0 must be rejected: %q", stdout)
	}
	t.Fatalf("overflowing override did not fall back to the clock: %q", stdout)
}
