package logging

import (
	"log"
	"os"
	"os/exec"

	"github.com/nodetec/rwz/pkg/relays"
	"github.com/nodetec/rwz/pkg/utils/files"
	"github.com/pterm/pterm"
)

// Function to create rwz log file
func AppendRWZLogFile(currentUsername, rwzLogFilePath, logMessage string) {
	if currentUsername == relays.RootUser {
		rwzLogFile, err := os.OpenFile(rwzLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to open rwz log file: %v", err)
			os.Exit(1)
		}
		defer rwzLogFile.Close()

		log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
		log.SetOutput(rwzLogFile)
		log.Println(logMessage)

		files.SetPermissions(rwzLogFilePath, 0644)
	} else {
		err := exec.Command("sudo", "touch", rwzLogFilePath).Run()
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to create rwz log file: %v", err)
			os.Exit(1)
		}

		files.SetPermissionsUsingLinux(currentUsername, rwzLogFilePath, "0666")

		rwzLogFile, err := os.OpenFile(rwzLogFilePath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			pterm.Println()
			pterm.Error.Printfln("Failed to open rwz log file: %v", err)
			os.Exit(1)
		}
		defer rwzLogFile.Close()

		log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)
		log.SetOutput(rwzLogFile)
		log.Println(logMessage)

		files.SetPermissionsUsingLinux(currentUsername, rwzLogFilePath, "0644")
	}
}
