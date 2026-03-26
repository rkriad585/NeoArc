// main.go

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Config struct {
	ServerURL string `json:"server_url"`
}

type AliasResponse struct {
	Success  bool   `json:"success"`
	Alias    string `json:"alias"`
	Command  string `json:"command"`
	ExecType string `json:"exec_type"`
	Message  string `json:"message"`
}

func getConfigPath() string {
	var configDir string
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(os.Getenv("APPDATA"), "neoarc")
	} else {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".neoarc")
	}
	os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "config.json")
}

func loadConfig() Config {
	path := getConfigPath()
	file, err := os.ReadFile(path)
	if err != nil {
		defaultCfg := Config{ServerURL: "http://localhost:5000"} // Change to your public IP/Domain
		saveConfig(defaultCfg)
		return defaultCfg
	}
	var cfg Config
	json.Unmarshal(file, &cfg)
	return cfg
}

func saveConfig(cfg Config) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(getConfigPath(), data, 0644)
}

func fetchAlias(aliasName string) (*AliasResponse, error) {
	cfg := loadConfig()
	resp, err := http.Get(fmt.Sprintf("%s/api/alias/%s", cfg.ServerURL, aliasName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var aliasResp AliasResponse
	json.Unmarshal(body, &aliasResp)
	return &aliasResp, nil
}

func executeCommand(cmdStr, execType string) {
	var cmd *exec.Cmd

	// Determine proper file extension
	ext := ".sh"
	if execType == "powershell" {
		ext = ".ps1"
	} else if execType == "python" {
		ext = ".py"
	} else if execType == "go" {
		ext = ".go"
	} else if runtime.GOOS == "windows" && (execType == "cmd" || execType == "") {
		ext = ".bat"
	}

	tmpFile, err := os.CreateTemp("", "neoarc-*"+ext)
	if err != nil {
		fmt.Println("Error creating temp script:", err)
		return
	}

	// Add Go module structure automatically if the user just provides plain Go code without a package declaration
	if execType == "go" {
		// If the user forgot 'package main', prepend it so 'go run' works natively
		if len(cmdStr) < 12 || cmdStr[0:12] != "package main" {
			tmpFile.WriteString("package main\n\n")
		}
	}

	defer os.Remove(tmpFile.Name()) // Clean up after execution

	tmpFile.WriteString(cmdStr)
	tmpFile.Close()

	switch execType {
	case "bash":
		cmd = exec.Command("bash", tmpFile.Name())
	case "powershell":
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tmpFile.Name())
	case "python":
		cmd = exec.Command("python", tmpFile.Name())
	case "go":
		cmd = exec.Command("go", "run", tmpFile.Name()) // Execute Go code
	default:
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", tmpFile.Name())
		} else {
			cmd = exec.Command("sh", tmpFile.Name())
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func showHelp() {
	fmt.Println(`NeoArc - Cross-Platform Alias Executor
Usage:
  neoarc get <alias>         : Print the code for the alias
  neoarc <alias>             : Run the alias command
  neoarc run <alias>         : Run the alias command
  neoarc <alias>             : Run the alias command (shorthand)
  neoarc config <server-url> : Set the NeoArc Web Server URL
  neoarc help                : Show help menu`)
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	cmd := os.Args[1]

	if cmd == "help" {
		showHelp()
		return
	}

	if cmd == "config" && len(os.Args) == 3 {
		cfg := loadConfig()
		cfg.ServerURL = os.Args[2]
		saveConfig(cfg)
		fmt.Println("Config updated! Server URL:", cfg.ServerURL)
		return
	}

	var aliasName string
	isGet := false

	if cmd == "get" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: neoarc get <alias-name>")
			return
		}
		aliasName = os.Args[2]
		isGet = true
	} else if cmd == "run" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: neoarc run <alias-name>")
			return
		}
		aliasName = os.Args[2]
	} else {
		aliasName = cmd
	}

	alias, err := fetchAlias(aliasName)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	if !alias.Success {
		fmt.Println("Error:", alias.Message)
		return
	}

	if isGet {
		fmt.Printf("--- NeoArc Alias: %s (%s) ---\n%s\n------------------------\n", alias.Alias, alias.ExecType, alias.Command)
	} else {
		executeCommand(alias.Command, alias.ExecType)
	}
}
