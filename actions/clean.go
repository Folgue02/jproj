package actions

import (
	"fmt"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

type CleanConfiguration struct {
    Directory string
}

func NewCleanConfiguration(args []string) (*CleanConfiguration, error) {
    parser := argparse.NewParser("clean", "Cleans/Remove generated files.")    
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Project's directory."})

    if err := parser.Parse(args); err != nil {
        return nil, err
    }

    return &CleanConfiguration {
        Directory: *projectDirectory,
    }, nil
}

func clean(args []string) error {
    cleanConfig, err := NewCleanConfiguration(args)

    if err != nil {
        return fmt.Errorf("Wrong arguments: %v", err)
    }

    projectConfig, err := configuration.LoadConfigurationFromFile(cleanConfig.Directory)

    targetPath := path.Join(cleanConfig.Directory, projectConfig.ProjectTarget)
    
    entries, err := os.ReadDir(targetPath)

    if err != nil {
        return fmt.Errorf("Error: Cannot clean the target directory due to the following error: %v\n", err)
    }

    for _, entry := range entries {
        entryPath := path.Join(targetPath, entry.Name())

        err := os.RemoveAll(entryPath)

        if err != nil {
            return fmt.Errorf("Error: Cannot remove file/dir while cleaning the target directory: %v\n", err)
        }
    }
    return nil
}
