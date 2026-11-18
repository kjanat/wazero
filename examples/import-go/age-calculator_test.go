package main

import (
	"fmt"
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
	// Allow for the year rolling over while the module runs.
	for _, year := range []int{before, time.Now().Year()} {
		if stdout == fmt.Sprintf("println >> %d\nlog_i32 >> %d\n", year, year) {
			return
		}
	}
	t.Fatalf("overflowing override did not fall back to the clock: %q", stdout)
}
