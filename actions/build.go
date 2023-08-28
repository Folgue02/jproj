package actions

import (
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
)

func buildProject(args []string) {
    parser := argparse.NewParser("build", "Builds/Compiles the current project")
    projectDirectory := parser.String(
        "d",
        "directory", 
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the project's directory" })
    parser.Parse(args)

    projectDirectoryStat, err := os.Stat(*projectDirectory)

    if err != nil {
        log.Printf("Error: Cannot stat project's directory due to the following error: %v\n", err)
        return
    } 

    if !projectDirectoryStat.IsDir() {
        log.Printf("Error: The path to the project's directory specified ('%s') doesn't point to a directory.\n", *projectDirectory)
        return
    }

    javaFiles, err := utils.GrepFilesByExtension(path.Join(*projectDirectory, "src"), "java", utils.GrepFiles)

    if err != nil {
        log.Printf("Error: Cannot find java source files due to the following error: %v\n", err)
        return
    } else if len(javaFiles) == 0 {
        log.Printf("Error: No java source files found in the './src' directory.\n")
        return
    }
    projectConfiguration, err := configuration.LoadConfigurationFromFile(path.Join(*projectDirectory, "jproj.json"))
    if err != nil {
        log.Printf("Cannot load project's configuration due to the following error: %v\n", err)
        return
    }

    // Build
    err = buildSources(*projectDirectory, javaFiles, *projectConfiguration)

    if err != nil {
        log.Printf("Error: Error while compiling with 'javac': %v\n", err)
    } else {
        log.Println("Done.")
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
