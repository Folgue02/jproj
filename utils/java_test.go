package utils_test

import (
	"testing"

	"github.com/folgue02/jproj/utils"
)

func TestDividePackageName(t *testing.T) {
    expected := []string { "com", "company", "app", "App" }
    result := utils.DividePackageName("com.company.app.App")

    if len(result) != len(expected) {
        t.Errorf("Different length (expected: %v, result: %v)", expected, result)
    }

    for i, e := range expected {
        if e != result[i] {
            t.Errorf("Element no.%d (expected: %v, result: %v)", i, expected, result)
        }
    }
}

func TestJavaPathToFilePath(t *testing.T) {
    expected := "com/company/app/App.java"
    result := utils.JavaPathToFilePath("com.company.app.App")

    if expected != result {
        t.Errorf("(expected: %s, result: %s)", expected, result)
    }
}

func TestToJavaFilePath(t *testing.T) {
    expected := "com/company/app/App.java"
    jpp := utils.NewJavaPackagePath("com.company.app.App")

    if jpp.ToJavaFilePath() != expected {
        t.Errorf("Expected: %v, Result: %v", expected, jpp.ToJavaFilePath())
    }
}

func TestToJavaDirPath(t *testing.T) {
    expected := "com/company/app"
    jpp := utils.NewJavaPackagePath("com.company.app")

    if jpp.ToJavaDirPath() != expected {
        t.Errorf("Expected: %v, Result: %v", expected, jpp.ToJavaDirPath())
    }
}
