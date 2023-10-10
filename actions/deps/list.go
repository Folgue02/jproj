package deps

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

type ListActionConfiguration struct {
    Directory string
    LocalOnly bool
}

func ListActionHandler(args []string) error {
    parser := argparse.NewParser("list", "Lists all of the project's dependencies.")

    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the project's directory." })

    localOnly := parser.Flag(
        "l",
        "local-only",
        &argparse.Options { Required: false, Default: false, Help: "Lists all local dependencies that are not libFiles in the project's configuration (or in other words, jar's manually put in the lib directory of the project)." })

    if err := parser.Parse(args); err != nil {
        return fmt.Errorf("Wrong arguments: %v", err)
    }

    return ListAction(ListActionConfiguration {
        Directory: *projectDirectory,
        LocalOnly: *localOnly,
    })
}

func ListAction(listConfig ListActionConfiguration) error {
    projectConfig, err := configuration.LoadConfigurationFromFile(listConfig.Directory)
    if err != nil {
        return fmt.Errorf("Cannot load project's configuration: %v", err)
    }

    if listConfig.LocalOnly {
        jarsInLib, err := projectConfig.ListJarInLib(listConfig.Directory)

        if err != nil {
            return fmt.Errorf("Cannot project's lib directory: %v", err)
        }
        
        localDepsCount := 0
        for _, jarName := range jarsInLib {
            stat, err := os.Stat(jarName)

            if err != nil {
                return fmt.Errorf("Cannot stat jar path '%s': %v", jarName, err)
            }

            log.Println("Local dependencies:")
            if !stat.IsDir() && filepath.Ext(jarName) == ".jar" && 
                    !projectConfig.IsJarInDependencies(jarName) {
                localDepsCount += 1
                log.Printf("   [%d]:  %s\n", localDepsCount, jarName)
            }
        }

        if localDepsCount == 0 {
            log.Printf("  There are no local dependencies.")
        }
        
    } else {
        if len(projectConfig.Dependencies) == 0 {
            log.Println("This project doesn't have any libFiles dependencies.")
            return nil
        }

        log.Println("Project's libFiles dependencies: ")
        for i, dep := range projectConfig.Dependencies {
            log.Printf("   [%d]: %s\n", i + 1, dep)
        }

    }
    return nil
}

