package utils

import (
	"log"
	"os"
	"os/exec"
)

// Equivalent to 'exec.Command(...)', but also also inherits the proccess'
// io. It also displays a message showing the whole command.
func CMD(bin string, args ...string) error {
    command := exec.Command(bin, args...)
    command.Stdout = os.Stdout
    command.Stderr = os.Stderr
    command.Stdin = os.Stdin
    
    log.Printf("==> Executing command '%s' with args -< %v >-...\n", bin, args)
    return command.Run()
}
