package actions

import (
	"fmt"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

func listproject(args []string) error{
    parser := argparse.NewParser("listproject", "Displays the information about the project")
    projectDirectory := parser.String(
        "d", 
        "directory", 
        &argparse.Options { Required: false, Default: ".", Help: "The directory containing the project and its configuration."},
    )
    parser.Parse(args)

    stat, err := os.Stat(*projectDirectory)

    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("The project's directory specified doesn't exist. ('%s')", *projectDirectory)
        } else {
            return fmt.Errorf("Cannot stat the project's directory ('%s')", *projectDirectory)
        }
    }

    if !stat.IsDir() {
        return fmt.Errorf("The project directory specified ('%s') is not a directory.", *projectDirectory)
    }

    config, err := configuration.LoadConfigurationFromFile(path.Join(*projectDirectory, "jproj.json"))

    if err != nil {
        return fmt.Errorf("Error: Cannot load the project's configuration due to the following error: %v", err)
    } else {
        fmt.Println(config)
    }

    return nil
}

