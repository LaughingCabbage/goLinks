package cmd

import (
	"os"
	"testing"
)

func TestVerifyPath(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	good, err := verifyPath(cwd)
	if !good || (err != nil) {
		t.Fail()
	}

	good, err = verifyPath(cwd + "/asdf")

	if good {
		t.Error("bad path should fail")
	}
}
