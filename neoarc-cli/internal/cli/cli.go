package cli

import (
	"bufio"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	APIToken  string
	Version   string
	BuildTime string
	Commit    string
)

type Config struct {
	ServerURL   string `json:"server_url"`
	APIToken    string `json:"api_token,omitempty"`
	InsecureTLS bool   `json:"insecure_tls,omitempty"`
}

type AliasResponse struct {
	Success  bool   `json:"success"`
	Alias    string `json:"alias"`
	Command  string `json:"command"`
	ExecType string `json:"exec_type"`
	Message  string `json:"message"`
}

type CacheEntry struct {
	Command  string `json:"command"`
	ExecType string `json:"exec_type"`
	ETag     string `json:"etag"`
	Time     int64  `json:"time"`
}

type AliasCache struct {
	Entries map[string]CacheEntry `json:"entries"`
}

type TrustStore struct {
	Trusted map[string]int64 `json:"trusted"`
}

func ConfigDir() string {
	var configDir string
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(os.Getenv("APPDATA"), "neoarc")
	} else {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".neoarc")
	}
	os.MkdirAll(configDir, 0755)
	return configDir
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func CachePath() string {
	return filepath.Join(ConfigDir(), "alias_cache.json")
}

func TrustPath() string {
	return filepath.Join(ConfigDir(), "trusted.json")
}

func LoadConfig() Config {
	path := ConfigPath()
	file, err := os.ReadFile(path)
	if err != nil {
		defaultCfg := Config{ServerURL: "http://localhost:59248"}
		SaveConfig(defaultCfg)
		return defaultCfg
	}
	var cfg Config
	json.Unmarshal(file, &cfg)
	return cfg
}

func SaveConfig(cfg Config) {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(ConfigPath(), data, 0644)
}

func LoadCache() AliasCache {
	path := CachePath()
	file, err := os.ReadFile(path)
	if err != nil {
		return AliasCache{Entries: make(map[string]CacheEntry)}
	}
	var c AliasCache
	json.Unmarshal(file, &c)
	if c.Entries == nil {
		c.Entries = make(map[string]CacheEntry)
	}
	return c
}

func SaveCache(c AliasCache) {
	data, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(CachePath(), data, 0644)
}

func LoadTrusted() TrustStore {
	path := TrustPath()
	file, err := os.ReadFile(path)
	if err != nil {
		return TrustStore{Trusted: make(map[string]int64)}
	}
	var t TrustStore
	json.Unmarshal(file, &t)
	if t.Trusted == nil {
		t.Trusted = make(map[string]int64)
	}
	return t
}

func SaveTrusted(t TrustStore) {
	data, _ := json.MarshalIndent(t, "", "  ")
	os.WriteFile(TrustPath(), data, 0644)
}

func ResolveAPIToken(cfg Config) string {
	if cfg.APIToken != "" {
		return cfg.APIToken
	}
	return APIToken
}

func NewHTTPClient(cfg Config) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureTLS},
	}
	return &http.Client{Timeout: 15 * time.Second, Transport: tr}
}

func ConfirmExecution(aliasName string) bool {
	trusted := LoadTrusted()
	if _, ok := trusted.Trusted[aliasName]; ok {
		return true
	}

	fmt.Printf("Execute alias '%s'? This will run code from the remote server. [y/N]: ", aliasName)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))

	if line == "y" || line == "yes" {
		trusted.Trusted[aliasName] = time.Now().Unix()
		SaveTrusted(trusted)
		return true
	}
	return false
}

