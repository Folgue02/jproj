package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/folgue02/jproj/configuration"
    "github.com/akamensky/argparse"
)

func createproject(args []string) {
    parser := argparse.NewParser("createproj", "Creates a new project")
    projectName := parser.String("n", "name", &argparse.Options { Required: true, Help: "Name of the new project"})
    projectBaseDirectory := parser.String("b", "basedirectory", &argparse.Options { Required: false, Default: ".", Help: "Directory in where the new project should be created." })
    err := parser.Parse(args)

    if err != nil {
        log.Printf("Error on argument parsing: %v\n", err)
        return
    }

    baseDirectoryStat, err := os.Stat(*projectBaseDirectory)
    
    if err != nil {
        log.Printf("Error: Cannot stat the base directory specified ('%s'): %v\n", *projectBaseDirectory, err)
        return
    }

    if !baseDirectoryStat.IsDir() {
        log.Printf("Error: The base directory specified ('%s') isn't a directory.\n", *projectBaseDirectory)
        return
    }

    fmt.Printf("Project name: %s, Base directory: %s\n", *projectName, *projectBaseDirectory)
    projPath, configf, err := setupProject(*projectName, *projectBaseDirectory)

    if err != nil {
        log.Printf("Error: Cannot setup the project's directory due to the following error: '%v'\n", err)
        return
    }

    projectConfiguration := configuration.Configuration { 
        ProjectName: *projectName,
        ProjectTarget: "./target",
    }

    jsonConfiguration, _ := json.Marshal(projectConfiguration)

    if _, err := configf.Write(jsonConfiguration); err != nil {
        log.Printf("Error: Cannot save configuration into the configuration file in '%s' due to the following error: %v\n", *projPath, err)
        return
    }

    log.Printf("Project '%s' created.\n", *projectName)
}

// Creates the project's directory, as well as its configuration file.
// Returns:
// 1. The project's path
// 2. A `File` object of the configuration file within the project.
func setupProject(projectName string, directory string) (*string, *os.File, error) {
    projectPath := path.Join(directory, projectName)
    projectConfigurationPath := path.Join(projectPath, "jproj.json")
    // TODO: Control errors from the other creations
    err := os.Mkdir(projectPath, 0750) // rwxr-x---
    os.Mkdir(path.Join(projectPath, "src"), 0750)
    os.Mkdir(path.Join(projectPath, "target"), 0750)

    if err != nil {
        return nil, nil, fmt.Errorf("Error: Cannot create directory '%s' due to the following error: %s\n", projectName, err)
    }

    configf, err := os.Create(projectConfigurationPath)

    if err != nil {
        return nil, nil, fmt.Errorf("Error: Cannot create config file in '%s' due to the following error: %s\n", projectConfigurationPath, err)
    }

    return &projectPath, configf, nil
}
