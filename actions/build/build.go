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
    Directory     string
    BuildIncludes bool
}

func NewBuildProjectConfiguration(args []string) (*BuildProjectConfiguration, error) {
    parser := argparse.NewParser("build", "Builds/Compiles the current project")
    projectDirectory := parser.String(
        "d",
        "directory", 
        &argparse.Options { Required: false, Default: ".", Help: "Specifies the project's directory" })

    buildIncludes := parser.Flag(
        "i",
        "build-includes",
        &argparse.Options { Required: false, Default: false, Help: "Builds the included projects specified in the project's configuration before building the current project." })

    if err := parser.Parse(args); err != nil {
        return nil, err
    }
    
    return &BuildProjectConfiguration {
        Directory:     *projectDirectory,
        BuildIncludes: *buildIncludes,
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

    if buildConfig.BuildIncludes && len(projectConfiguration.IncludedProjects) > 0 {
        log.Println("Building included projects...")
        if err = buildIncludedProjects(buildConfig.Directory, *projectConfiguration); err != nil {
            return fmt.Errorf("Error while building included projects: %v\n", err)
        }
    }

    // Get all java files in a slice
    javaFiles, err := utils.GrepFilesByExtension(path.Join(buildConfig.Directory, "src"), "java", utils.GrepFiles)

    if err != nil {
        return fmt.Errorf("Cannot find java source files due to the following error: %v", err)
    } else if len(javaFiles) == 0 {
        return fmt.Errorf("No java source files found in the './src' directory.")
    }

    // ------------------------
    // Get the project's classpaths
    // ------------------------

    // Get the paths to the jars in ./lib
    classpaths, err := utils.GrepFilesByExtension(path.Join(buildConfig.Directory, projectConfiguration.ProjectLib), "jar", utils.GrepFiles)

    if err != nil {
        return fmt.Errorf("Cannot list jar files in '%s': %v", projectConfiguration.ProjectLib, err)
    }

    // Get the paths of the target folders from the 
    // included projects
    includedTargets, err := projectConfiguration.ListIncludedTargetPaths(buildConfig.Directory)

    if err != nil {
        return fmt.Errorf("Cannot read included projects: %v", err)
    }

    classpaths = append(classpaths, includedTargets...)

    // Build
    javacCommand := java.NewJavacCommand(
        "javac", // TODO: Change for an actual compiler path given by configuration
        javaFiles,
        classpaths,
        filepath.Join(buildConfig.Directory, projectConfiguration.ProjectTarget))

    err = utils.CMD(javacCommand.CompilerPath, javacCommand.Arguments()...)

    if err != nil {
        return fmt.Errorf("Error: Error while compiling with 'javac': %v", err)
    } else {
        log.Println("Done.")
        return nil
    }
}

func buildIncludedProjects(mainProjectPath string, mainProjectConfig configuration.Configuration) error {
    for i, includedPath := range mainProjectConfig.IncludedProjects {
        includedPath = filepath.Join(mainProjectPath, includedPath)
        log.Printf("[%d]: Building project located in path '%s'...\n", i + 1, includedPath)

        includedConfiguration, err := configuration.LoadConfigurationFromFile(includedPath)

        if err != nil {
            return fmt.Errorf("Cannot build included project '%s': %v\n", includedPath, err)
        }

        javaSources, err := utils.GrepFilesByExtension(filepath.Join(includedPath, "./src"), "java", utils.GrepFiles)
        
        if err != nil {
            return fmt.Errorf("Cannot list java sources in included project '%s': %v\n", includedPath, err)
        }

        jarLibs, err := includedConfiguration.ListJarInLib(includedPath)

        if err != nil {
            return fmt.Errorf("Cannot list jars in included project '%s': %v\n", includedPath, err)
        }
        javac := java.NewJavacCommand("javac", javaSources, jarLibs, filepath.Join(includedPath, includedConfiguration.ProjectTarget))

        err = utils.CMD(javac.CompilerPath, javac.Arguments()...)

        if err != nil {
            return fmt.Errorf("Couldn't build included project '%s': %v", includedPath, err)
        }
    }
    log.Println("Done building included projects.")

    return nil
}
