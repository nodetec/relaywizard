package systemd

import (
	"log"
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

func CreateEnvFile(envFilePath, envTemplate string, envFileParams *EnvFileParams) {
	envFile, err := os.Create(envFilePath)
	if err != nil {
		log.Fatalf("Error creating environment file: %v", err)
	}
	defer envFile.Close()

	envTmpl, err := template.New("env").Parse(envTemplate)
	if err != nil {
		log.Fatalf("Error parsing environment template: %v", err)
	}

	var WSProtocol string
	if envFileParams.HTTPSEnabled {
		WSProtocol = "wss"
	} else {
		WSProtocol = "ws"
	}

	err = envTmpl.Execute(envFile, struct{ Domain, WSProtocol, PrivKey, PubKey, RelayContact string }{Domain: envFileParams.Domain, WSProtocol: WSProtocol, PrivKey: envFileParams.PrivKey, PubKey: envFileParams.PubKey, RelayContact: envFileParams.RelayContact})
	if err != nil {
		log.Fatalf("Error executing environment template: %v", err)
	}
}

func CreateServiceFile(serviceFilePath, serviceTemplate string) {
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		log.Fatalf("Error creating service file: %v", err)
	}
	defer serviceFile.Close()

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		log.Fatalf("Error parsing service template: %v", err)
	}

	err = tmpl.Execute(serviceFile, struct{}{})
	if err != nil {
		log.Fatalf("Error executing service template: %v", err)
	}
}

func Reload() {
	err := exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		log.Fatalf("Error reloading systemd daemon: %v", err)
	}
}

func EnableService(name string) {
	err := exec.Command("systemctl", "enable", name).Run()
	if err != nil {
		log.Fatalf("Error enabling %s service: %v", name, err)
	}
}

func StartService(name string) {
	err := exec.Command("systemctl", "start", name).Run()
	if err != nil {
		log.Fatalf("Error starting %s service: %v", name, err)
	}
}

func ReloadService(name string) {
	err := exec.Command("systemctl", "reload", name).Run()
	if err != nil {
		log.Fatalf("Error reloading %s service: %v", name, err)
	}
}

func RestartService(name string) {
	err := exec.Command("systemctl", "restart", name).Run()
	if err != nil {
		log.Fatalf("Error restarting %s service: %v", name, err)
	}
}
