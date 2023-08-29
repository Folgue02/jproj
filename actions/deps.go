package actions

import (
	"fmt"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/actions/deps"
	"github.com/folgue02/jproj/configuration"
)

func manageDependencies(args []string) error {
	parser := argparse.NewParser("deps", "Manages dependencies")
	action := parser.String(
		"a",
		"action",
		&argparse.Options { Required: false, Default: "fetch", Help: "Action to be done on the dependencies."})
	projectDirectory := parser.String(
		"d",
		"directory",
		&argparse.Options { Required: false, Default: ".", Help: "Directory where the project is located."},
	)

	dependencyList := parser.StringList(
		"e",
		"dependency-list",
		&argparse.Options { Required: false, Default: []string {}, Help: "List of dependencies."},
	)
	if err := parser.Parse(args);err != nil {
		return fmt.Errorf("Wrong arguments: %v", err)
	}

	projectConfiguration, err := configuration.LoadConfigurationFromFile(*projectDirectory)

	if err != nil {
		return fmt.Errorf("Cannot load configuration due to the following error: %v", err)
	}

	switch *action {
	case "fetch":
		if err := deps.Fetch(*projectDirectory, *projectConfiguration); err != nil {
			return err
		}
	case "clean":
	case "add":
		return deps.AddDependency(*projectDirectory, projectConfiguration, *dependencyList)
	default:
		return fmt.Errorf("Action not found: '%s'", *action)
	}

	return nil
}