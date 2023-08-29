package actions

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/akamensky/argparse"
	"github.com/folgue02/jproj/templates"
	"github.com/folgue02/jproj/utils"
)

const (
    ElementTypeClass = "class"
    ElementTypeInterface = "interface"
    ElementTypeEnum = "enum"
)

func newElement(args []string) error {
    parser := argparse.NewParser("new", "Add elements to the project.")    
    projectDirectory := parser.String(
        "d",
        "directory",
        &argparse.Options { Required: false, Default: ".", Help: "Specifies in what directory the project is." })
    elementType := parser.String(
        "t", 
        "type", 
        &argparse.Options { Required: false, Default: "class", Help: "Specifies what type of element." })
    elementName := parser.String(
        "n",
        "name",
        &argparse.Options { Required: true, Help: "Name of the new element" })
    err := parser.Parse(args)

    if err != nil {
        return fmt.Errorf("Error: Wrong arguments: %v", err)
    }
    
    jpp := utils.NewJavaPackagePath(*elementName)

    // Create package
    if *elementType == templates.ElementTypePackage {
        err := os.MkdirAll(path.Join(*projectDirectory, "src", jpp.ToJavaDirPath()), 0750)

        if err != nil {
            return fmt.Errorf("Error: Cannot create package due to the following error: %v", err)
        }
        return nil
    }

    // Create java source file (class, interface or enum)
    fileContent, ok := templates.GenerateTemplate(*elementType, jpp.Base(), jpp.DirPackagePath())

    if !ok {
        return fmt.Errorf("Error: '%s' is not identified as en element type.", *elementType)
    } 

    javaFilePath := path.Join(*projectDirectory, "src", jpp.ToJavaFilePath())
    javaFilePathDir := path.Dir(javaFilePath)

    err = os.MkdirAll(javaFilePathDir, 0750)

    if err != nil {
        return fmt.Errorf("Error: Cannot create the file's base path ('%s') due to the following error: %v", javaFilePathDir, err)
    }

    err = os.WriteFile(javaFilePath, []byte(*fileContent), 0750)
    
    if err != nil {
        return fmt.Errorf("Error: Cannot write to file '%s' due to the following error: %v", javaFilePath, err)
    } else {
        log.Println("Done.")
        return nil
    }
}

// Creates a Java file and all of its required directories.
func createJavaFile(packagePath, directory, content string) error {
    javaFilePath := utils.JavaPathToFilePath(packagePath)
    if err := os.MkdirAll(path.Base(javaFilePath), 0750); err != nil {
        return err
    }

    file, err := os.Create(javaFilePath)

    if err != nil {
        return err
    }

    _, err = file.Write([]byte(content))

    if err != nil {
        return err
    } else {
        return nil
    }
}
