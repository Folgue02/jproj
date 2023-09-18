package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
)

type BuildProjectConfiguration struct {
    Directory string
}

func NewBuildProjectConfiguration(args []string) (*BuildProjectConfiguration, error) {
    parser := argparse.NewParser("build", "Builds/Compiles the current project")
    projectDirectory := parser.String(
        "d",
        "directory", 
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the project's directory" })

    if err := parser.Parse(args); err != nil {
        return nil, err
    }
    
    return &BuildProjectConfiguration {
        Directory: *projectDirectory,
    }, nil
}

func buildProject(args []string) error {
    buildConfig, err := NewBuildProjectConfiguration(args)
    if err != nil {
        return fmt.Errorf("Error: Something wrong with the arguments: %v", err)
    }
    projectConfiguration, err := configuration.LoadConfigurationFromFile(path.Join(buildConfig.Directory, "jproj.json"))
    if err != nil {
        return fmt.Errorf("Cannot load project's configuration due to the following error: %v", err)
    }

    // Get all java files in a slice
    javaFiles, err := utils.GrepFilesByExtension(path.Join(buildConfig.Directory, "src"), "java", utils.GrepFiles)

    if err != nil {
        return fmt.Errorf("Error: Cannot find java source files due to the following error: %v", err)
    } else if len(javaFiles) == 0 {
        return fmt.Errorf("Error: No java source files found in the './src' directory.")
    }

    // Get all jar libs in a slice.
    jarLibs, err := utils.GrepFilesByExtension(path.Join(buildConfig.Directory, projectConfiguration.ProjectLib), "jar", utils.GrepFiles)

    if err != nil {
        return fmt.Errorf("Cannot list jar files in '%s': %v", projectConfiguration.ProjectLib, err)
    }

    // Build
    err = buildSources(buildConfig.Directory, javaFiles, jarLibs, *projectConfiguration)

    if err != nil {
        return fmt.Errorf("Error: Error while compiling with 'javac': %v", err)
    } else {
        log.Println("Done.")
        return nil
    }
}

// Builds up the command, and builds the sources.
// (NOTE: This method doesn't check the validity of none of 
// the arguments passed.)
func buildSources(projectDirectory string,
    javaSources []string,
    jarLibs []string,
    projectConfiguration configuration.Configuration) error {
    // Build the command
    cliCommand := javaSources 
    cliCommand = append(cliCommand, "-d")
    cliCommand = append(cliCommand, path.Join(projectDirectory, projectConfiguration.ProjectTarget))

    if len(jarLibs) > 0 {
        cliCommand = append(cliCommand, "-cp", strings.Join(jarLibs, ":"))
    }

    log.Printf("Building project '%s'...\n", projectConfiguration.ProjectName)

    log.Printf("Build command: --< %v >--\n", cliCommand)
    command := exec.Command("javac", cliCommand...)
    command.Stdout = os.Stdout
    command.Stderr = os.Stderr
    command.Stdin = os.Stdin
    return command.Run()
}
