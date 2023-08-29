package actions

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/configuration"
)

func createproject(args []string) error {
    parser := argparse.NewParser("createproj", "Creates a new project")
    projectName := parser.String("n", "name", &argparse.Options { Required: true, Help: "Name of the new project"})
    projectBaseDirectory := parser.String("b", "basedirectory", &argparse.Options { Required: false, Default: ".", Help: "Directory in where the new project should be created." })
    if err := parser.Parse(args); err != nil {
        return fmt.Errorf("Error on argument parsing: %v", err)

    }
    config := configuration.NewConfiguration(*projectName)

    if err := config.CreateProject(*projectBaseDirectory); err != nil {
        return err
    }

    if err :=os.WriteFile(path.Join(*projectBaseDirectory, config.ProjectName, "src", "App.java"), []byte("public class App {\n\tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello World!\");\n\t}\n}"), 0750); err != nil {
        return fmt.Errorf("Cannot create the start main class java file (%v)", err)
    }

    log.Printf("Project '%s' created.\n", *projectName)
    return nil
}