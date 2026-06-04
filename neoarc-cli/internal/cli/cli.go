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
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"neoarc/internal/config"
	"neoarc/internal/version"
)

const (
	RepoPath  = "rkriad585/NeoArc"
	RepoOwner = "rkriad585"
)

var (
	APIToken   string
	Version    string
	BuildTime  string
	Commit     string
	configPath string
)

func init() {
	if version.Version == "" {
		version.Version = Version
	}
	if version.Commit == "" {
		version.Commit = Commit
	}
	if version.BuildTime == "" {
		version.BuildTime = BuildTime
	}
}

type Config struct {
	ServerURL   string `toml:"server_url" json:"server_url"`
	APIToken    string `toml:"api_token,omitempty" json:"api_token,omitempty"`
	InsecureTLS bool   `toml:"insecure_tls,omitempty" json:"insecure_tls,omitempty"`
	Theme       string `toml:"theme,omitempty" json:"theme,omitempty"`
}

func resolveTheme(cfg Config) config.RoleColors {
	return config.ResolveTheme(cfg.Theme)
}

func sprintTheme(cfg Config, role config.Role, text string) string {
	t, ok := config.FindTheme(cfg.Theme)
	if !ok {
		t = config.DefaultTheme()
	}
	return config.Colorize(text, t.Hex(role))
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
	config.EnsureConfigDir()
	return config.ConfigDir()
}

func migrateOldConfig() {
	newDir := config.ConfigDir()

	oldDirs := []string{
		filepath.Join(os.Getenv("APPDATA"), "neoarc"),
	}
	if h, err := os.UserHomeDir(); err == nil {
		oldDirs = append(oldDirs, filepath.Join(h, ".neoarc"))
	}

	for _, old := range oldDirs {
		if old == "" || old == newDir {
			continue
		}
		entries, err := os.ReadDir(old)
		if err != nil {
			continue
		}
		if len(entries) == 0 {
			os.Remove(old)
			continue
		}
		os.MkdirAll(newDir, 0755)
		for _, e := range entries {
			oldPath := filepath.Join(old, e.Name())
			newPath := filepath.Join(newDir, e.Name())
			if _, err := os.Stat(newPath); err == nil {
				continue
			}
			os.Rename(oldPath, newPath)
		}
		remaining, _ := os.ReadDir(old)
		if len(remaining) == 0 {
			os.Remove(old)
		}
	}
}

func ConfigPath() string {
	if configPath != "" {
		return configPath
	}
	return config.ConfigFile("config.toml")
}

func CachePath() string {
	return config.ConfigFile("alias_cache.json")
}

func TrustPath() string {
	return config.ConfigFile("trusted.json")
}

func migrateConfigFormat() {
	oldPath := config.ConfigFile("config.json")
	newPath := ConfigPath()
	if _, err := os.Stat(newPath); err == nil {
		return
	}
	data, err := os.ReadFile(oldPath)
	if err != nil {
		return
	}
	var old Config
	if err := json.Unmarshal(data, &old); err != nil {
		return
	}
	SaveConfig(old)
	os.Remove(oldPath)
}

func LoadConfig() Config {
	migrateOldConfig()
	migrateConfigFormat()
	config.EnsureConfigDir()
	path := ConfigPath()
	file, err := os.ReadFile(path)
	if err != nil {
		defaultCfg := Config{ServerURL: "http://localhost:59248"}
		SaveConfig(defaultCfg)
		return defaultCfg
	}
	var cfg Config
	if err := toml.Unmarshal(file, &cfg); err != nil {
		defaultCfg := Config{ServerURL: "http://localhost:59248"}
		SaveConfig(defaultCfg)
		return defaultCfg
	}
	return cfg
}

func SaveConfig(cfg Config) {
	config.EnsureConfigDir()
	var buf strings.Builder
	if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
		return
	}
	os.WriteFile(ConfigPath(), []byte(buf.String()), 0644)
}

func LoadCache() AliasCache {
	migrateOldConfig()
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
	config.EnsureConfigDir()
	data, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(CachePath(), data, 0644)
}

