package actions

import (
	"io/ioutil"
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

func newElement(args []string) {
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
        log.Printf("Error: Wrong arguments: %v\n", err)
        return
    }
    
    jpp := utils.NewJavaPackagePath(*elementName)

    // Create package
    if *elementType == templates.ElementTypePackage {
        err := os.MkdirAll(path.Join(*projectDirectory, "src", jpp.ToJavaDirPath()), 0750)

        if err != nil {
            log.Printf("Error: Cannot create package due to the following error: %v\n", err)
            return
        }
        return
    }

    // Create java source file (class, interface or enum)
    fileContent, ok := templates.GenerateTemplate(*elementType, jpp.Base(), jpp.DirPackagePath())

    if !ok {
        log.Printf("Error: '%s' is not identified as en element type.\n", *elementType)
        return
    } 

    javaFilePath := path.Join(*projectDirectory, "src", jpp.ToJavaFilePath())
    javaFilePathDir := path.Dir(javaFilePath)

    err = os.MkdirAll(javaFilePathDir, 0750)

    if err != nil {
        log.Printf("Error: Cannot create the file's base path ('%s') due to the following error: %v\n", javaFilePathDir, err)
        return
    }

    err = ioutil.WriteFile(javaFilePath, []byte(*fileContent), 0750)
    
    if err != nil {
        log.Printf("Error: Cannot write to file '%s' due to the following error: %v\n", javaFilePath, err)
    } else {
        log.Println("Done.")
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
