// main.go
package main

import (
	"fmt"
	"os"

	"ventana.com/pkg/arts"
	"ventana.com/pkg/config"
	"ventana.com/pkg/localization"
	"ventana.com/pkg/menu"
)

func main() {
	settings := config.NewAppConfig()
	if err := settings.Load(); err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	loc := localization.NewLocalization(settings)
	if err := loc.Initialize(); err != nil {
		fmt.Println("Error initializing localization:", err)
		os.Exit(1)
	}

	arts.DisplayWelcomeMessage()
	lang := localization.NewLanguage(settings, loc)
	m := menu.NewMenu(loc, lang, settings)

	for {
		m.Display()
		loc.ChangeLanguage(settings.Language)
	}
}
