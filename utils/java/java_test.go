package java_test

import (
	"fmt"
	"testing"

	"github.com/folgue02/jproj/utils/java"
)

func TestArgs(t *testing.T) {
    expected := []string {
        "-cp", "./target/:./lib/file.jar",
        "me.user.app.App",
        "--my-flag", "--another-flag",
    }

    javaCommand := java.JavaCommand {
        JavaPath: "java",
        MainClass: "me.user.app.App",
        ClassPaths: []string { "./target/", "./lib/file.jar" },
        ExecArguments: []string { "--my-flag", "--another-flag" },
    }

    result := javaCommand.Arguments()

    printBoth := func() string {
        return fmt.Sprintf("Expected: %s, Result: %s", expected, result)
    }

    if len(result) != len(expected)  {
        t.Errorf("Length of expected (%d) differs from the length of result (%d): %s", len(expected), len(result), printBoth())
    }

    for i, expectedItem := range expected {
        if expectedItem != result[i] {
            t.Errorf("Elements differ in index '%d' (%s != %s): %s", i, expectedItem, result[i], printBoth())
        }
    }
}
