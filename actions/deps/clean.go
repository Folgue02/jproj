package deps

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/folgue02/jproj/configuration"
)

func CleanDependencies(projectDirectory string, configuration *configuration.Configuration) error {
	projectDirectory, err  := filepath.Abs(projectDirectory)

	if err != nil {
		return fmt.Errorf("The path specified doesn't exist: %v", err)
	}
	filepath.Walk(path.Join(projectDirectory, configuration.ProjectLib), func(targetPath string, fileInfo fs.FileInfo, err error) error {
		jarPath, _ := filepath.Abs(targetPath)

		// Ignore the root path
		if jarPath == projectDirectory {
			return nil
		}

		log.Printf("Removing .jar ('%s')...\n", jarPath)
		return os.RemoveAll(targetPath)
	})

	return nil
}