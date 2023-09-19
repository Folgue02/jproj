package createproject

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

type CreateProjectConfiguration struct {
    Name string
    Directory string
}

func NewCreateProjectConfiguration(args []string) (*CreateProjectConfiguration, error) {
    parser := argparse.NewParser("createproj", "Creates a new project")
    projectName := parser.String("n", "name", &argparse.Options { Required: true, Help: "Name of the new project"})
    projectBaseDirectory := parser.String("b", "basedirectory", &argparse.Options { Required: false, Default: ".", Help: "Directory in where the new project should be created." })

    if err := parser.Parse(args); err != nil {
        return nil, err
    }
    
    return &CreateProjectConfiguration {
        Name: *projectName,
        Directory: *projectBaseDirectory,
    }, nil
}

func CreateProjectActionHandler(args []string) error {
    createProjectConfiguration, err := NewCreateProjectConfiguration(args)

    if  err != nil {
        return fmt.Errorf("Error on argument parsing: %v", err)
    }

    return CreateProjectAction(*createProjectConfiguration)
}

func CreateProjectAction(createProjectConfiguration CreateProjectConfiguration) error {
    config := configuration.NewConfiguration(createProjectConfiguration.Name)

    if err := config.CreateProject(createProjectConfiguration.Directory); err != nil {
        return err
    }

    if err :=os.WriteFile(path.Join(createProjectConfiguration.Directory, config.ProjectName, "src", "App.java"), []byte("public class App {\n\tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello World!\");\n\t}\n}"), 0750); err != nil {
        return fmt.Errorf("Cannot create the start main class java file (%v)", err)
    }

    log.Printf("Project '%s' created.\n", createProjectConfiguration.Name)
    return nil
}
