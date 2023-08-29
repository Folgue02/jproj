package deps

import "github.com/folgue02/jproj/configuration"

func Fetch(directory string, config configuration.Configuration) error {
	// TODO: Either remove this function entirely (redundant) or
	// give it some functionality.
	if err := config.FetchDependencies(directory); err != nil {
		return err
	}
	return nil
}