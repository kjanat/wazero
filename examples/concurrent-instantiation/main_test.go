package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/kjanat/wazero/internal/testing/maintester"
	"github.com/kjanat/wazero/internal/testing/require"
)

// Test_main ensures the following will work:
//
//	go run main.go
func Test_main(t *testing.T) {
	stdout, _ := maintester.TestMain(t, main, "main")

	for i := range 50 {
		require.True(t, strings.Contains(stdout, strconv.Itoa(i*2)))
	}
}
