package main

import (
	"log"
	"os"

	"github.com/folgue02/jproj/actions"
	"github.com/folgue02/jproj/actions/build"
	"github.com/folgue02/jproj/actions/clean"
	"github.com/folgue02/jproj/actions/createproject"
	"github.com/folgue02/jproj/actions/deps"
	jar "github.com/folgue02/jproj/actions/jar"
	"github.com/folgue02/jproj/actions/listproject"
	newElement "github.com/folgue02/jproj/actions/newElement"
	"github.com/folgue02/jproj/actions/run"
	"github.com/folgue02/jproj/actions/validate"
	"github.com/folgue02/jproj/utils"
	"github.com/folgue02/jproj/utils/actionmanager"
)

func main() {
    actions := actionmanager.ActionCollection {
        "createproject": { createproject.CreateProjectActionHandler, "Creates a new project." },
        "listproject": { listproject.ListProjectActionHandler, "Displays the information about the project." },
        "build": { build.BuildActionHandler, "Builds/Compiles the current project." },
        "clean": { clean.CleanActionHandler, "Cleans/Removes the compiled sources of the project." },
        "new": { newElement.NewElementActionHandler, "Adds/Creates a new element to the project." },
        "run": { run.RunProjectActionHandler, "Runs the current project." },
        "deps": { deps.DepMgrActionHandler, "Manage dependencies." },
        "jar": { jar.CreateJarActionHandler, "Creates a jar based on the project specified." },
        "version": { actions.VersionAction, "Displays jproj's version." },
        "validate": { validate.ValidateActionHandler, "Validates the environment." },
    }
	if len(os.Args) < 2 {
		actions.ExecuteAction("help", []string{})
		log.Println("Done.")
		os.Exit(0)
	}

    err, ok := actions.ExecuteAction(os.Args[1], os.Args[1:])
	if !ok {
		log.Printf("No action with name '%s' found.\n", os.Args[1])
		os.Exit(1)
	} else if err != nil {
        log.Printf("Error while executing action with name %s: \n", os.Args[1])
        utils.BeautifyError(err, func (section string) { log.Printf(section) })
        os.Exit(2)
    }
}