func FetchAlias(aliasName string) (*AliasResponse, error) {
	cfg := LoadConfig()
	cache := LoadCache()

	now := time.Now().Unix()
	entry, found := cache.Entries[aliasName]
	if found && (now-entry.Time) < 30 {
		return &AliasResponse{
			Success:  true,
			Alias:    aliasName,
			Command:  entry.Command,
			ExecType: entry.ExecType,
		}, nil
	}

	client := NewHTTPClient(cfg)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/alias/%s", cfg.ServerURL, aliasName), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token := ResolveAPIToken(cfg)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	if found && entry.ETag != "" {
		req.Header.Set("If-None-Match", entry.ETag)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w\nHint: Check server URL with: neoarc config <url> (default: http://localhost:59248)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 304 {
		return &AliasResponse{
			Success:  true,
			Alias:    aliasName,
			Command:  entry.Command,
			ExecType: entry.ExecType,
		}, nil
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response: %w", readErr)
	}

	var aliasResp AliasResponse
	if err := json.Unmarshal(body, &aliasResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		aliasResp.Success = false
		aliasResp.Message = "API authentication failed. Use 'neoarc config-token <token>' or rebuild the binary."
	}

	if aliasResp.Success {
		etag := resp.Header.Get("ETag")
		if etag == "" {
			h := sha256.Sum256(body)
			etag = fmt.Sprintf("%x", h[:16])
		}
		cache.Entries[aliasName] = CacheEntry{
			Command:  aliasResp.Command,
			ExecType: aliasResp.ExecType,
			ETag:     etag,
			Time:     now,
		}
		SaveCache(cache)
	}

	return &aliasResp, nil
}

func ExecuteCommand(cmdStr, execType string, aliasArgs []string) int {
	var cmd *exec.Cmd

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
		fmt.Fprintln(os.Stderr, "Error creating temp script:", err)
		return 1
	}

	if execType == "go" {
		if len(cmdStr) < 12 || cmdStr[:12] != "package main" {
			if _, err := tmpFile.WriteString("package main\n\n"); err != nil {
				fmt.Fprintln(os.Stderr, "Error writing temp script:", err)
				tmpFile.Close()
				os.Remove(tmpFile.Name())
				return 1
			}
		}
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(cmdStr); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing temp script:", err)
		tmpFile.Close()
		return 1
	}

	if err := tmpFile.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "Error closing temp script:", err)
		return 1
	}

	switch execType {
	case "bash":
		args := []string{tmpFile.Name()}
		args = append(args, aliasArgs...)
		cmd = exec.Command("bash", args...)
	case "powershell":
		args := []string{"-ExecutionPolicy", "Bypass", "-File", tmpFile.Name()}
		args = append(args, aliasArgs...)
		cmd = exec.Command("powershell", args...)
	case "python":
		args := []string{tmpFile.Name()}
		args = append(args, aliasArgs...)
		cmd = exec.Command("python", args...)
	case "go":
		args := []string{"run", tmpFile.Name()}
		args = append(args, aliasArgs...)
		cmd = exec.Command("go", args...)
	default:
		if runtime.GOOS == "windows" {
			args := []string{"/C", tmpFile.Name()}
			args = append(args, aliasArgs...)
			cmd = exec.Command("cmd", args...)
		} else {
			args := []string{tmpFile.Name()}
			args = append(args, aliasArgs...)
			cmd = exec.Command("sh", args...)
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintln(os.Stderr, "Execution error:", err)
		return 1
	}
	return 0
}

type AliasesListResponse struct {
	Success bool     `json:"success"`
	Aliases []string `json:"aliases"`
	Message string   `json:"message"`
}

func FetchAliases() ([]string, error) {
	cfg := LoadConfig()
	client := NewHTTPClient(cfg)

	req, err := http.NewRequest("GET", cfg.ServerURL+"/api/aliases", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token := ResolveAPIToken(cfg)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response: %w", readErr)
	}

	var listResp AliasesListResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return []string{}, nil
	}

	if !listResp.Success {
		return nil, fmt.Errorf("server error: %s", listResp.Message)
	}

	return listResp.Aliases, nil
}

func installDir() string {
	dir := filepath.Join(homeDir(), ".config", "neostore", "neoarc", "bin")
	os.MkdirAll(dir, 0755)
	return dir
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	h, _ := os.UserHomeDir()
	return h
}

