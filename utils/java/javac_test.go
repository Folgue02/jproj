package java

import (
	"fmt"
	"testing"
)

func TestArgsNoCP(t *testing.T) {
    expected := []string {
        "src/File.java",
        "-d", "target",
    }

    javacCommand := JavacCommand {
        CompilerPath: "javac",
        TargetDir: "target",
    }

    javacCommand.AddSources("src/File.java")

    result := javacCommand.Arguments()

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

func TestArgs(t *testing.T) {
    expected := []string {
        "src/File.java",
        "-d", "target",
        "-cp", "lib/dep.jar:lib/dep2.jar",
    }

    javacCommand := JavacCommand {
        CompilerPath: "javac",
        ClassPaths: []string { "lib/dep.jar", "lib/dep2.jar" },
        TargetDir: "target",
    }

    javacCommand.AddSources("src/File.java")

    result := javacCommand.Arguments()

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
