package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
)

func buildProject(args []string) error {
    parser := argparse.NewParser("build", "Builds/Compiles the current project")
    projectDirectory := parser.String(
        "d",
        "directory", 
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the project's directory" })
    if err := parser.Parse(args); err != nil {
        return fmt.Errorf("Error: Something wrong with the arguments: %v", err)
    }

    projectDirectoryStat, err := os.Stat(*projectDirectory)

    if err != nil {
        return fmt.Errorf("Error: Cannot stat project's directory due to the following error: %v", err)
    } else if !projectDirectoryStat.IsDir() {
        return fmt.Errorf("Error: The path to the project's directory specified ('%s') doesn't point to a directory.", *projectDirectory)
    }

    javaFiles, err := utils.GrepFilesByExtension(path.Join(*projectDirectory, "src"), "java", utils.GrepFiles)

    if err != nil {
        return fmt.Errorf("Error: Cannot find java source files due to the following error: %v", err)
    } else if len(javaFiles) == 0 {
        return fmt.Errorf("Error: No java source files found in the './src' directory.")
    }
    projectConfiguration, err := configuration.LoadConfigurationFromFile(path.Join(*projectDirectory, "jproj.json"))
    if err != nil {
        return fmt.Errorf("Cannot load project's configuration due to the following error: %v", err)
    }

    // Build
    err = buildSources(*projectDirectory, javaFiles, *projectConfiguration)

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
    projectConfiguration configuration.Configuration) error {
    // Build the command
    cliCommand := javaSources 
    cliCommand = append(cliCommand, "-d")
    cliCommand = append(cliCommand, path.Join(projectDirectory, "./target/"))

    log.Printf("Building project '%s'...\n", projectConfiguration.ProjectName)

    log.Printf("Build command: --< %v >--\n", cliCommand)
    command := exec.Command("javac", cliCommand...)
    command.Stdout = os.Stdout
    command.Stderr = os.Stderr
    command.Stdin = os.Stdin
    return command.Run()
}