func LoadTrusted() TrustStore {
	migrateOldConfig()
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
	config.EnsureConfigDir()
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

	cfg := LoadConfig()
	p := resolveTheme(cfg)
	fmt.Printf("Execute alias '%s'? %s [y/N]: ", p.Accent.Sprint(aliasName), p.Warning.Sprint("This will run code from the remote server."))
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
		entry.Time = now
		cache.Entries[aliasName] = entry
		SaveCache(cache)
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
	dir := filepath.Join(config.ConfigDir(), "bin")
	os.MkdirAll(dir, 0755)
	return dir
}

func resolveVersion(repo string) string {
	if Version != "" {
		return Version
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/main/.version", repo))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body))
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
	version := resolveVersion(repo)
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

	var url string
	if version != "" {
		url = fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, version, downloadName)
	} else {
		url = fmt.Sprintf("https://github.com/%s/releases/latest/download/%s", repo, downloadName)
	}

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
			fmt.Printf("    [Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path','User')+';%s','User')\n", targetDir)
			fmt.Println("    Or run: installer.ps1")
		} else {
			fmt.Println("OK   Already in PATH.")
		}
	} else {
		rcFile := ""
		switch {
		case os.Getenv("SHELL") != "" && strings.Contains(os.Getenv("SHELL"), "zsh"):
			rcFile = filepath.Join(config.HomeDir(), ".zshrc")
		case os.Getenv("SHELL") != "" && strings.Contains(os.Getenv("SHELL"), "bash"):
			if runtime.GOOS == "darwin" {
				rcFile = filepath.Join(config.HomeDir(), ".bash_profile")
			} else {
				rcFile = filepath.Join(config.HomeDir(), ".bashrc")
			}
		default:
			rcFile = filepath.Join(config.HomeDir(), ".profile")
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

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving binary path: %v\n", err)
		return 1
	}
	realPath, err := filepath.EvalSymlinks(exePath)
	if err == nil {
		exePath = realPath
	}

	// Check whether the running binary lives inside the config directory.
	binaryInside := strings.HasPrefix(exePath, configDir+string(os.PathSeparator))

	if runtime.GOOS == "windows" && binaryInside {
		// Windows cannot delete a running executable.  Remove everything
		// *except* the running binary, then launch a deferred batch script
		// that waits 1s, removes the binary and the now-empty directories.
		filepath.WalkDir(configDir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if !strings.EqualFold(path, exePath) {
				os.Remove(path)
			}
			return nil
		})
		// Remove the empty bin/ sub-directory (the exe is still there so
		// the rmdir below will handle it).
		binDir := filepath.Dir(exePath)
		if entries, _ := os.ReadDir(binDir); len(entries) == 1 {
			os.Remove(binDir)
		}
		fmt.Println("OK   Removed config files.")

		batContent := fmt.Sprintf("@echo off\r\ntimeout /t 1 /nobreak >nul\r\nrmdir /s /q \"%s\" 2>nul\r\necho OK   NeoArc has been uninstalled.\r\ndel /f /q \"%%~f0\"\r\n", configDir)
		batPath := filepath.Join(os.TempDir(), "neoarc-uninstall.bat")
		if err := os.WriteFile(batPath, []byte(batContent), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating uninstall script: %v\n", err)
			fmt.Println("Please manually delete:", exePath)
			return 1
		}
		fmt.Println("OK   Uninstall script created. NeoArc will be fully removed shortly.")
		exec.Command("cmd", "/C", "start", "/B", batPath).Start()
	} else {
		// Unix: RemoveAll works even with a running binary (inode lives).
		// Windows (binary outside config): RemoveAll works fine.
		if _, err := os.Stat(configDir); err == nil {
			if err := os.RemoveAll(configDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing config directory: %v\n", err)
				return 1
			}
			fmt.Println("OK   Removed config directory:", configDir)
		} else {
			fmt.Println("OK   No config directory found.")
		}

		if runtime.GOOS == "windows" {
			batContent := fmt.Sprintf("@echo off\r\ntimeout /t 1 /nobreak >nul\r\ndel /f /q \"%s\" 2>nul\r\necho OK   NeoArc has been uninstalled.\r\ndel /f /q \"%%~f0\"\r\n", exePath)
			batPath := filepath.Join(os.TempDir(), "neoarc-uninstall.bat")
			if err := os.WriteFile(batPath, []byte(batContent), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating uninstall script: %v\n", err)
				fmt.Println("Please manually delete:", exePath)
				return 1
			}
			fmt.Println("OK   Uninstall script created. Binary will be deleted shortly.")
			exec.Command("cmd", "/C", "start", "/B", batPath).Start()
		} else {
			if err := os.Remove(exePath); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting binary: %v\n", err)
				fmt.Println("Please manually delete:", exePath)
				return 1
			}
			fmt.Println("OK   Deleted binary:", exePath)
		}
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

