// menu/menu.go
package menu

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"ventana.com/pkg/arts"
	"ventana.com/pkg/config"
	"ventana.com/pkg/history"
	"ventana.com/pkg/localization"
	"ventana.com/pkg/server"
)

type Menu struct {
	localization *localization.Localization
	language     *localization.Language
	appConfig    *config.AppConfig
	serverMgr    *server.ServerManager
	scriptMgr    *server.ScriptManager
}

func NewMenu(localization *localization.Localization, language *localization.Language, appConfig *config.AppConfig) *Menu {
	serverMgr := server.NewServerManager(localization, appConfig)
	return &Menu{
		localization: localization,
		language:     language,
		appConfig:    appConfig,
		serverMgr:    serverMgr,
		scriptMgr:    server.NewScriptManager(appConfig, history.NewHistoryManager(appConfig), localization),
	}
}

func (m *Menu) Display() {
	menuOptions := []string{
		config.Purple(m.localization.Msg("connect_new")),
		config.DarkBlue(m.localization.Msg("my_servers")),
		config.Purple(m.localization.Msg("remove_server")),
		config.DarkBlue(m.localization.Msg("update_server")),
		config.Purple(m.localization.Msg("toggle_favorite")),
		config.DarkBlue(m.localization.Msg("change_language")),
		config.Purple(m.localization.Msg("run_script")),
		config.Red(m.localization.Msg("exit_option")),
	}

	prompt := promptui.Select{
		Label: m.localization.Msg("enter_choice"),
		Items: menuOptions,
		Templates: &promptui.SelectTemplates{
			Active:   "👉 {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✅ {{ . | green }}",
		},
		Size: len(menuOptions),
	}

	_, choice, err := prompt.Run()
	if err != nil {
		fmt.Println("Invalid choice.")
		return
	}

	switch choice {
	case menuOptions[0]:
		m.serverMgr.ConnectToNewServer()
	case menuOptions[1]:
		m.serverMgr.SelectServerFromHistory()
	case menuOptions[2]:
		m.serverMgr.RemoveServerFromHistory()
	case menuOptions[3]:
		m.serverMgr.UpdateServerInHistory()
	case menuOptions[4]:
		m.serverMgr.ToggleFavoriteServer()
	case menuOptions[5]:
		m.language.ChangeLanguage()
	case menuOptions[6]:
		m.scriptMgr.RunScript()
	case menuOptions[7]:
		fmt.Print("\033[H\033[2J")
		arts.DisplayWelcomeMessage()
		fmt.Println(m.localization.Msg("good_bye"))
		os.Exit(0)
	default:
		fmt.Println("Invalid choice.")
	}
	fmt.Print("\033[H\033[2J")
}
