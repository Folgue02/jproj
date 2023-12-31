package clean

import (
	"fmt"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
)

type CleanConfiguration struct {
    Directory   string
    BinCleaning bool
}

func NewCleanConfiguration(args []string) (*CleanConfiguration, error) {
    parser := argparse.NewParser("clean", "Cleans/Remove generated files.")    
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Project's directory."})

    binCleaningFlag := parser.Flag(
        "b",
        "clean-bin-path",
        &argparse.Options { Required: false, Default: false, Help: "Also clean the binary output directory." })

    if err := parser.Parse(args); err != nil {
        return nil, err
    }

    return &CleanConfiguration {
        Directory: *projectDirectory,
        BinCleaning: *binCleaningFlag,
    }, nil
}

func CleanActionHandler(args []string) error {
    cleanConfig, err := NewCleanConfiguration(args)

    if err != nil {
        return fmt.Errorf("Wrong arguments: %v", err)
    }

    return CleanAction(*cleanConfig)
}

func CleanAction(cleanConfig CleanConfiguration) error {
    projectConfig, err := configuration.LoadConfigurationFromFile(cleanConfig.Directory)

    if err != nil {
        return fmt.Errorf("Couldn't load project's configuration: %v", err)
    }

    targetPath := path.Join(cleanConfig.Directory, projectConfig.ProjectTarget)
    
    err = utils.CleanAllDirectory(targetPath)

    if cleanConfig.BinCleaning {
        binPath := path.Join(cleanConfig.Directory, projectConfig.ProjectBin)

        err = utils.CleanAllDirectory(binPath)

        if err != nil {
            return fmt.Errorf("Error: Cannot clean the bin directory due to the following error: %v\n", err)
        }
    }

    if err != nil {
        return fmt.Errorf("Error: Cannot clean the target directory due to the following error: %v\n", err)
    } else {
        return nil
    }
}
