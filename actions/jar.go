package actions

import (
	"fmt"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
	"github.com/folgue02/jproj/utils"
)

type JarCommandConfiguration struct {
    Directory string
}

func NewJarCommandConfiguration(args []string) (*JarCommandConfiguration, error) {
    parser := argparse.NewParser("jar", "Creates a JAR file based on the current project.")
    
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Project's directory" })

    if err := parser.Parse(args); err != nil {
        return nil, err
    }

    return &JarCommandConfiguration {
        Directory: *projectDirectory,
    }, nil
}


func CreateJar(args []string) error {
    jarConfig, err := NewJarCommandConfiguration(args)

    if err != nil  {
        return err
    }
    
    projectConfig, err := configuration.LoadConfigurationFromFile(jarConfig.Directory)

    if err != nil {
        return err
    }

    // Save manifest
    err = projectConfig.Manifest.WriteToFile(path.Join(jarConfig.Directory, projectConfig.ProjectTarget, "MANIFEST.mf"))

    if err != nil {
        return fmt.Errorf("Couldn't save manifest file: %v", err)
    }

    cliArgs := buildJarCommand(*jarConfig, *projectConfig)

    err = utils.CMD("jar", cliArgs...)

    if err != nil {
        return fmt.Errorf("Error while executing command: %v", err)
    }

    return nil
}

func buildJarCommand(jarConfig JarCommandConfiguration, projectConfig configuration.Configuration) []string {
    args := []string { "--create", "--file" }
    args = append(args, path.Join(jarConfig.Directory, projectConfig.ProjectBin, projectConfig.ProjectName) + ".jar")
    args = append(args, "--manifest", path.Join(jarConfig.Directory, projectConfig.ProjectTarget, "MANIFEST.mf"))
    if projectConfig.IsExecutableProject() {
        args = append(args, "--main-class", projectConfig.MainClassPath)
    }
    args = append(args, "-C", path.Join(jarConfig.Directory, projectConfig.ProjectTarget), ".")
    return args 
}
