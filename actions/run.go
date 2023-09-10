package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

func runProject(args []string) error {
    parser := argparse.NewParser("run", "Run the specified project.")
    projectMainClass := parser.String(
        "c",
        "mainclass",
        &argparse.Options { Required: false, Default: "The one specified by 'jproj.json'", Help: "Class containing the 'main()' method." },
    )
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the project's directory." })
    if err := parser.Parse(args); err != nil {
        return fmt.Errorf("Error: Wrong arguments: %v", err)
    }

    // Running the project requires building it first
    if err := buildProject([]string { "build", "-d", *projectDirectory }); err != nil {
        return fmt.Errorf("Error: Cannot build the project: %v", err)
    }

    projectConfiguration, err := configuration.LoadConfigurationFromFile(*projectDirectory)

    if err != nil {
        return fmt.Errorf("Error: Cannot load configuration due to the following error: %v", err)
    }

    // Execute the java command
    mainClass := *projectMainClass

    // If the main class is not defined, it will default to the projectConfiguration's
    if *projectMainClass == "The one specified by 'jproj.json'" {
        mainClass = projectConfiguration.MainClassPath
    }
    
    javaArgs := []string { "-cp", projectConfiguration.ProjectTarget, mainClass }
    log.Printf("---< Executing 'java' with args %v >---\n", javaArgs)
    cmd := exec.Command("java", javaArgs...)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    err = cmd.Run()

    if err != nil {
        return fmt.Errorf("Error: Error while running with 'java': %v", err)
    } else {
        log.Println("Done.")
        return nil
    }

}