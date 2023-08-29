package deps

import (
	"fmt"

	"github.com/folgue02/jproj/configuration"
)

func AddDependency(projectDirectory string, projectConfiguration *configuration.Configuration, depList []string) error {
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