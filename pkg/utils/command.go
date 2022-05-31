package utils

import (
	"log"
	"os"
	"os/exec"
)

func RunCommand(command []string) error {
	log.Printf("Command: %v", command)
	cmd := exec.Command(command[0], command[1:]...) // TODO: possible security risk if run without sanitation
	log.Printf("Running command")
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start();
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}