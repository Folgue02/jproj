package actions

import (
	"fmt"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/actions/deps"
	"github.com/folgue02/jproj/configuration"
)

type ManageDependenciesConfiguration struct {
    ActionName string
    Directory string
    DependencyList []string
}

func NewManageDependenciesConfiguration(args []string) (*ManageDependenciesConfiguration, error) {
	parser := argparse.NewParser("deps", "Manages dependencies")
	action := parser.String(
		"a",
		"action",
		&argparse.Options{Required: false, Default: "fetch", Help: "Action to be done on the dependencies.(fetch, clean, add)"})
	projectDirectory := parser.String(
		"d",
		"directory",
		&argparse.Options{Required: false, Default: ".", Help: "Directory where the project is located."},
	)

	dependencyList := parser.StringList(
		"e",
		"dependency-list",
		&argparse.Options{Required: false, Default: []string{}, Help: "List of dependencies."},
	)
	if err := parser.Parse(args); err != nil {
		return nil, err
	}

    return &ManageDependenciesConfiguration {
        ActionName: *action,
        Directory: *projectDirectory,
        DependencyList: *dependencyList,
    }, nil
}

func manageDependencies(args []string) error {
    depConfig, err := NewManageDependenciesConfiguration(args)

	if err != nil {
		return fmt.Errorf("Wrong arguments: %v", err)
	}

	projectConfiguration, err := configuration.LoadConfigurationFromFile(depConfig.Directory)

	if err != nil {
		return fmt.Errorf("Cannot load configuration due to the following error: %v", err)
	}

	switch depConfig.ActionName {
	case "fetch":
		if err := deps.Fetch(depConfig.Directory, *projectConfiguration); err != nil {
			return err
		}
	case "clean":
		return deps.CleanDependencies(depConfig.Directory, projectConfiguration)
	case "add":
		return deps.AddDependency(depConfig.Directory, projectConfiguration, depConfig.DependencyList)
	default:
		return fmt.Errorf("Action not found: '%s'", depConfig.ActionName)
	}

	return nil
}
