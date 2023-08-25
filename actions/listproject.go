package actions

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

func listproject(args []string) {
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
            log.Printf("The project's directory specified doesn't exist. ('%s')\n", *projectDirectory)
        } else {
            log.Printf("Cannot stat the project's directory ('%s')\n", *projectDirectory)
        }
        return
    }

    if !stat.IsDir() {
        log.Printf("The project directory specified ('%s') is not a directory.\n", *projectDirectory)
        return
    }

    config, err := configuration.LoadConfigurationFromFile(path.Join(*projectDirectory, "jproj.json"))

    if err != nil {
        log.Printf("Error: Cannot load the project's configuration due to the following error: %v\n", err)
        return
    } else {
        fmt.Println(config)
    }
}