func fetchLatestVersion() string {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://raw.githubusercontent.com/%s/main/.version", RepoPath))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	v := strings.TrimSpace(string(body))
	if v == "" {
		return ""
	}
	return v
}

func compareVersions(a, b string) int {
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")
	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}
	for i := 0; i < maxLen; i++ {
		var na, nb int
		if i < len(partsA) {
			na, _ = strconv.Atoi(partsA[i])
		}
		if i < len(partsB) {
			nb, _ = strconv.Atoi(partsB[i])
		}
		if na < nb {
			return -1
		}
		if na > nb {
			return 1
		}
	}
	return 0
}

func resolveDownloadName() string {
	switch runtime.GOOS {
	case "windows":
		return config.ProjectName + "-windows-amd64.exe"
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return config.ProjectName + "-darwin-arm64"
		}
		return config.ProjectName + "-darwin-amd64"
	case "linux":
		if runtime.GOARCH == "arm64" {
			return config.ProjectName + "-linux-arm64"
		}
		return config.ProjectName + "-linux-amd64"
	}
	return ""
}

type DownloadProgress struct {
	Total uint64
}

func (dp *DownloadProgress) Write(p []byte) (int, error) {
	n := len(p)
	dp.Total += uint64(n)
	fmt.Printf("\r  Downloaded: %.2f MB", float64(dp.Total)/1024/1024)
	return n, nil
}

