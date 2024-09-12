package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Function to check if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && !info.IsDir()
}

// Function to extract the directory name from the domain
func GetDirectoryName(domainName string) string {
	domainParts := strings.Split(domainName, ".")
	if len(domainParts) > 2 {
		return domainParts[1]
	}
	return domainParts[0]
}

func UpdateDatabasePath(filePath string, newDBPath string) {
	// Define the file path and new DB path
	// filePath := "your_file.txt"         // Replace with your actual file path
	// newDBPath := "./new-db-location/"   // Replace with the new location you want

	// Read the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a slice to store the modified lines
	var lines []string

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Replace the line if it matches the pattern
		if strings.HasPrefix(line, "db = \"") {
			line = fmt.Sprintf("db = \"%s\"", newDBPath)
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Write the modified content back to the file
	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Database path updated successfully!")
}
