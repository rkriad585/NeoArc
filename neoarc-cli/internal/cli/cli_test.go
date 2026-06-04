package cli

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d := t.TempDir()
	os.Setenv("APPDATA", d)
	os.Setenv("HOME", d)
	os.Setenv("USERPROFILE", d)
	return d
}

func TestConfigLoadSave(t *testing.T) {
	d := tempDir(t)

	cfg := LoadConfig()
	if cfg.ServerURL != "http://localhost:59248" {
		t.Fatalf("expected default server URL, got %q", cfg.ServerURL)
	}

	cfg.ServerURL = "http://example.com:9999"
	cfg.APIToken = "test-token"
	cfg.InsecureTLS = true
	SaveConfig(cfg)

	cfg2 := LoadConfig()
	if cfg2.ServerURL != "http://example.com:9999" {
		t.Fatalf("expected custom server URL, got %q", cfg2.ServerURL)
	}
	if cfg2.APIToken != "test-token" {
		t.Fatalf("expected api token, got %q", cfg2.APIToken)
	}
	if !cfg2.InsecureTLS {
		t.Fatal("expected insecure TLS to be true")
	}

	path := filepath.Join(d, ".config", "neostore", "neoarc", "config.toml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("config.toml was not created")
	}
}

func TestCacheLoadSave(t *testing.T) {
	tempDir(t)

	cache := LoadCache()
	if cache.Entries == nil {
		t.Fatal("expected non-nil entries map")
	}
	if len(cache.Entries) != 0 {
		t.Fatal("expected empty cache")
	}

	now := time.Now().Unix()
	cache.Entries["test-alias"] = CacheEntry{
		Command:  "echo hi",
		ExecType: "bash",
		ETag:     `"abc123"`,
		Time:     now,
	}
	SaveCache(cache)

	cache2 := LoadCache()
	entry, ok := cache2.Entries["test-alias"]
	if !ok {
		t.Fatal("expected cached entry to exist after save")
	}
	if entry.Command != "echo hi" {
		t.Fatalf("expected command 'echo hi', got %q", entry.Command)
	}
	if entry.ETag != `"abc123"` {
		t.Fatalf("expected etag, got %q", entry.ETag)
	}
	if entry.Time != now {
		t.Fatalf("expected time %d, got %d", now, entry.Time)
	}
}

func TestCacheCorruptedFile(t *testing.T) {
	tempDir(t)

	path := filepath.Join(os.Getenv("USERPROFILE"), ".config", "neostore", "neoarc", "alias_cache.json")
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte("{broken json"), 0644)

	cache := LoadCache()
	if cache.Entries == nil {
		t.Fatal("expected non-nil entries even with corrupted cache")
	}
}

func TestCacheNilEntries(t *testing.T) {
	tempDir(t)

	path := filepath.Join(os.Getenv("USERPROFILE"), ".config", "neostore", "neoarc", "alias_cache.json")
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte(`{}`), 0644)

	cache := LoadCache()
	if cache.Entries == nil {
		t.Fatal("expected non-nil entries when json has no entries key")
	}
}

func TestTrustStoreLoadSave(t *testing.T) {
	tempDir(t)

	trusted := LoadTrusted()
	if trusted.Trusted == nil {
		t.Fatal("expected non-nil trusted map")
	}

	trusted.Trusted["test-alias"] = 1234567890
	SaveTrusted(trusted)

	trusted2 := LoadTrusted()
	ts, ok := trusted2.Trusted["test-alias"]
	if !ok {
		t.Fatal("expected trusted entry to exist")
	}
	if ts != 1234567890 {
		t.Fatalf("expected timestamp 1234567890, got %d", ts)
	}
}

func TestResolveAPIToken(t *testing.T) {
	cfg := Config{APIToken: "config-token"}
	APIToken = "compile-token"
	if tok := ResolveAPIToken(cfg); tok != "config-token" {
		t.Fatalf("expected config-token, got %q", tok)
	}

	cfg2 := Config{}
	if tok := ResolveAPIToken(cfg2); tok != "compile-token" {
		t.Fatalf("expected compile-token, got %q", tok)
	}

	APIToken = ""
	if tok := ResolveAPIToken(cfg2); tok != "" {
		t.Fatalf("expected empty token, got %q", tok)
	}

	APIToken = "compile-token"
}

