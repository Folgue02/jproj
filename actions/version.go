package actions

import (
	"fmt"
	"runtime"
)

const VERSION = 0.25

func versionAction(args []string) error {
    fmt.Printf("JProj Current version: %.2f\n", VERSION)
    fmt.Printf("Running on %s_%s\n", runtime.GOOS, runtime.GOARCH)
    return nil
}
