package main

import (
	"fmt"
	"testing"
)

type genDepNodeModulePathTest struct {
	input     string
	depedency string
	expected  []string
}

func TestGenDepNodeModulePath(t *testing.T) {
	tests := []genDepNodeModulePathTest{
		{"node_modules/foo/node_modules/bar", "fizzbuzz", []string{
			"node_modules/foo/node_modules/bar/node_modules/fizzbuzz",
			"node_modules/foo/node_modules/fizzbuzz",
			"node_modules/fizzbuzz",
		}},
		{"node_modules/@clever/foo/node_modules/@clever/bar", "@types/fizzbuzz", []string{
			"node_modules/@clever/foo/node_modules/@clever/bar/node_modules/@types/fizzbuzz",
			"node_modules/@clever/foo/node_modules/@types/fizzbuzz",
			"node_modules/@types/fizzbuzz",
		}},
		{"", "@types/fizzbuzz", []string{
			"node_modules/@types/fizzbuzz",
		}},
		{"", "cron-service", []string{
			"node_modules/cron-service",
		}},
		{"local-package", "foobar", []string{
			"local-package/node_modules/foobar",
			"node_modules/foobar",
		}},
		{"some/local/package/node_modules/foo/node_modules/bar", "fizzbuzz", []string{
			"some/local/package/node_modules/foo/node_modules/bar/node_modules/fizzbuzz",
			"some/local/package/node_modules/foo/node_modules/fizzbuzz",
			"some/local/package/node_modules/fizzbuzz",
			"node_modules/fizzbuzz",
		}},
		{"some/local/package/node_modules/another/local/package/node_modules/foo", "fizzbuzz", []string{
			"some/local/package/node_modules/another/local/package/node_modules/foo/node_modules/fizzbuzz",
			"some/local/package/node_modules/another/local/package/node_modules/fizzbuzz",
			"some/local/package/node_modules/fizzbuzz",
			"node_modules/fizzbuzz",
		}},
	}

	for _, tt := range tests {
		results := genDepNodeModulePath(tt.input, tt.depedency)
		if len(tt.expected) != len(results) {
			t.Errorf("expected length to be %d, got=%d", len(tt.expected), len(results))
			fmt.Printf("results = %+v\n", results)
			continue
		}
		for i, result := range results {
			if tt.expected[i] != result {
				t.Errorf("index %d got=%q, want=%q", i, result, tt.expected[i])
			}
		}
	}
}
