package main

import (
	"log"
	"os"

	"github.com/folgue02/jproj/actions"
)

func main() {
	if len(os.Args) < 2 {
		actions.ExecuteAction("listactions", []string{})
		log.Println("Done.")
		os.Exit(0)
	}

	if !actions.ExecuteAction(os.Args[1], os.Args[1:]) {
		log.Printf("No action with name '%s' found.\n", os.Args[1])
		os.Exit(1)
	}
}
