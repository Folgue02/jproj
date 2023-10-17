package actions

import (
    "fmt"
    "path"
    "path/filepath"

    "github.com/akamensky/argparse"
    "github.com/folgue02/jproj/actions/build"
    "github.com/folgue02/jproj/configuration"
    "github.com/folgue02/jproj/utils"
    "github.com/folgue02/jproj/utils/java"
)

type JarCommandConfiguration struct {
    Directory  string
    Static     bool
    NoBuild    bool
    OutputPath string
}

func NewJarCommandConfiguration(args []string) (*JarCommandConfiguration, error) {
    parser := argparse.NewParser("jar", "Creates a JAR file based on the current project.")

    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options{Required: false, Default: ".", Help: "Project's directory"})

    staticLinking := parser.Flag(
        "s",
        "static",
        &argparse.Options{Required: false, Default: false, Help: "Statically links the jar to the dependencies (includes all dependencies in the jar)."},
    )

    outputPath := parser.String(
        "o",
        "output",
        &argparse.Options{Required: false, Default: "", Help: "Specify the place to output the jar to."})

    noBuild := parser.Flag(
        "n",
        "no-build",
        &argparse.Options{Required: false, Default: false, Help: "Do not build the project before creating the jar (use the previously compiled classes of the project)."})

    if err := parser.Parse(args); err != nil {
        return nil, err
    }

    return &JarCommandConfiguration{
        Directory:  *projectDirectory,
        NoBuild:    *noBuild,
        Static:     *staticLinking,
        OutputPath: *outputPath,
    }, nil
}

func CreateJarActionHandler(args []string) error {
    jarConfig, err := NewJarCommandConfiguration(args)

    if err != nil {
        return err
    }

    return CreateJarAction(*jarConfig)
}

func CreateJarAction(jarConfig JarCommandConfiguration) error {
    projectConfig, err := configuration.LoadConfigurationFromFile(jarConfig.Directory)
    outputManifestPath := path.Join(jarConfig.Directory, projectConfig.ProjectBin, "MANIFEST.mf")

    if err != nil {
        return err
    }

    if !jarConfig.NoBuild {
        if err = build.BuildAction(build.BuildProjectConfiguration{Directory: jarConfig.Directory}); err != nil {
            return fmt.Errorf("Error while compiling sources pre-jar: %v", err)
        }
    }


    // Save manifest
    err = projectConfig.Manifest.WriteToFile(outputManifestPath)

    if err != nil {
        return fmt.Errorf("Couldn't save manifest file: %v", err)
    }

    jarCommand := java.NewJarCommand(
        "jar",
        filepath.Join(jarConfig.Directory, projectConfig.ProjectBin, projectConfig.OutputJarName()),
        [][]string{{filepath.Join(jarConfig.Directory, projectConfig.ProjectTarget)}},
        projectConfig.MainClassPath)

    jarCommand.ManifestFile = outputManifestPath

    if jarConfig.Static {
        includedTargets, err := projectConfig.ListIncludedTargetPaths(jarConfig.Directory)

        if err != nil {
            return fmt.Errorf("Cannot list included target paths: %v", err)
        }

        for _, includedTarget := range includedTargets {
            jarCommand.Sources = append(jarCommand.Sources, []string{ includedTarget })
        }
    }
    err = utils.CMD(jarCommand.JarPath, jarCommand.Arguments()...)

    if err != nil {
        return fmt.Errorf("Error while executing command: %v", err)
    }

    return nil
}
