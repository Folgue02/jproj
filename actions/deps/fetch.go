package deps

import (
	"fmt"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

type FetchActionConfiguration struct {
    Directory string
}

func FetchActionHandler(args []string) error {
    parser := argparse.NewParser("fetch", "Fetches the dependencies listed in the project's configuration.")

    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the project's directory." })

    // TODO: Add 'tolerate-errors' (allows the action to keep going
    // even if some dependencies cannot be fetched) and 'no-overwrite' 
    // (doesn't overwrite already existing jars) flags 

    if err := parser.Parse(args); err != nil {
        return fmt.Errorf("Wrong arguments: %v", err)
    }

    projectConfig, err := configuration.LoadConfigurationFromFile(*projectDirectory)

    if err != nil {
        return fmt.Errorf("Cannot load project's configuration: %v", err)
    }

    return Fetch(FetchActionConfiguration {
        Directory: *projectDirectory,
    }, *projectConfig)

}

func Fetch(fetchConfig FetchActionConfiguration, config configuration.Configuration) error {
	// TODO: Either remove this function entirely (redundant) or
	// give it some functionality.
	if err := config.FetchDependencies(fetchConfig.Directory); err != nil {
		return err
	}
	return nil
}
