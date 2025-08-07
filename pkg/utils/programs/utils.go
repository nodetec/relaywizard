package programs

import (
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"strings"
)

// Function to determine an array of process IDs (pids) as strings for a given path to a program
func DeterminePidsOfProgram(programFilePath string) []string {
	var pidsOfProgram []string

	out, err := exec.Command("pidof", programFilePath).CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			errorCode := exitError.ExitCode()
			// pid for program not found
			if errorCode == 1 {
				return pidsOfProgram
			} else {
				pterm.Println()
				pterm.Error.Printfln("Failed to get the pid(s) for the program located at %s: %v", programFilePath, err)
				os.Exit(1)
			}
		}
	}

	pidofOutput := string(out)

	pidsOfProgram = strings.Fields(pidofOutput)

	return pidsOfProgram
}
