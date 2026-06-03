package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Role int

const (
	RolePrimary Role = iota
	RoleSuccess
	RoleWarning
	RoleError
	RoleAccent
)

type Theme struct {
	Name   string
	Label  string
	Colors []string
}

func (t Theme) Hex(role Role) string {
	idx := int(role)
	if idx >= len(t.Colors) {
		idx %= len(t.Colors)
	}
	return t.Colors[idx]
}

func (t Theme) Color(role Role) *color.Color {
	r, g, b := parseHex(t.Hex(role))
	return color.RGB(int(r), int(g), int(b)).Add(color.Bold)
}

var defaultTheme = "sunny_beach_day"

var Themes = []Theme{
	{
		Name: "dark", Label: "Dark Theme",
		Colors: []string{"#0f172a", "#111827", "#1e293b", "#334155", "#64748b", "#e2e8f0"},
	},
	{
		Name: "light", Label: "Light Theme",
		Colors: []string{"#ffffff", "#f8fafc", "#e2e8f0", "#cbd5e1", "#475569", "#0f172a"},
	},
	{
		Name: "sunny_beach_day", Label: "Sunny Beach Day (Default)",
		Colors: []string{"#264653", "#2a9d8f", "#e9c46a", "#f4a261", "#e76f51"},
	},
	{
		Name: "olive_garden_feast", Label: "Olive Garden Feast",
		Colors: []string{"#606c38", "#283618", "#fefae0", "#dda15e", "#bc6c25"},
	},
	{
		Name: "summer_ocean_breeze", Label: "Summer Ocean Breeze",
		Colors: []string{"#e63946", "#f1faee", "#a8dadc", "#457b9d", "#1d3557"},
	},
	{
		Name: "refreshing_summer_fun", Label: "Refreshing Summer Fun",
		Colors: []string{"#8ecae6", "#219ebc", "#023047", "#ffb703", "#fb8500"},
	},
	{
		Name: "black_gold_elegance", Label: "Black & Gold Elegance",
		Colors: []string{"#000000", "#14213d", "#fca311", "#e5e5e5", "#ffffff"},
	},
	{
		Name: "vibrant_color_fiesta", Label: "Vibrant Color Fiesta",
		Colors: []string{"#ffbe0b", "#fb5607", "#ff006e", "#8338ec", "#3a86ff"},
	},
	{
		Name: "light_steel", Label: "Light Steel",
		Colors: []string{"#f8f9fa", "#e9ecef", "#dee2e6", "#ced4da", "#adb5bd", "#6c757d", "#495057", "#343a40", "#212529"},
	},
	{
		Name: "golden_twilight", Label: "Golden Twilight",
		Colors: []string{"#000814", "#001d3d", "#003566", "#ffc300", "#ffd60a"},
	},
	{
		Name: "deep_sea", Label: "Deep Sea",
		Colors: []string{"#0d1b2a", "#1b263b", "#415a77", "#778da9", "#e0e1dd"},
	},
	{
		Name: "bright_green", Label: "Bright Green",
		Colors: []string{"#004b23", "#006400", "#007200", "#008000", "#38b000", "#70e000", "#9ef01a", "#ccff33"},
	},
	{
		Name: "vivid_nightfall", Label: "Vivid Nightfall",
		Colors: []string{"#10002b", "#240046", "#3c096c", "#5a189a", "#7b2cbf", "#9d4edd", "#c77dff", "#e0aaff"},
	},
}

func FindTheme(name string) (Theme, bool) {
	for _, t := range Themes {
		if strings.EqualFold(t.Name, name) {
			return t, true
		}
	}
	return Theme{}, false
}

func ThemeNames() []string {
	names := make([]string, len(Themes))
	for i, t := range Themes {
		names[i] = t.Name
	}
	return names
}

func ThemeLabels() []string {
	labels := make([]string, len(Themes))
	for i, t := range Themes {
		labels[i] = t.Label
	}
	return labels
}

func DefaultTheme() Theme {
	t, ok := FindTheme(defaultTheme)
	if !ok {
		return Themes[0]
	}
	return t
}

type RoleColors struct {
	Primary, Success, Warning, Error, Accent *color.Color
}

func ResolveTheme(name string) RoleColors {
	t, ok := FindTheme(name)
	if !ok {
		t = DefaultTheme()
	}
	return RoleColors{
		Primary: t.Color(RolePrimary),
		Success: t.Color(RoleSuccess),
		Warning: t.Color(RoleWarning),
		Error:   t.Color(RoleError),
		Accent:  t.Color(RoleAccent),
	}
}

func parseHex(hex string) (uint8, uint8, uint8) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255
	}
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return 255, 255, 255
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return 255, 255, 255
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return 255, 255, 255
	}
	return uint8(r), uint8(g), uint8(b)
}

func Colorize(text, hex string) string {
	r, g, b := parseHex(hex)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}

func ColorizeBold(text, hex string) string {
	r, g, b := parseHex(hex)
	return fmt.Sprintf("\x1b[1;38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}
