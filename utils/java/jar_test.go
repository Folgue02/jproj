package java_test

import (
	"fmt"
	"testing"

	"github.com/folgue02/jproj/utils/java"
)

func TestJarArguments(t *testing.T) {
	expected := []string{
		"--create", "--file", "file.jar",
		"--main-class", "me.user.app.App",
		"-C", "./src", ".", "-C", "./lib", "file.jar",
	}

	jarCommand := java.NewJarCommand(
		"jar",
		"file.jar",
		[][]string{{"./src"}, {"./lib", "file.jar"}},
		"me.user.app.App",
	)

	result := jarCommand.Arguments()
	printBoth := func() string {
		return fmt.Sprintf("Expected: %s, Result: %s", expected, result)
	}

	if len(result) != len(expected) {
		t.Errorf("Length of expected (%d) differs from the length of result (%d): %s", len(expected), len(result), printBoth())
		return
	}

	for i, expectedItem := range expected {
		if expectedItem != result[i] {
			t.Errorf("Elements differ in index '%d' (%s != %s): %s", i, expectedItem, result[i], printBoth())
		}
	}
}
