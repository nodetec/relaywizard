package programs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nodetec/rwz/pkg/logs"
	"github.com/nodetec/rwz/pkg/utils/logging"
	"github.com/pterm/pterm"
)

// Function to determine an array of process IDs (pids) as strings for a given path to a program
func DeterminePidsOfProgram(currentUsername, programFilePath string) []string {
	var pidsOfProgram []string

	out, err := exec.Command("pidof", programFilePath).CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			errorCode := exitError.ExitCode()
			// pid for program not found or if pidof cannot access the process information
			if errorCode == 1 {
				return pidsOfProgram
			} else {
				logging.AppendRWZLogFile(currentUsername, logs.RWZLogFilePath, fmt.Sprintf("Failed to get the pid(s) for the program located at %s: %v", programFilePath, err))
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
