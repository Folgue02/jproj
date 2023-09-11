package actions

import (
	"fmt"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

type ListProjectConfiguration struct {
    Directory string
}

func NewListProjectConfiguration(args []string) (*ListProjectConfiguration, error) {
    parser := argparse.NewParser("listproject", "Displays the information about the project")
    projectDirectory := parser.String(
        "d", 
        "directory", 
        &argparse.Options { Required: false, Default: ".", Help: "The directory containing the project and its configuration."},
    )

    if err := parser.Parse(args); err != nil {
        return nil, err
    }

    return &ListProjectConfiguration {
        Directory: *projectDirectory,
    }, nil
}

func listproject(args []string) error{
    listConfig, err := NewListProjectConfiguration(args)

    if err != nil {
        return fmt.Errorf("Wrong arguments: %v", err)
    }

    stat, err := os.Stat(listConfig.Directory)

    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("The project's directory specified doesn't exist. ('%s')", listConfig.Directory)
        } else {
            return fmt.Errorf("Cannot stat the project's directory ('%s')", listConfig.Directory)
        }
    }

    if !stat.IsDir() {
        return fmt.Errorf("The project directory specified ('%s') is not a directory.", listConfig.Directory)
    }

    config, err := configuration.LoadConfigurationFromFile(path.Join(listConfig.Directory, "jproj.json"))

    if err != nil {
        return fmt.Errorf("Error: Cannot load the project's configuration due to the following error: %v", err)
    } else {
        fmt.Println(config)
    }

    return nil
}

