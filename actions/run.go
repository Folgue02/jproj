package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

type RunConfiguration struct {
	MainClass string
	Directory string
}

func NewRunConfiguration(args []string) (*RunConfiguration, error) {
	parser := argparse.NewParser("run", "Run the specified project.")
	projectMainClass := parser.String(
		"c",
		"mainclass",
		&argparse.Options{Required: false, Default: "The one specified by 'jproj.json'", Help: "Class containing the 'main()' method."},
	)
	projectDirectory := parser.String(
		"d",
		"directory",
		&argparse.Options{Required: false, Default: ".", Help: "Specifies the project's directory."})

	if err := parser.Parse(args); err != nil {
		return nil, err
	}

	return &RunConfiguration{
		MainClass: *projectMainClass,
		Directory: *projectDirectory,
	}, nil
}

func runProject(args []string) error {

	runConfig, err := NewRunConfiguration(args)

	if err != nil {
		return fmt.Errorf("Error with arguments: %v", err)
	}

	// Running the project requires building it first
	if err := buildProject([]string{"build", "-d", runConfig.Directory}); err != nil {
		return fmt.Errorf("Error: Cannot build the project: %v", err)
	}

	projectConfiguration, err := configuration.LoadConfigurationFromFile(runConfig.Directory)

	if err != nil {
		return fmt.Errorf("Error: Cannot load configuration due to the following error: %v", err)
	}

	// Execute the java command
	mainClass := runConfig.MainClass

	// If the main class is not defined, it will default to the projectConfiguration's
	if runConfig.MainClass == "The one specified by 'jproj.json'" {
		mainClass = projectConfiguration.MainClassPath
	}

	javaArgs := []string{"-cp", projectConfiguration.ProjectTarget, mainClass}
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
