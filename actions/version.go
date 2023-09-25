package actions

import (
	"fmt"
	"runtime"

    "github.com/folgue02/jproj/version"
)

func VersionAction(args []string) error {
    fmt.Printf("JProj Current version: %s\n", version.GetJprojVersion().String())
    fmt.Printf("Running on %s_%s\n", runtime.GOOS, runtime.GOARCH)
    return nil
}
