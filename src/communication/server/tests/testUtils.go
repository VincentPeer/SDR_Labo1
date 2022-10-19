package tests

import (
	"path/filepath"
	"testing"
)

const (
	TEST_CONFIG_FILE_PATH = "../tests/config_test.json"
)

func GetTestData(t *testing.T) string {
	path, err := filepath.Abs(TEST_CONFIG_FILE_PATH)

	if err != nil {
		t.Error(err)
	}
	return path
}
