package run

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
	"github.com/folgue02/jproj/utils/java"
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

    classpaths, err := projectConfiguration.ListJarInLib(runConfig.Directory)

    if err != nil {
        return fmt.Errorf("Couldn't list jars in the lib folder: %v", err)
    }

    includedTargets, err := projectConfiguration.ListIncludedTargetPaths(runConfig.Directory)

    if err != nil {
        return fmt.Errorf("Couldn't list included targets: %v", err)
    }

    classpaths = append(classpaths, includedTargets...)

    javaSources, err := utils.GrepFilesByExtension(filepath.Join(runConfig.Directory, "./src/"), "java", utils.GrepFiles)

    if err != nil {
        return fmt.Errorf("Cannot list sources in dir '%s': %v", filepath.Join(runConfig.Directory, "./src/"), err)
    }
    
	// Running the project requires building it first (if '--no-build' hasn't been specified)
    if !runConfig.NoBuild { 
        javac := java.NewJavacCommand(
            "javac", // TODO: Change it for a configuration parameter
            javaSources,
            classpaths,
            filepath.Join(runConfig.Directory, projectConfiguration.ProjectTarget))
        
        if err := utils.CMD(javac.CompilerPath, javac.Arguments()...); err != nil {
            return fmt.Errorf("Cannot build the project: %v", err)
        }
    }

    javaCommand := java.NewJavaCommand(
        "java", // TODO: Replace with something from a configuration
        append(classpaths, filepath.Join(runConfig.Directory, projectConfiguration.ProjectTarget)),
        mainClass,
        []string {})

    err = utils.CMD(javaCommand.JavaPath, javaCommand.Arguments()...)

	if err != nil {
		return fmt.Errorf("Error: Error while running with 'java': %v", err)
	} else {
		log.Println("Done.")
		return nil
	}
}