func selfUpdate(proxyURL string) int {
	cfg := LoadConfig()
	p := resolveTheme(cfg)

	if proxyURL != "" {
		fmt.Printf("  Using proxy: %s\n", proxyURL)
		os.Setenv("HTTP_PROXY", proxyURL)
		os.Setenv("HTTPS_PROXY", proxyURL)
	}

	fmt.Println(p.Primary.Sprint(">>> Checking for updates..."))

	current := resolveVersion(RepoPath)
	if current == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine current version.")
		return 1
	}

	latest := fetchLatestVersion()
	if latest == "" {
		fmt.Fprintln(os.Stderr, "Error: could not fetch latest version from GitHub.")
		return 1
	}

	fmt.Printf("    Current: %s\n", p.Success.Sprint(current))
	fmt.Printf("    Latest : %s\n", p.Accent.Sprint(latest))

	cmp := compareVersions(current, latest)
	if cmp >= 0 {
		fmt.Println(p.Success.Sprint("OK   Already up to date."))
		return 0
	}

	fmt.Printf(">>> New version %s available. Updating...\n", p.Accent.Sprint(latest))

	downloadName := resolveDownloadName()
	if downloadName == "" {
		fmt.Fprintln(os.Stderr, "Error: unsupported platform:", runtime.GOOS+"/"+runtime.GOARCH)
		return 1
	}

	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", RepoPath, latest, downloadName)

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: download failed: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error: download failed (HTTP %d): %s\n", resp.StatusCode, strings.TrimSpace(string(body)))
		return 1
	}

	tmpFile, err := os.CreateTemp("", config.ProjectName+"-update-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create temp file: %v\n", err)
		return 1
	}
	tmpPath := tmpFile.Name()

	progress := &DownloadProgress{}
	teeReader := io.TeeReader(resp.Body, progress)
	written, err := io.Copy(tmpFile, teeReader)
	tmpFile.Close()
	fmt.Println()
	if err != nil {
		os.Remove(tmpPath)
		fmt.Fprintf(os.Stderr, "Error: could not write update: %v\n", err)
		return 1
	}
	if written == 0 {
		os.Remove(tmpPath)
		fmt.Fprintln(os.Stderr, "Error: downloaded empty file.")
		return 1
	}

	chmod := runtime.GOOS != "windows"
	if chmod {
		os.Chmod(tmpPath, 0755)
	}

	exePath, err := os.Executable()
	if err != nil {
		os.Remove(tmpPath)
		fmt.Fprintf(os.Stderr, "Error: could not resolve binary path: %v\n", err)
		return 1
	}

	if runtime.GOOS == "windows" {
		oldPath := exePath + ".old"
		os.Remove(oldPath)
		if err := os.Rename(exePath, oldPath); err != nil {
			os.Remove(tmpPath)
			fmt.Fprintf(os.Stderr, "Error: could not backup current binary: %v\n", err)
			return 1
		}
		if err := os.Rename(tmpPath, exePath); err != nil {
			os.Rename(oldPath, exePath)
			os.Remove(tmpPath)
			fmt.Fprintf(os.Stderr, "Error: could not replace binary, restored original: %v\n", err)
			return 1
		}
		fmt.Printf("OK   Updated to %s (%s)\n", p.Success.Sprint(latest), exePath)
		fmt.Println("Note: you can safely delete " + oldPath + " after this session.")
	} else {
		if err := os.Rename(tmpPath, exePath); err != nil {
			os.Remove(tmpPath)
			fmt.Fprintf(os.Stderr, "Error: could not replace binary: %v\n", err)
			return 1
		}
		os.Chmod(exePath, 0755)
		fmt.Printf("OK   Updated to %s (%s)\n", p.Success.Sprint(latest), exePath)
	}

	return 0
}

func ShowHelp() {
	cfg := LoadConfig()
	p := resolveTheme(cfg)

	fmt.Println(p.Primary.Sprint("NeoArc - Cross-Platform Alias Executor"))
	fmt.Println(`Usage:
  neoarc get <alias> [args...]       : Print the code for the alias
  neoarc run <alias> [args...]       : Run the alias command with args
  neoarc <alias> [args...]           : Run the alias command (shorthand)
  neoarc config <server-url>         : Set the NeoArc Web Server URL
  neoarc config-token <tok>          : Set the API authentication token
  neoarc config insecure             : Enable insecure TLS (skip certificate verify)
  neoarc config secure               : Disable insecure TLS (default, verify certs)
  neoarc config theme <name>         : Set the active color theme ('list' for all)
  neoarc config theme edit           : Open TUI theme picker
  neoarc edit                        : Open TUI configuration editor
  neoarc update                      : Check for updates and self-update the binary
  neoarc self-update                 : Alias for update
  neoarc version                     : Show the installed version
  neoarc help                        : Show this help menu
  neoarc completion <shell>          : Generate shell completion script (bash|zsh|powershell)

Flags:
  -v, --version                      : Show the installed version
  -h, --help                         : Show this help menu
  --install                          : Download and install NeoArc to ~/.config/neostore/neoarc/bin/
  --selfuninstall                    : Remove NeoArc config, cache, and binary from the system
  --config <path>                    : Use a custom config file path (before any command)
  --proxy <url>, -p <url>            : Use proxy for self-update download (after update command)
  --dry-run                          : Print the alias code without executing
  --yes                              : Skip execution confirmation prompt

Arguments after <alias> are passed through to the executed command.
  bash/sh    : accessible via $1, $2, $@
  powershell : accessible via $args[0], $args[1]
  python     : accessible via sys.argv[1], sys.argv[2]
  cmd        : accessible via %1, %2`)

	themeLabel := cfg.Theme
	if themeLabel == "" {
		themeLabel = config.DefaultTheme().Label
	} else if t, ok := config.FindTheme(themeLabel); ok {
		themeLabel = t.Label
	}
	fmt.Println()
	fmt.Printf("Version  : %s\n", p.Success.Sprint(resolveVersion(RepoPath)))
	fmt.Printf("Theme    : %s\n", p.Accent.Sprint(themeLabel))
	fmt.Printf("Config   : %s\n", ConfigPath())
}

