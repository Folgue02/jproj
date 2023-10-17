package validate

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

type ValidateConfiguration struct {
    Directory       string
    ValidateJdk     bool
    ValidateProject bool
}

func NewValidateConfiguration(args []string) (*ValidateConfiguration, error) {
    parser := argparse.NewParser("validate", "Validate elements such as JDK, project...")
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Project's directory." })

    validateJdk := parser.Flag(
        "j",
        "validate-jdk",
        &argparse.Options { Required: false, Default: false, Help: "Checks the validity of the JDK." })

    validateProject := parser.Flag(
        "p",
        "validate-project",
        &argparse.Options { Required: false, Default: false, Help: "Checks the validity of the project located at the specified directory." })

    if err := parser.Parse(args); err != nil {
        return nil, err
    }

    return &ValidateConfiguration {
        Directory: *projectDirectory,
        ValidateJdk: *validateJdk,
        ValidateProject: *validateProject,
    }, nil
}

func ValidateActionHandler(args []string) error {
    vConfig, err := NewValidateConfiguration(args)

    if err != nil {
        return err
    }  

    return ValidateAction(*vConfig)
}

func ValidateAction(vConfig ValidateConfiguration) error {
    if vConfig.ValidateJdk {
        validateJdk(vConfig)
    }

    if vConfig.ValidateProject {
        if err := validateProject(vConfig); err != nil {
            log.Println("==> Invalid project, check the logs below to find out why.")
            return err
        }

        log.Println("Project validated (valid).")
    }

    return  nil
}

func validateProject(vConfig ValidateConfiguration) error {
    f, err := os.Stat(vConfig.Directory)

    if err != nil {
        return fmt.Errorf("Cannot read directory: %s", err)
    } else if !f.IsDir() {
        return fmt.Errorf("The directory doesn't exist.\n")
    }

    pConfig, err := configuration.LoadConfigurationFromFile(vConfig.Directory)
    
    if err != nil {
        return fmt.Errorf("Cannot load project configuration: %s", err)
    }

    err = pConfig.Validate(vConfig.Directory)

    if err != nil {
        return err 
    }

    if err = pConfig.ValidateIncluded(vConfig.Directory); err != nil {
        return err
    }
    
    return nil
}

func validateJdk(vConfig ValidateConfiguration) {
    requiredExe := []string { 
        "javac",
        "java",
        "jar",
    }

    found := 0
    for _, exe := range requiredExe {
        exeName := exe

        // I <3 Windows
        if runtime.GOOS == "windows" && !strings.HasSuffix(exeName, ".exe") {
            exeName += fmt.Sprintf("%s.exe", exeName) 
        }
        exePath, err := exec.LookPath(exeName)

        if err != nil {
            log.Printf("WARNING: Executable '%s' is missing.\n", exeName)
        } else {
            log.Printf("Executable '%s' found in path '%s'\n", exeName, exePath)
            found += 1
        }
    }

    if found < len(requiredExe)  {
        log.Printf("==> WARNING: Not all executables were found in $PATH (%d/%d found)\n", found, len(requiredExe))
    } else {
        log.Printf("==> SUCCESS: All executables were found (%d/%d found).\n", found, len(requiredExe))
    }
}


