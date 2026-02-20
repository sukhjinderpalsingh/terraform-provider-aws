// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package flex

import (
	"strings"
	"testing"

	testingiface "github.com/mitchellh/go-testing-interface"
)

func TestGenerateGoldenPath(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		fullTestName string
		expectedPath string
	}{
		"root test": {
			fullTestName: "TestRoot",
			expectedPath: "autoflex/unknown/root/testroot.golden",
		},
		"single level": {
			fullTestName: "TestRoot/test case",
			expectedPath: "autoflex/unknown/root/test_case.golden",
		},
		"multiple levels": {
			fullTestName: "TestRoot/Outer Test-Case/inner test case",
			expectedPath: "autoflex/unknown/root/outer_testcase/inner_test_case.golden",
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			fooT := &testingiface.RuntimeT{}

			subtestName := testCase.fullTestName
			if last := strings.LastIndex(subtestName, "/"); last > 0 {
				subtestName = subtestName[last+1:]
			}
			path := autoGenerateGoldenPath(fooT, testCase.fullTestName, subtestName)

			if path != testCase.expectedPath {
				t.Errorf("Incorrect path %q, expected %q", path, testCase.expectedPath)
			}
		})
	}
}

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
