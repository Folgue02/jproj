package deps

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
)

type CleanActionConfiguration struct {
    Directory string
    IgnoreLocal bool
}

func CleanActionHandler(args []string) error {
    parser := argparse.NewParser("clean", "Cleans the project's lib directory.")
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the project's directory." })

    // TODO: Finish the logic of this flag
    ignoreLocal := parser.Flag(
        "i",
        "ignore-local",
        &argparse.Options { Required: false, Default: false, Help: "Prevents this action from removing local dependencies in the lib directory that are not listed in the project's configuration." })

    if err := parser.Parse(args); err != nil {
        return fmt.Errorf("Wrong arguments: %v", err)
    }

    projectConfig, err := configuration.LoadConfigurationFromFile(*projectDirectory)

    if err != nil {
        return fmt.Errorf("Cannot load project's configuration: %v", err)
    }
    
    return CleanDependenciesAction(CleanActionConfiguration { 
        Directory: *projectDirectory,
        IgnoreLocal: *ignoreLocal,
    }, projectConfig)
}

func CleanDependenciesAction(cleanConfig CleanActionConfiguration, configuration *configuration.Configuration) error {
	projectDirectory, err  := filepath.Abs(cleanConfig.Directory)

	if err != nil {
		return fmt.Errorf("The path specified doesn't exist: %v", err)
	}

    utils.CleanDirectory(path.Join(projectDirectory, configuration.ProjectLib), func(targetPath string) bool {
        if strings.HasSuffix(targetPath, ".jar") {
            log.Printf("Removing .jar ('%s')...\n", targetPath)
            return true
        }

        return false
    })
    /*
	filepath.Walk(path.Join(projectDirectory, configuration.ProjectLib), func(targetPath string, fileInfo fs.FileInfo, err error) error {
		jarPath, _ := filepath.Abs(targetPath)
        parentPath, _ := filepath.Abs("..")

		// Ignore the root path
		if jarPath == projectDirectory || jarPath == parentPath {
			return nil
		}

		log.Printf("Removing .jar ('%s')...\n", jarPath)
		return os.RemoveAll(targetPath)
	})*/

	return nil
}