func Run(args []string) int {
	migrateOldConfig()

	argIdx := 1

	if len(args) > 2 && args[1] == "--config" {
		customPath := args[2]
		if info, err := os.Stat(customPath); err == nil && !info.IsDir() {
			configPath = customPath
		} else {
			fmt.Fprintf(os.Stderr, "Error: config file not found: %s\n", customPath)
			return 1
		}
		argIdx = 3
	}

	shifted := append([]string{args[0]}, args[argIdx:]...)
	args = shifted

	if len(args) < 2 {
		ShowHelp()
		return 0
	}

	switch args[1] {
	case "help", "-h", "--help":
		ShowHelp()
		return 0
	case "-v", "--version", "version":
		fmt.Println("NeoArc version", resolveVersion(RepoPath))
		return 0
	case "edit":
		return editConfig()
	}

	if args[1] == "config" && len(args) >= 3 {
		sub := args[2]
		switch sub {
		case "edit":
			return editConfig()
		case "theme":
			if len(args) >= 4 && args[3] == "edit" {
				return editTheme()
			}
			if len(args) < 4 {
				fmt.Println("Usage: neoarc config theme <name>")
				fmt.Println("       neoarc config theme list")
				fmt.Println("       neoarc config theme edit")
				return 0
			}
			opt := args[3]
			if opt == "-h" || opt == "--help" {
				fmt.Println("Usage: neoarc config theme <name>")
				fmt.Println("       neoarc config theme list")
				fmt.Println("       neoarc config theme edit")
				return 0
			}
			if opt == "list" {
				fmt.Println("Available themes:")
				for _, t := range config.Themes {
					line := fmt.Sprintf("  %-30s %s", t.Name, t.Label)
					if t.Name == config.DefaultTheme().Name {
						line = config.ColorizeBold(line, t.Hex(config.RolePrimary))
					}
					fmt.Println(line)
				}
				return 0
			}
			if _, ok := config.FindTheme(opt); !ok {
				fmt.Fprintf(os.Stderr, "Error: unknown theme %q. Use 'neoarc config theme list' to see available themes.\n", opt)
				return 1
			}
			cfg := LoadConfig()
			cfg.Theme = opt
			SaveConfig(cfg)
			fmt.Println(sprintTheme(cfg, config.RoleSuccess, "Theme updated successfully!"))
			fmt.Printf("Active theme: %s\n", sprintTheme(cfg, config.RoleAccent, opt))
			return 0
		case "insecure":
			cfg := LoadConfig()
			cfg.InsecureTLS = true
			SaveConfig(cfg)
			fmt.Println("Insecure TLS enabled.")
			return 0
		case "secure":
			cfg := LoadConfig()
			cfg.InsecureTLS = false
			SaveConfig(cfg)
			fmt.Println("Insecure TLS disabled.")
			return 0
		default:
			if strings.HasPrefix(sub, "-") {
				fmt.Println("Usage: neoarc config <server-url>")
				fmt.Println("       neoarc config insecure")
				fmt.Println("       neoarc config secure")
				fmt.Println("       neoarc config theme <name>")
				fmt.Println("       neoarc config edit")
				return 0
			}
			if len(args) > 3 {
				fmt.Fprintf(os.Stderr, "Error: unexpected argument %q after %q\n", args[3], sub)
				fmt.Println("Usage: neoarc config <server-url>")
				return 1
			}
			cfg := LoadConfig()
			cfg.ServerURL = sub
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

	if args[1] == "update" || args[1] == "self-update" {
		proxyURL := ""
		for i := 2; i < len(args)-1; i++ {
			if args[i] == "--proxy" || args[i] == "-p" {
				proxyURL = args[i+1]
				break
			}
		}
		return selfUpdate(proxyURL)
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
	argIdx = 1

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