func selfInstall() int {
	targetDir := installDir()
	binName := "neoarc"
	if runtime.GOOS == "windows" {
		binName = "neoarc.exe"
	}
	targetPath := filepath.Join(targetDir, binName)

	fmt.Println(">>> Installing NeoArc...")

	// Determine download URL for latest version
	repo := "rkriad585/NeoArc"
	version := "v3.0.1"
	var downloadName string

	switch runtime.GOOS {
	case "windows":
		downloadName = "neoarc-windows-amd64.exe"
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			downloadName = "neoarc-darwin-arm64"
		default:
			downloadName = "neoarc-darwin-amd64"
		}
	case "linux":
		switch runtime.GOARCH {
		case "arm64":
			downloadName = "neoarc-linux-arm64"
		default:
			downloadName = "neoarc-linux-amd64"
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return 1
	}

	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, version, downloadName)

	fmt.Printf(">>> Downloading %s\n", url)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		fmt.Println("Falling back to copying the current binary...")
		return copySelf(targetPath)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Download failed: HTTP %d\n", resp.StatusCode)
		fmt.Println("Falling back to copying the current binary...")
		return copySelf(targetPath)
	}

	out, err := os.Create(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		return 1
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		out.Close()
		os.Remove(targetPath)
		return 1
	}
	out.Close()

	if written == 0 {
		fmt.Fprintln(os.Stderr, "Downloaded empty file")
		os.Remove(targetPath)
		fmt.Println("Falling back to copying the current binary...")
		return copySelf(targetPath)
	}

	if runtime.GOOS != "windows" {
		os.Chmod(targetPath, 0755)
	}

	fmt.Printf("OK   Installed to %s (%d bytes)\n", targetPath, written)

	return addToPath(targetDir)
}

func copySelf(targetPath string) int {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving binary path: %v\n", err)
		return 1
	}

	src, err := os.Open(exe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening current binary: %v\n", err)
		return 1
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		return 1
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		fmt.Fprintf(os.Stderr, "Error copying binary: %v\n", err)
		dst.Close()
		os.Remove(targetPath)
		return 1
	}
	dst.Close()

	if runtime.GOOS != "windows" {
		os.Chmod(targetPath, 0755)
	}

	fmt.Printf("OK   Installed to %s\n", targetPath)
	return addToPath(filepath.Dir(targetPath))
}

func addToPath(targetDir string) int {
	if runtime.GOOS == "windows" {
		currentPath := os.Getenv("Path")
		if !strings.Contains(currentPath, targetDir) {
			newPath := currentPath + ";" + targetDir
			os.Setenv("Path", newPath)
			fmt.Println("OK   Added to PATH for this session.")
			fmt.Println(">>> To make it permanent, run the following in an Administrator PowerShell:")
			fmt.Printf("    [Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path','User')+'%s;%s','User')\n", ";", targetDir)
			fmt.Println("    Or run: installer.ps1")
		} else {
			fmt.Println("OK   Already in PATH.")
		}
	} else {
		rcFile := ""
		switch {
		case os.Getenv("SHELL") != "" && strings.Contains(os.Getenv("SHELL"), "zsh"):
			rcFile = filepath.Join(homeDir(), ".zshrc")
		case os.Getenv("SHELL") != "" && strings.Contains(os.Getenv("SHELL"), "bash"):
			if runtime.GOOS == "darwin" {
				rcFile = filepath.Join(homeDir(), ".bash_profile")
			} else {
				rcFile = filepath.Join(homeDir(), ".bashrc")
			}
		default:
			rcFile = filepath.Join(homeDir(), ".profile")
		}

		line := fmt.Sprintf("export PATH=\"$PATH:%s\"", targetDir)
		data, _ := os.ReadFile(rcFile)
		if !strings.Contains(string(data), targetDir) {
			f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not write to %s: %v\n", rcFile, err)
				fmt.Printf("Add this line manually:\n  %s\n", line)
			} else {
				defer f.Close()
				fmt.Fprintln(f)
				fmt.Fprintln(f, "# Added by NeoArc installer")
				fmt.Fprintln(f, line)
				fmt.Printf("OK   Added to PATH in %s\n", rcFile)
				fmt.Println(">>> Run 'source", rcFile, "' or restart your terminal.")
			}
		} else {
			fmt.Println("OK   Already in PATH.")
		}
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  NeoArc installed successfully!")
	fmt.Printf("  Binary : %s\n", targetDir)
	fmt.Println("  Usage  : neoarc help")
	fmt.Println("========================================")
	return 0
}

