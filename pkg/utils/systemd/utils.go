package systemd

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"text/template"
)

type EnvFileParams struct {
	Domain       string
	HTTPSEnabled bool
	PrivKey      string
	PubKey       string
	RelayContact string
}

type ServiceFileParams struct {
	EnvFilePath    string
	BinaryFilePath string
	ConfigFilePath string
}

func CreateEnvFile(envFilePath, envTemplate string, envFileParams *EnvFileParams) {
	envFile, err := os.Create(envFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to create environment file: %v", err))
		os.Exit(1)
	}
	defer envFile.Close()

	envTmpl, err := template.New("env").Parse(envTemplate)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to parse environment template: %v", err))
		os.Exit(1)
	}

	var WSProtocol string
	if envFileParams.HTTPSEnabled {
		WSProtocol = "wss"
	} else {
		WSProtocol = "ws"
	}

	err = envTmpl.Execute(envFile, struct{ Domain, WSProtocol, PrivKey, PubKey, RelayContact string }{Domain: envFileParams.Domain, WSProtocol: WSProtocol, PrivKey: envFileParams.PrivKey, PubKey: envFileParams.PubKey, RelayContact: envFileParams.RelayContact})
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to execute environment template: %v", err))
		os.Exit(1)
	}
}

func CreateServiceFile(serviceFilePath, serviceTemplate string, serviceFileParams *ServiceFileParams) {
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to create service file: %v", err))
		os.Exit(1)
	}
	defer serviceFile.Close()

	serviceTmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to parse service template: %v", err))
		os.Exit(1)
	}

	err = serviceTmpl.Execute(serviceFile, struct{ EnvFilePath, BinaryFilePath, ConfigFilePath string }{EnvFilePath: serviceFileParams.EnvFilePath, BinaryFilePath: serviceFileParams.BinaryFilePath, ConfigFilePath: serviceFileParams.ConfigFilePath})
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to execute service template: %v", err))
		os.Exit(1)
	}
}

func Reload() {
	err := exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to reload systemd daemon: %v", err))
		os.Exit(1)
	}
}

func EnableService(name string) {
	err := exec.Command("systemctl", "enable", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to enable %s service: %v", name, err))
		os.Exit(1)
	}
}

func StartService(name string) {
	err := exec.Command("systemctl", "start", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to start %s service: %v", name, err))
		os.Exit(1)
	}
}

func ReloadService(name string) {
	err := exec.Command("systemctl", "reload", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to reload %s service: %v", name, err))
		os.Exit(1)
	}
}

func RestartService(name string) {
	err := exec.Command("systemctl", "restart", name).Run()
	if err != nil {
		pterm.Println()
		pterm.Error.Println(fmt.Sprintf("Failed to restart %s service: %v", name, err))
		os.Exit(1)
	}
}
