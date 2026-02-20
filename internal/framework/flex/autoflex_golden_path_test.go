// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package flex

import (
	"testing"
)

func TestNormalizeTestName(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testName string
		expected string
	}{
		// This shouldn't happen, tests always start with "Test"
		"no prefix": {
			testName: "ImpossibleTestCase",
			expected: "impossible_test_case",
		},
		"normal": {
			testName: "TestExpandLogging_collections",
			expected: "expand_logging_collections",
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			path := normalizeTestName(testCase.testName)

			if path != testCase.expected {
				t.Errorf("Incorrect name %q, expected %q", path, testCase.expected)
			}
		})
	}
}

func TestNormalizeTestCaseName(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		testName string
		expected string
	}{
		"star": {
			testName: "*struct",
			expected: "pointer_struct",
		},
		"with spaces": {
			testName: "with a space",
			expected: "with_a_space",
		},
		"with upper case": {
			testName: "With Uppercase",
			expected: "with_uppercase",
		},
		"with hyphen": {
			testName: "with-hyphen",
			expected: "withhyphen",
		},
		"with comma": {
			testName: "with,comma",
			expected: "withcomma",
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			path := normalizeTestCaseName(testCase.testName)

			if path != testCase.expected {
				t.Errorf("Incorrect name %q, expected %q", path, testCase.expected)
			}
		})
	}
}
