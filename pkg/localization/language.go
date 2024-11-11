package localization

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"ventana.com/pkg/config"
)

type Language struct {
	settings     *config.AppConfig
	localization *Localization
}

func NewLanguage(settings *config.AppConfig, localization *Localization) *Language {
	return &Language{
		settings:     settings,
		localization: localization,
	}
}

func (l *Language) ChangeLanguage() {
	fmt.Println(config.LightGreen(l.localization.Msg("change_language")))

	languages := []string{
		"1️⃣ 🇺🇸 English",
		"2️⃣ 🇧🇷 Português",
		"3️⃣ 🇷🇺 Русский",
	}

	prompt := promptui.Select{
		Label: config.LightGreen(l.localization.Msg("enter_choice")),
		Items: languages,
		Templates: &promptui.SelectTemplates{
			Active:   "👉 {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✅ {{ . | green }}",
		},
		Size: len(languages),
	}

	_, langChoice, err := prompt.Run()
	if err != nil {
		fmt.Println(config.Red(l.localization.Msg("invalid_choice")))
		return
	}
	lang := "en"
	switch langChoice {
	case languages[1]:
		lang = "pt"
	case languages[2]:
		lang = "ru"
	}
	l.settings.Update(map[string]interface{}{"Language": lang})

}
