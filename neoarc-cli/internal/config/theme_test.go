package config

import (
	"strings"
	"testing"
)

func TestDefaultTheme(t *testing.T) {
	tm := DefaultTheme()
	if tm.Name != "sunny_beach_day" {
		t.Fatalf("expected default theme 'sunny_beach_day', got %q", tm.Name)
	}
}

func TestFindTheme(t *testing.T) {
	tm, ok := FindTheme("dark")
	if !ok {
		t.Fatal("expected to find 'dark' theme")
	}
	if tm.Name != "dark" {
		t.Fatalf("expected name 'dark', got %q", tm.Name)
	}
	if len(tm.Colors) != 6 {
		t.Fatalf("expected 6 colors, got %d", len(tm.Colors))
	}
}

func TestFindThemeCaseInsensitive(t *testing.T) {
	tm, ok := FindTheme("DARK")
	if !ok {
		t.Fatal("expected case-insensitive find")
	}
	if tm.Name != "dark" {
		t.Fatalf("expected 'dark', got %q", tm.Name)
	}
}

func TestFindThemeNotFound(t *testing.T) {
	_, ok := FindTheme("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestThemeNames(t *testing.T) {
	names := ThemeNames()
	if len(names) != len(Themes) {
		t.Fatalf("expected %d names, got %d", len(Themes), len(names))
	}
	if names[0] != "dark" {
		t.Fatalf("expected first theme 'dark', got %q", names[0])
	}
}

func TestThemeLabels(t *testing.T) {
	labels := ThemeLabels()
	if len(labels) != len(Themes) {
		t.Fatalf("expected %d labels, got %d", len(Themes), len(labels))
	}
}

func TestResolveTheme(t *testing.T) {
	rc := ResolveTheme("dark")
	if rc.Primary == nil || rc.Success == nil || rc.Warning == nil || rc.Error == nil || rc.Accent == nil {
		t.Fatal("expected all role colors to be non-nil")
	}
}

func TestResolveThemeDefault(t *testing.T) {
	rc := ResolveTheme("")
	if rc.Primary == nil {
		t.Fatal("expected non-nil Primary for default theme")
	}
}

func TestThemeHex(t *testing.T) {
	tm := Theme{Name: "test", Colors: []string{"#ff0000", "#00ff00", "#0000ff"}}
	if h := tm.Hex(RolePrimary); h != "#ff0000" {
		t.Fatalf("expected #ff0000, got %q", h)
	}
	if h := tm.Hex(RoleSuccess); h != "#00ff00" {
		t.Fatalf("expected #00ff00, got %q", h)
	}
	if h := tm.Hex(RoleWarning); h != "#0000ff" {
		t.Fatalf("expected #0000ff, got %q", h)
	}
}

func TestThemeHexWraps(t *testing.T) {
	tm := Theme{Name: "test", Colors: []string{"#ff0000", "#00ff00"}}
	if h := tm.Hex(RoleAccent); h != "#ff0000" {
		t.Fatalf("expected #ff0000 (wrapped), got %q", h)
	}
}

func TestParseHex(t *testing.T) {
	r, g, b := parseHex("#ff0000")
	if r != 255 || g != 0 || b != 0 {
		t.Fatalf("expected 255,0,0 got %d,%d,%d", r, g, b)
	}
	r, g, b = parseHex("#00ff00")
	if r != 0 || g != 255 || b != 0 {
		t.Fatalf("expected 0,255,0 got %d,%d,%d", r, g, b)
	}
	r, g, b = parseHex("#0000ff")
	if r != 0 || g != 0 || b != 255 {
		t.Fatalf("expected 0,0,255 got %d,%d,%d", r, g, b)
	}
}

func TestParseHexInvalid(t *testing.T) {
	r, g, b := parseHex("nothex")
	if r != 255 || g != 255 || b != 255 {
		t.Fatalf("expected 255,255,255 fallback, got %d,%d,%d", r, g, b)
	}
}

func TestColorize(t *testing.T) {
	result := Colorize("hello", "#ff0000")
	if !strings.HasPrefix(result, "\x1b[") {
		t.Fatal("expected ANSI escape sequence")
	}
	if !strings.Contains(result, "hello") {
		t.Fatal("expected text to be included")
	}
	if !strings.HasSuffix(result, "\x1b[0m") {
		t.Fatal("expected reset sequence at end")
	}
}

func TestColorizeBold(t *testing.T) {
	result := ColorizeBold("bold text", "#00ff00")
	if !strings.Contains(result, "[1;") {
		t.Fatal("expected bold ANSI sequence")
	}
	if !strings.Contains(result, "bold text") {
		t.Fatal("expected bold text")
	}
}

func TestAllThemesHaveNames(t *testing.T) {
	for _, tm := range Themes {
		if tm.Name == "" {
			t.Fatal("all themes must have a name")
		}
		if tm.Label == "" {
			t.Fatalf("theme %q must have a label", tm.Name)
		}
		if len(tm.Colors) == 0 {
			t.Fatalf("theme %q must have at least one color", tm.Name)
		}
	}
}

func TestAllThemesUniqueNames(t *testing.T) {
	seen := make(map[string]bool)
	for _, tm := range Themes {
		if seen[tm.Name] {
			t.Fatalf("duplicate theme name: %q", tm.Name)
		}
		seen[tm.Name] = true
	}
}
