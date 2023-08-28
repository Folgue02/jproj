package actions

import (
	"log"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

func clean(args []string) {
    parser := argparse.NewParser("clean", "Cleans/Remove generated files.")    
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Project's directory."})
    parser.Parse(args)

    projectConfig, err := configuration.LoadConfigurationFromFile(*projectDirectory)

    targetPath := path.Join(*projectDirectory, projectConfig.ProjectTarget)
    
    entries, err := os.ReadDir(targetPath)

    if err != nil {
        log.Printf("Error: Cannot clean the target directory due to the following error: %v\n", err)
        return
    }

    for _, entry := range entries {
        entryPath := path.Join(targetPath, entry.Name())

        err := os.RemoveAll(entryPath)

        if err != nil {
            log.Printf("Error: Cannot remove file/dir while cleaning the target directory: %v\n", err)
            return
        }
    }
}
