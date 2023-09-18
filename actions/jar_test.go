package actions

import (
	"fmt"
	"path"
	"testing"

	"github.com/folgue02/jproj/configuration"
)

func TestBuildJarCommand(t *testing.T) {
    jarConfig := JarCommandConfiguration { Directory: "." }
    projectConfig := configuration.Configuration {
        ProjectName: "testing",
        ProjectTarget: "target",
    }

    expected := []string { "--create", "--file", "testing.jar", 
        "--manifest", path.Join(jarConfig.Directory, "target", "MANIFEST.mf"),
        "-C", "target", "." }
    output := buildJarCommand(jarConfig, projectConfig)

    display := func(expected, output []string) string {
        return fmt.Sprintf("Expected: %s, Got: %s", expected, output) 
    }

    if len(expected) != len(output) {
        t.Errorf("Different lengths: %d != %d (%s)", len(expected), len(output), display(expected, output))
    }

    for i, _ := range expected {
        if output[i] != expected[i] {
            t.Errorf("Different value at index %d (%s != %s) Data: %s", i, expected[i], output[i], display(expected, output))
        }
    }
}
