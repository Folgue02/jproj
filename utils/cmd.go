package utils

import (
	"log"
	"os"
	"os/exec"
)

func CMD(bin string, args ...string) error {
    command := exec.Command(bin, args...)
    command.Stdout = os.Stdout
    command.Stderr = os.Stderr
    command.Stdin = os.Stdin
    
    log.Printf("==> Executing command '%s' with args -< %v >-...\n", bin, args)
    return command.Run()
}
