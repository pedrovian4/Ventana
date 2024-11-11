// localization/localization.go
package localization

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"ventana.com/pkg/config"
)

type Localization struct {
	messages  map[string]string
	appConfig *config.AppConfig
}

func NewLocalization(appConfig *config.AppConfig) *Localization {
	return &Localization{
		messages:  make(map[string]string),
		appConfig: appConfig,
	}
}

func (l *Localization) Initialize() error {
	lang := l.appConfig.Language
	if err := l.loadMessages(lang); err != nil {
		fmt.Println("Error loading language file:", err)
		if lang != "en" {
			fmt.Println("Attempting to load default language: English.")
			if fallbackErr := l.loadMessages("en"); fallbackErr != nil {
				return fmt.Errorf("error loading default language (English): %v", fallbackErr)
			}
		} else {
			return fmt.Errorf("critical error: default language could not be loaded: %v", err)
		}
	}
	return nil
}

func (l *Localization) loadMessages(lang string) error {
	filePath := filepath.Join(l.appConfig.MessageDirectory, lang+".json")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening language file (%s): %v", lang, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&l.messages); err != nil {
		return fmt.Errorf("error decoding language file (%s): %v", lang, err)
	}
	return nil
}

func (l *Localization) Msg(key string) string {
	msg, exists := l.messages[key]
	if !exists {
		return "Message not found: " + key
	}
	return msg
}

func (l *Localization) ChangeLanguage(newLang string) error {
	l.appConfig.Language = newLang
	if err := l.appConfig.Save(); err != nil {
		return err
	}
	return l.Initialize()
}
