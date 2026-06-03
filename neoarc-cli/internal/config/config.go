package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const ProjectName = "neoarc"

func HomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	h, _ := os.UserHomeDir()
	return h
}

func ConfigDir() string {
	return filepath.Join(HomeDir(), ".config", "neostore", ProjectName)
}

func EnsureConfigDir() error {
	return os.MkdirAll(ConfigDir(), 0755)
}

func ConfigFile(name string) string {
	return filepath.Join(ConfigDir(), name)
}

func LogFile(name string) string {
	return filepath.Join(ConfigDir(), name)
}

func SaveDir() string {
	dir := filepath.Join(HomeDir(), "Downloads", "neostore", ProjectName)
	os.MkdirAll(dir, 0755)
	return dir
}
