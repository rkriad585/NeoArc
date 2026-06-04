package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"neoarc/internal/banner"
	"neoarc/internal/config"
)

func editConfig() int {
	banner.Print()
	cfg := LoadConfig()

	themeNames := config.ThemeNames()
	themeLabels := config.ThemeLabels()
	themeOptions := make([]huh.Option[string], len(themeNames))
	currentThemeIdx := 0
	for i, name := range themeNames {
		label := themeLabels[i]
		themeOptions[i] = huh.NewOption(label, name)
		if name == cfg.Theme || (cfg.Theme == "" && name == config.DefaultTheme().Name) {
			currentThemeIdx = i
		}
	}

	var serverURL, apiToken string
	var insecureTLS bool
	var theme string

	serverURL = cfg.ServerURL
	apiToken = cfg.APIToken
	insecureTLS = cfg.InsecureTLS
	theme = cfg.Theme
	if theme == "" {
		theme = themeNames[currentThemeIdx]
	}

	themeSelect := huh.NewSelect[string]().
		Title("Theme").
		Options(themeOptions...).
		Value(&theme)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server URL").
				Description("NeoArc Web Server address").
				Value(&serverURL),

			huh.NewInput().
				Title("API Token").
				Description("Authentication token (leave blank to keep current)").
				Value(&apiToken),

			huh.NewConfirm().
				Title("Insecure TLS").
				Description("Skip TLS certificate verification").
				Value(&insecureTLS),
		),

		huh.NewGroup(themeSelect),
	)

	err := form.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: form cancelled: %v\n", err)
		return 1
	}

	cfg.ServerURL = serverURL
	cfg.APIToken = apiToken
	cfg.InsecureTLS = insecureTLS
	cfg.Theme = theme
	SaveConfig(cfg)
	fmt.Println(config.Colorize("\n✓ Configuration saved!", config.DefaultTheme().Hex(config.RoleSuccess)))
	return 0
}

func editTheme() int {
	banner.Print()
	cfg := LoadConfig()

	themeNames := config.ThemeNames()
	themeLabels := config.ThemeLabels()
	themeOptions := make([]huh.Option[string], len(themeNames))
	currentThemeIdx := 0
	for i, name := range themeNames {
		label := themeLabels[i]
		themeOptions[i] = huh.NewOption(label, name)
		if name == cfg.Theme || (cfg.Theme == "" && name == config.DefaultTheme().Name) {
			currentThemeIdx = i
		}
	}

	theme := cfg.Theme
	if theme == "" {
		theme = themeNames[currentThemeIdx]
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Theme").
				Description("Select a color theme for the CLI").
				Options(themeOptions...).
				Value(&theme),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: form cancelled: %v\n", err)
		return 1
	}

	cfg.Theme = theme
	SaveConfig(cfg)
	fmt.Println(config.Colorize("\n✓ Theme updated successfully!", config.DefaultTheme().Hex(config.RoleSuccess)))
	return 0
}
