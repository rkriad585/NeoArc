package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "neoarc.exe")

	cmd := exec.Command("go", "build", "-o", bin, "./cmd/neoarc")
	cmd.Dir = filepath.Join("..")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	return bin
}

func TestBinaryHelp(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "help")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(string(out), "NeoArc - Cross-Platform Alias Executor") {
		t.Fatal("help output should contain title")
	}
}

func TestBinaryNoArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, _ := cmd.Output()
	if !strings.Contains(string(out), "NeoArc - Cross-Platform Alias Executor") {
		t.Fatal("should show help with no args")
	}
}

func TestBinaryConfigHelp(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "help")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	output := string(out)
	if !strings.Contains(output, "neoarc get") {
		t.Fatal("help should list 'get' command")
	}
	if !strings.Contains(output, "--dry-run") {
		t.Fatal("help should list --dry-run flag")
	}
}
