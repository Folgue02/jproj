package run

import (
	"fmt"
	"log"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/actions/build"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
)

type RunConfiguration struct {
	MainClass string
	Directory string
    NoBuild   bool
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

    noBuild := parser.Flag(
        "n",
        "no-build",
        &argparse.Options { Required: false, Default: false, Help: "Do not build the project before running it (run the previous compilation of it)." })

	if err := parser.Parse(args); err != nil {
		return nil, err
	}

	return &RunConfiguration{
		MainClass: *projectMainClass,
		Directory: *projectDirectory,
        NoBuild:   *noBuild,
	}, nil
}

func RunProjectActionHandler(args []string) error {
	runConfig, err := NewRunConfiguration(args)

	if err != nil {
		return fmt.Errorf("Error with arguments: %v", err)
	}
    
    return RunProjectAction(*runConfig)
}

func RunProjectAction(runConfig RunConfiguration) error {
	// Running the project requires building it first (if '--no-build' hasn't been specified)
    if !runConfig.NoBuild { 
        if err := build.BuildAction(build.BuildProjectConfiguration { Directory: runConfig.Directory }); err != nil {
            return fmt.Errorf("Cannot build the project: %v", err)
        }
    }

	projectConfiguration, err := configuration.LoadConfigurationFromFile(runConfig.Directory)

	if err != nil {
		return fmt.Errorf("Cannot load configuration due to the following error: %v", err)
	}

    if !projectConfiguration.IsExecutableProject() {
        return fmt.Errorf("The project is not executable (empty/non specified mainclass)");
    }

	// Execute the java command
	mainClass := runConfig.MainClass

	// If the main class is not defined, it will default to the projectConfiguration's
	if runConfig.MainClass == "The one specified by 'jproj.json'" {
		mainClass = projectConfiguration.MainClassPath
	}

    jarLibs, err := projectConfiguration.ListJarInLib(runConfig.Directory)

    if err != nil {
        return fmt.Errorf("Couldn't list jars in the lib folder: %v", err)
    }
    
    javaArgs := buildRunJavaArgs(runConfig, *projectConfiguration, mainClass, jarLibs)
    err = utils.CMD("java", javaArgs...)

	if err != nil {
		return fmt.Errorf("Error: Error while running with 'java': %v", err)
	} else {
		log.Println("Done.")
		return nil
	}
}


// Generates the required arguments for the 'java' binary with the purpose of running the project specified.
func buildRunJavaArgs(runConfig RunConfiguration, projectConfig configuration.Configuration, mainClass string, jarLibs []string) []string {
    jarLibs = append(jarLibs, projectConfig.ProjectTarget)
    classPath := strings.Join(jarLibs, ":")
	javaArgs := []string{"-cp", classPath, mainClass}

    return javaArgs
}