func selfUninstall() int {
	configDir := ConfigDir()

	fmt.Println(">>> Uninstalling NeoArc...")

	if _, err := os.Stat(configDir); err == nil {
		if err := os.RemoveAll(configDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing config directory: %v\n", err)
			return 1
		}
		fmt.Println("OK   Removed config directory:", configDir)
	} else {
		fmt.Println("OK   No config directory found.")
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving binary path: %v\n", err)
		return 1
	}

	realPath, err := filepath.EvalSymlinks(exePath)
	if err == nil {
		exePath = realPath
	}

	if runtime.GOOS == "windows" {
		batContent := fmt.Sprintf("@echo off\r\ntimeout /t 1 /nobreak >nul\r\ndel /f /q \"%s\"\r\nif exist \"%s\" (echo ERR Failed to delete binary) else (echo OK   Deleted binary: %s)\r\ndel /f /q \"%%~f0\"\r\n", exePath, exePath, exePath)
		batPath := filepath.Join(os.TempDir(), "neoarc-uninstall.bat")
		if err := os.WriteFile(batPath, []byte(batContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating uninstall script: %v\n", err)
			fmt.Println("Please manually delete:", exePath)
			return 1
		}
		fmt.Println("OK   Uninstall script created. Binary will be deleted shortly.")
		cmd := exec.Command("cmd", "/C", "start", "/B", batPath)
		cmd.Stderr = os.Stderr
		cmd.Start()
	} else {
		if err := os.Remove(exePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting binary: %v\n", err)
			fmt.Println("Please manually delete:", exePath)
			return 1
		}
		fmt.Println("OK   Deleted binary:", exePath)
	}

	fmt.Println()
	fmt.Println("To remove NeoArc from your PATH, edit your shell rc file")
	fmt.Println("and delete the line containing 'neostore/neoarc/bin'.")

	if runtime.GOOS == "windows" {
		fmt.Println("Or run: installer.ps1 --selfuninstall")
	} else {
		fmt.Println("Or run: installer.sh --selfuninstall")
	}

	fmt.Println("Restart your terminal for PATH changes to take effect.")
	return 0
}

func ShowHelp() {
	fmt.Println(`NeoArc - Cross-Platform Alias Executor
Usage:
  neoarc get <alias> [args...]       : Print the code for the alias
  neoarc run <alias> [args...]       : Run the alias command with args
  neoarc <alias> [args...]           : Run the alias command (shorthand)
  neoarc config <server-url>         : Set the NeoArc Web Server URL
  neoarc config-token <tok>          : Set the API authentication token
  neoarc config insecure             : Enable insecure TLS (skip certificate verify)
  neoarc config secure               : Disable insecure TLS (default, verify certs)
  neoarc help                        : Show help menu
  neoarc completion <shell>          : Generate shell completion script (bash|zsh|powershell)

Standalone flags:
  --install                          : Download and install NeoArc to ~/.config/neostore/neoarc/bin/
  --selfuninstall                    : Remove NeoArc config, cache, and binary from the system

Options (place before <alias>):
  --dry-run                          : Print the alias code without executing
  --yes                              : Skip execution confirmation prompt

Arguments after <alias> are passed through to the executed command.
  bash/sh    : accessible via $1, $2, $@
  powershell : accessible via $args[0], $args[1]
  python     : accessible via sys.argv[1], sys.argv[2]
  cmd        : accessible via %1, %2`)
}

func Run(args []string) int {
	if len(args) < 2 {
		ShowHelp()
		return 0
	}

	if args[1] == "help" {
		ShowHelp()
		return 0
	}

	if args[1] == "config" && len(args) >= 3 {
		if args[2] == "insecure" {
			cfg := LoadConfig()
			cfg.InsecureTLS = true
			SaveConfig(cfg)
			fmt.Println("Insecure TLS enabled.")
			return 0
		}
		if args[2] == "secure" {
			cfg := LoadConfig()
			cfg.InsecureTLS = false
			SaveConfig(cfg)
			fmt.Println("Insecure TLS disabled.")
			return 0
		}
		if len(args) == 3 {
			cfg := LoadConfig()
			cfg.ServerURL = args[2]
			SaveConfig(cfg)
			fmt.Println("Config updated! Server URL:", cfg.ServerURL)
			return 0
		}
	}

	if args[1] == "config-token" && len(args) == 3 {
		cfg := LoadConfig()
		cfg.APIToken = args[2]
		SaveConfig(cfg)
		fmt.Println("API token updated.")
		return 0
	}

	if args[1] == "completion" {
		if len(args) < 3 {
			fmt.Println("Usage: neoarc completion <shell>\nSupported shells: bash, zsh, powershell")
			return 0
		}
		aliases, _ := FetchAliases()
		if aliases == nil {
			aliases = []string{}
		}
		return GenerateCompletion(args[2], aliases)
	}

	if args[1] == "--install" {
		return selfInstall()
	}

	if args[1] == "--selfuninstall" {
		return selfUninstall()
	}

	if args[1] == "_list_aliases" {
		aliases, _ := FetchAliases()
		for _, a := range aliases {
			fmt.Println(a)
		}
		return 0
	}

	dryRun := false
	yesMode := false
	argIdx := 1

	for argIdx < len(args) && args[argIdx][0] == '-' {
		switch args[argIdx] {
		case "--dry-run":
			dryRun = true
		case "--yes":
			yesMode = true
		default:
			fmt.Fprintln(os.Stderr, "Unknown flag:", args[argIdx])
			return 1
		}
		argIdx++
	}

	if argIdx >= len(args) {
		ShowHelp()
		return 0
	}

	cmd := args[argIdx]
	argIdx++

	var aliasName string
	var aliasArgs []string
	isGet := false

	if cmd == "get" {
		if argIdx >= len(args) {
			fmt.Println("Usage: neoarc get <alias-name>")
			return 1
		}
		aliasName = args[argIdx]
		aliasArgs = args[argIdx+1:]
		isGet = true
	} else if cmd == "run" {
		if argIdx >= len(args) {
			fmt.Println("Usage: neoarc run <alias-name>")
			return 1
		}
		aliasName = args[argIdx]
		aliasArgs = args[argIdx+1:]
	} else {
		aliasName = cmd
		aliasArgs = args[argIdx:]
	}

	alias, err := FetchAlias(aliasName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return 1
	}

	if !alias.Success {
		fmt.Println("Error:", alias.Message)
		return 1
	}

	if isGet || dryRun {
		if dryRun {
			fmt.Println("--- Dry Run ---")
		}
		fmt.Printf("--- NeoArc Alias: %s (%s) ---\n%s\n------------------------\n", alias.Alias, alias.ExecType, alias.Command)
		return 0
	}

	if !yesMode && !ConfirmExecution(aliasName) {
		fmt.Println("Execution cancelled.")
		return 1
	}

	return ExecuteCommand(alias.Command, alias.ExecType, aliasArgs)
}
