package actions

import (
	"fmt"
	"runtime"
)

const VERSION = 0.2

func version(args []string) error {
    fmt.Printf("JProj Current version: %.1f\n", VERSION)
    fmt.Printf("Running on %s_%s\n", runtime.GOOS, runtime.GOARCH)
    return nil
}
