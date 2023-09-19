package actions

import (
	"fmt"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/actions/build"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
)

type JarCommandConfiguration struct {
    Directory  string
    Static     bool
    OutputPath string
}

func NewJarCommandConfiguration(args []string) (*JarCommandConfiguration, error) {
    parser := argparse.NewParser("jar", "Creates a JAR file based on the current project.")
    
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Project's directory" })

    staticLinking := parser.Flag(
        "s",
        "static",
        &argparse.Options { Required: false, Default: false, Help: "Statically links the jar to the dependencies (includes all dependencies in the jar)."},
    )

    outputPath := parser.String(
        "o",
        "output",
        &argparse.Options { Required: false, Default: "", Help: "Specify the place to output the jar to." })

    if err := parser.Parse(args); err != nil {
        return nil, err
    }

    return &JarCommandConfiguration {
        Directory:  *projectDirectory,
        Static:     *staticLinking,
        OutputPath: *outputPath,
    }, nil
}

func CreateJarActionHandler(args []string) error {
    jarConfig, err := NewJarCommandConfiguration(args)

    if err != nil  {
        return err
    }
    
    return CreateJarAction(*jarConfig)
}

func CreateJarAction(jarConfig JarCommandConfiguration) error {
    projectConfig, err := configuration.LoadConfigurationFromFile(jarConfig.Directory)

    if err != nil {
        return err
    }

    if err = build.BuildAction(build.BuildProjectConfiguration { Directory: jarConfig.Directory });
        err != nil {
            return fmt.Errorf("Error while compiling sources pre-jar: %v", err)
    }


    // Save manifest
    err = projectConfig.Manifest.WriteToFile(path.Join(jarConfig.Directory, projectConfig.ProjectTarget, "MANIFEST.mf"))

    if err != nil {
        return fmt.Errorf("Couldn't save manifest file: %v", err)
    }

    cliArgs := buildJarCommand(jarConfig, *projectConfig)

    err = utils.CMD("jar", cliArgs...)

    if err != nil {
        return fmt.Errorf("Error while executing command: %v", err)
    }

    return nil
}

func buildJarCommand(jarConfig JarCommandConfiguration, projectConfig configuration.Configuration) []string {
    args := []string { "--create", "--file" }
    
    if jarConfig.OutputPath == "" {
        args = append(args, path.Join(jarConfig.Directory, projectConfig.ProjectBin, projectConfig.ProjectName) + ".jar")
    } else {
        args = append(args, jarConfig.OutputPath)
    }

    args = append(args, "--manifest", path.Join(jarConfig.Directory, projectConfig.ProjectTarget, "MANIFEST.mf"))
    if projectConfig.IsExecutableProject() {
        args = append(args, "--main-class", projectConfig.MainClassPath)
    }

    args = append(args, "-C", path.Join(jarConfig.Directory, projectConfig.ProjectTarget), ".")

    if jarConfig.Static {
        args = append(args, "-C", path.Join(jarConfig.Directory, projectConfig.ProjectLib), ".")
    }
    return args 
}
