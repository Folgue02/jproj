package deps

import (
	"fmt"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

type AddActionConfiguration struct {
    Directory string
    Dependencies []string
}

func AddActionHandler(args []string) error {
    parser := argparse.NewParser("add", "Adds dependencies to the configuration of the project.")
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the directory of the project."})
    dependencies := parser.List(
        "e",
        "dependency-list",
        &argparse.Options { Required: true, Help: "Dependencies to add." })

    if err := parser.Parse(args); err != nil {
        return fmt.Errorf("Wrong arguments: %v", err);
    }

    projectConfiguration, err := configuration.LoadConfigurationFromFile(*projectDirectory)

    if err != nil {
        return fmt.Errorf("Cannot load project's configuration: %v", err)
    }

    return addDependency(*projectDirectory, projectConfiguration, *dependencies)
}

func addDependency(projectDirectory string, projectConfiguration *configuration.Configuration, depList []string) error {
	for _, depString := range depList {
		dependency, err := configuration.NewDependencyFromString(depString)

		if err != nil {
			return err
		}

		if projectConfiguration.DependencyExists(dependency.Name) {
			return fmt.Errorf("Cannot add dependency '%s', there is already a dependency with that name ('%s')", dependency, dependency.Name)
		}

		projectConfiguration.Dependencies = append(projectConfiguration.Dependencies, *dependency)
	}

	if err := projectConfiguration.SaveConfiguration(projectDirectory); err != nil {
		return fmt.Errorf("Cannot save configuration: %v", err)
	} else {
		return nil
	}
}