func TestNewHTTPClient(t *testing.T) {
	cfgSecure := Config{InsecureTLS: false}
	client := NewHTTPClient(cfgSecure)
	if client.Timeout != 15*time.Second {
		t.Fatalf("expected 15s timeout, got %v", client.Timeout)
	}
	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("expected *http.Transport")
	}
	if tr.TLSClientConfig == nil || tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("expected secure TLS (nil or false)")
	}

	cfgInsecure := Config{InsecureTLS: true}
	client2 := NewHTTPClient(cfgInsecure)
	tr2, _ := client2.Transport.(*http.Transport)
	if tr2.TLSClientConfig == nil || !tr2.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("expected insecure TLS")
	}
}

func TestShowHelp(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	ShowHelp()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "NeoArc - Cross-Platform Alias Executor") {
		t.Fatal("help should contain title")
	}
	if !strings.Contains(output, "neoarc get") {
		t.Fatal("help should contain get command")
	}
	if !strings.Contains(output, "neoarc run") {
		t.Fatal("help should contain run command")
	}
	if !strings.Contains(output, "neoarc config") {
		t.Fatal("help should contain config command")
	}
	if !strings.Contains(output, "--dry-run") {
		t.Fatal("help should contain dry-run flag")
	}
	if !strings.Contains(output, "--yes") {
		t.Fatal("help should contain yes flag")
	}
	if !strings.Contains(output, "completion") {
		t.Fatal("help should contain completion command")
	}
	if !strings.Contains(output, "-v, --version") {
		t.Fatal("help should contain -v flag")
	}
	if !strings.Contains(output, "-h, --help") {
		t.Fatal("help should contain -h flag")
	}
	if !strings.Contains(output, "version") {
		t.Fatal("help should contain version command")
	}
}

func TestCacheTTL(t *testing.T) {
	tempDir(t)

	cache := LoadCache()
	now := time.Now().Unix()
	cache.Entries["ttl-test"] = CacheEntry{
		Command:  "old",
		ExecType: "bash",
		Time:     now - 31,
	}
	SaveCache(cache)

	cache2 := LoadCache()
	entry, ok := cache2.Entries["ttl-test"]
	if !ok {
		t.Fatal("expected cache entry")
	}
	if time.Now().Unix()-entry.Time < 30 {
		t.Fatal("expected entry to be older than 30s")
	}
}

func TestConfigPaths(t *testing.T) {
	d := tempDir(t)
	cfgDir := ConfigDir()
	expected := filepath.Join(d, ".config", "neostore", "neoarc")
	if cfgDir != expected {
		t.Fatalf("expected config dir %q, got %q", expected, cfgDir)
	}
	if ConfigPath() != filepath.Join(expected, "config.toml") {
		t.Fatal("config path mismatch")
	}
	if CachePath() != filepath.Join(expected, "alias_cache.json") {
		t.Fatal("cache path mismatch")
	}
	if TrustPath() != filepath.Join(expected, "trusted.json") {
		t.Fatal("trust path mismatch")
	}
}

func TestDirCreation(t *testing.T) {
	d := tempDir(t)
	cfgDir := filepath.Join(d, ".config", "neostore", "neoarc")
	os.RemoveAll(cfgDir)

	result := ConfigDir()
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		t.Fatal("ConfigDir should create the directory")
	}
	if result != cfgDir {
		t.Fatalf("expected %q, got %q", cfgDir, result)
	}
}

