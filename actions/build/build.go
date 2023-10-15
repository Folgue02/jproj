package build

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
	"github.com/folgue02/jproj/utils/java"
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

func BuildActionHandler(args []string) error {
    buildConfig, err := NewBuildProjectConfiguration(args)

    if err != nil {
        return err
    }

    return BuildAction(*buildConfig)
}

func BuildAction(buildConfig BuildProjectConfiguration) error {
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
    javacCommand := java.NewJavacCommand(
        "javac", // TODO: Change for an actual compiler path given by configuration
        javaFiles,
        jarLibs,
        filepath.Join(buildConfig.Directory, projectConfiguration.ProjectTarget))

    err = utils.CMD(javacCommand.CompilerPath, javacCommand.Arguments()...)

    if err != nil {
        return fmt.Errorf("Error: Error while compiling with 'javac': %v", err)
    } else {
        log.Println("Done.")
        return nil
    }
}