func TestRunHelp(t *testing.T) {
	code := Run([]string{"neoarc", "help"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRunNoArgs(t *testing.T) {
	code := Run([]string{"neoarc"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRunConfigThemeSet(t *testing.T) {
	_ = tempDir(t)
	code := Run([]string{"neoarc", "config", "theme", "dark"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	cfg := LoadConfig()
	if cfg.Theme != "dark" {
		t.Fatalf("expected theme 'dark', got %q", cfg.Theme)
	}
}

func TestRunConfigThemeNoArg(t *testing.T) {
	code := Run([]string{"neoarc", "config", "theme"})
	if code != 0 {
		t.Fatalf("expected exit code 0 (usage message), got %d", code)
	}
}

func TestRunConfigThemeHelp(t *testing.T) {
	code := Run([]string{"neoarc", "config", "theme", "--help"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRunConfigHelpFlag(t *testing.T) {
	code := Run([]string{"neoarc", "config", "--help"})
	if code != 0 {
		t.Fatalf("expected exit code 0 (usage), got %d", code)
	}
}

func TestRunConfigUnknownSubcommand(t *testing.T) {
	code := Run([]string{"neoarc", "config", "bogus-sub", "extra"})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRunConfigThemeUnknown(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w

	code := Run([]string{"neoarc", "config", "theme", "nonexistent"})

	w.Close()
	os.Stderr = old

	var buf strings.Builder
	io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "unknown theme") {
		t.Fatal("expected 'unknown theme' error on stderr")
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRunConfigThemeList(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	code := Run([]string{"neoarc", "config", "theme", "list"})

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "Available themes") {
		t.Fatal("expected 'Available themes' in output")
	}
	if !strings.Contains(buf.String(), "dark") {
		t.Fatal("expected 'dark' theme in list")
	}
	if !strings.Contains(buf.String(), "sunny_beach_day") {
		t.Fatal("expected 'sunny_beach_day' in list")
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestHelpContainsConfigTheme(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	ShowHelp()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "config theme") {
		t.Fatal("help should contain config theme command")
	}
	if !strings.Contains(output, "Theme") {
		t.Fatal("help should show theme info")
	}
}

func TestRunConfigServerURL(t *testing.T) {
	_ = tempDir(t)
	code := Run([]string{"neoarc", "config", "http://test:1234"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	cfg := LoadConfig()
	if cfg.ServerURL != "http://test:1234" {
		t.Fatalf("expected http://test:1234, got %q", cfg.ServerURL)
	}
}

func TestRunConfigToken(t *testing.T) {
	_ = tempDir(t)
	code := Run([]string{"neoarc", "config-token", "mytoken"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	cfg := LoadConfig()
	if cfg.APIToken != "mytoken" {
		t.Fatalf("expected mytoken, got %q", cfg.APIToken)
	}
}

func TestVersionFlag(t *testing.T) {
	code := Run([]string{"neoarc", "-v"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestVersionLongFlag(t *testing.T) {
	code := Run([]string{"neoarc", "--version"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestVersionCommand(t *testing.T) {
	code := Run([]string{"neoarc", "version"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestHelpShortFlag(t *testing.T) {
	code := Run([]string{"neoarc", "-h"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestHelpLongFlag(t *testing.T) {
	code := Run([]string{"neoarc", "--help"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestHelpContainsEditCommand(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	ShowHelp()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "neoarc edit") {
		t.Fatal("help should contain edit command")
	}
	if !strings.Contains(output, "neoarc config theme edit") {
		t.Fatal("help should contain theme edit command")
	}
}

func TestRunConfigInsecure(t *testing.T) {
	_ = tempDir(t)
	code := Run([]string{"neoarc", "config", "insecure"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	cfg := LoadConfig()
	if !cfg.InsecureTLS {
		t.Fatal("expected insecure TLS")
	}
}

func TestRunConfigSecure(t *testing.T) {
	_ = tempDir(t)
	Run([]string{"neoarc", "config", "insecure"})
	code := Run([]string{"neoarc", "config", "secure"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	cfg := LoadConfig()
	if cfg.InsecureTLS {
		t.Fatal("expected secure TLS")
	}
}

func TestRunUnknownFlag(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w

	code := Run([]string{"neoarc", "--bogus"})

	w.Close()
	os.Stderr = old

	var buf strings.Builder
	io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "Unknown flag") {
		t.Fatal("expected 'Unknown flag' on stderr")
	}
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRunGetMissingArg(t *testing.T) {
	code := Run([]string{"neoarc", "get"})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestRunRunMissingArg(t *testing.T) {
	code := Run([]string{"neoarc", "run"})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
}

func TestHelpContainsCompletion(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	ShowHelp()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "completion") {
		t.Fatal("help should contain completion command")
	}
	if !strings.Contains(output, "bash|zsh|powershell") {
		t.Fatal("help should list supported shells")
	}
	if !strings.Contains(output, "[args...]") {
		t.Fatal("help should show args in usage")
	}
}

func TestGenerateCompletionBash(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	code := GenerateCompletion("bash", []string{"myalias", "other"})

	w.Close()
	os.Stdout = old

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "_neoarc_completions") {
		t.Fatal("bash completion should define function")
	}
	if !strings.Contains(output, "complete -F") {
		t.Fatal("bash completion should register complete")
	}
	if !strings.Contains(output, "--dry-run") {
		t.Fatal("bash completion should include dry-run flag")
	}
	if !strings.Contains(output, "--yes") {
		t.Fatal("bash completion should include yes flag")
	}
	if !strings.Contains(output, "myalias") {
		t.Fatal("bash completion should include alias names")
	}
}

func TestGenerateCompletionZsh(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	code := GenerateCompletion("zsh", []string{})

	w.Close()
	os.Stdout = old

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "#compdef neoarc") {
		t.Fatal("zsh completion should start with #compdef")
	}
	if !strings.Contains(output, "_arguments") {
		t.Fatal("zsh completion should use _arguments")
	}
}

func TestGenerateCompletionPowerShell(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	code := GenerateCompletion("powershell", []string{"myalias"})

	w.Close()
	os.Stdout = old

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Register-ArgumentCompleter") {
		t.Fatal("powershell completion should register completer")
	}
	if !strings.Contains(output, "CompletionResult") {
		t.Fatal("powershell completion should use CompletionResult")
	}
	if !strings.Contains(output, "myalias") {
		t.Fatal("powershell completion should include alias names")
	}
}

func TestGenerateCompletionUnknownShell(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w

	code := GenerateCompletion("fish", []string{})

	w.Close()
	os.Stderr = old

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}

	var buf strings.Builder
	io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "Unknown shell") {
		t.Fatal("expected 'Unknown shell' on stderr")
	}
}

func TestRunCompletionCommand(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	code := Run([]string{"neoarc", "completion", "bash"})

	w.Close()
	os.Stdout = old

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	var buf strings.Builder
	io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "complete -F") {
		t.Fatal("completion command should output bash completion")
	}
}

func TestRunCompletionMissingArg(t *testing.T) {
	code := Run([]string{"neoarc", "completion"})
	if code != 0 {
		t.Fatalf("expected exit code 0 (help), got %d", code)
	}
}

func TestRunListAliasesAuthFail(t *testing.T) {
	_ = tempDir(t)
	Run([]string{"neoarc", "config", "http://localhost:1"})

	code := Run([]string{"neoarc", "_list_aliases"})
	if code != 0 {
		t.Fatalf("expected exit code 0 (graceful fallback), got %d", code)
	}
}

func TestRunWithAliasArgs(t *testing.T) {
	_ = tempDir(t)
	Run([]string{"neoarc", "config", "http://localhost:1"})
	Run([]string{"neoarc", "config-token", "test-token"})

	code := Run([]string{"neoarc", "get", "myalias", "arg1", "arg2"})
	if code != 1 {
		t.Fatalf("expected exit code 1 (connection refused), got %d", code)
	}
}

func TestRunShorthandWithAliasArgs(t *testing.T) {
	_ = tempDir(t)
	Run([]string{"neoarc", "config", "http://localhost:1"})
	Run([]string{"neoarc", "config-token", "test-token"})

	code := Run([]string{"neoarc", "myalias", "arg1", "arg2"})
	if code != 1 {
		t.Fatalf("expected exit code 1 (connection refused), got %d", code)
	}
}

func TestFetchAliasesConnectionError(t *testing.T) {
	_ = tempDir(t)
	Run([]string{"neoarc", "config", "http://localhost:1"})
	Run([]string{"neoarc", "config-token", "test-token"})

	aliases, err := FetchAliases()
	if err == nil {
		t.Fatal("expected connection error")
	}
	if aliases != nil {
		t.Fatal("expected nil aliases on error")
	}
}

func TestFetchAliasesAuthFail(t *testing.T) {
	_ = tempDir(t)
	Run([]string{"neoarc", "config", "http://localhost:1"})

	aliases, err := FetchAliases()
	if err == nil {
		t.Fatal("expected connection error when no server is running")
	}
	if aliases != nil {
		t.Fatal("expected nil aliases on connection error")
	}
}

func TestRunSelfUninstall(t *testing.T) {
	d := tempDir(t)

	cfgDir := filepath.Join(d, ".config", "neostore", "neoarc")
	os.MkdirAll(cfgDir, 0755)
	os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte("server_url = \"http://test:1234\"\n"), 0644)

	code := Run([]string{"neoarc", "--selfuninstall"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	if _, err := os.Stat(cfgDir); !os.IsNotExist(err) {
		t.Fatal("expected config directory to be removed")
	}
}

func TestInstallDir(t *testing.T) {
	d := installDir()
	if !strings.Contains(d, ".config") || !strings.Contains(d, "neostore") {
		t.Fatalf("expected install dir to contain .config/neostore, got %q", d)
	}
}

func TestRunSelfInstall(t *testing.T) {
	code := Run([]string{"neoarc", "--install"})
	if code != 0 {
		t.Fatalf("expected exit code 0 (fallback copy), got %d", code)
	}

	targetDir := installDir()
	binName := "neoarc"
	if runtime.GOOS == "windows" {
		binName = "neoarc.exe"
	}
	targetPath := filepath.Join(targetDir, binName)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Fatal("expected binary to be installed")
	}
}

func TestHelpContainsInstall(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	ShowHelp()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "--install") {
		t.Fatal("help should contain --install flag")
	}
}

func TestCopySelf(t *testing.T) {
	exe, err := os.Executable()
	if err != nil || exe == "" {
		t.Skip("binary path not available")
	}
	if _, err := os.Stat(exe); os.IsNotExist(err) {
		t.Skip("test binary has been removed")
	}

	d := t.TempDir()
	target := filepath.Join(d, "neoarc.copy")

	code := copySelf(target)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Fatal("expected binary to be copied")
	}

	info, _ := os.Stat(target)
	if info.Size() == 0 {
		t.Fatal("expected non-empty copied binary")
	}
}

func TestHelpContainsSelfUninstall(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	ShowHelp()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "--selfuninstall") {
		t.Fatal("help should contain --selfuninstall flag")
	}
	if !strings.Contains(output, "-u") {
		t.Fatal("help should contain -u flag")
	}
	if !strings.Contains(output, "--uninstall") {
		t.Fatal("help should contain --uninstall flag")
	}
}

func TestHelpContainsUpdate(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	ShowHelp()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "update") {
		t.Fatal("help should contain update command")
	}
}

func TestRunUpdate(t *testing.T) {
	_ = tempDir(t)
	Version = "v0.0.1"
	defer func() { Version = "" }()

	code := Run([]string{"neoarc", "update"})
	if code != 0 {
		t.Fatalf("expected exit code 0 (update succeeded), got %d", code)
	}
}

func TestRunUpdateAlreadyLatest(t *testing.T) {
	_ = tempDir(t)

	code := Run([]string{"neoarc", "update"})
	if code != 0 {
		t.Fatalf("expected exit code 0 (already up to date), got %d", code)
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"v1.0.0", "v1.0.0", 0},
		{"v1.0.0", "v1.0.1", -1},
		{"v1.0.1", "v1.0.0", 1},
		{"v1.0.0", "v2.0.0", -1},
		{"v2.0.0", "v1.0.0", 1},
		{"v1.0", "v1.0.0", 0},
		{"v1.0.0", "v1.0", 0},
		{"1.0.0", "v1.0.0", 0},
		{"v1.0.0", "1.0.0", 0},
	}

	for _, tt := range tests {
		got := compareVersions(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestResolveDownloadName(t *testing.T) {
	name := resolveDownloadName()
	if name == "" {
		t.Fatal("expected non-empty download name")
	}
	if !strings.Contains(name, "neoarc-") {
		t.Fatal("expected download name to contain neoarc-")
	}
}

func TestMigrationFromOldPaths(t *testing.T) {
	d := tempDir(t)
	oldDir := filepath.Join(d, "neoarc")
	os.MkdirAll(oldDir, 0755)
	oldCfg := filepath.Join(oldDir, "config.json")
	os.WriteFile(oldCfg, []byte(`{"server_url":"http://migrated:1234"}`), 0644)

	cfg := LoadConfig()
	if cfg.ServerURL != "http://migrated:1234" {
		t.Fatalf("expected migrated config, got %q", cfg.ServerURL)
	}

	newDir := filepath.Join(d, ".config", "neostore", "neoarc")
	newCfg := filepath.Join(newDir, "config.toml")
	if _, err := os.Stat(newCfg); os.IsNotExist(err) {
		t.Fatal("expected config to be migrated to new directory as config.toml")
	}
}
