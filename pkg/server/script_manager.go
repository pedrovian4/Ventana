package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"ventana.com/pkg/config"
	"ventana.com/pkg/history"
	"ventana.com/pkg/localization"
)

type ScriptManager struct {
	appConfig  *config.AppConfig
	historyMgr *history.HistoryManager
	localizer  *localization.Localization
}

func NewScriptManager(appConfig *config.AppConfig, historyMgr *history.HistoryManager, localizer *localization.Localization) *ScriptManager {
	return &ScriptManager{
		appConfig:  appConfig,
		historyMgr: historyMgr,
		localizer:  localizer,
	}
}

func (s *ScriptManager) RunScript() error {
	scripts, err := s.listScripts()
	if err != nil {
		return err
	}

	if len(scripts) == 0 {
		fmt.Println(config.Red(s.localizer.Msg("scripts_not_found")))
		return nil
	}

	script, err := s.selectScript(scripts)
	if err != nil {
		return err
	}

	servers, err := s.historyMgr.LoadHistory()
	if err != nil || len(servers) == 0 {
		fmt.Println(config.Red(s.localizer.Msg("no_saved_servers")))
		return nil
	}

	server, err := s.selectServer(servers)
	if err != nil {
		return err
	}

	err = s.executeScript(server, script)
	if err != nil {
		fmt.Println(config.Red(s.localizer.Msg("script_execution_failed")))
		return err
	}

	fmt.Println(config.LightGreen(s.localizer.Msg("script_executed")))
	return nil
}

func (s *ScriptManager) listScripts() ([]string, error) {
	scriptsDir := s.appConfig.ScriptsDirectory
	fmt.Println(scriptsDir)
	files, err := os.ReadDir(scriptsDir)
	if err != nil {
		return nil, fmt.Errorf(s.localizer.Msg("failed_to_read_scripts_dir"), err)
	}

	var scripts []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sh") {
			scripts = append(scripts, file.Name())
		}
	}

	return scripts, nil
}

func (s *ScriptManager) selectScript(scripts []string) (string, error) {
	prompt := promptui.Select{
		Label:     config.Purple(s.localizer.Msg("select_script")),
		Items:     scripts,
		Templates: s.getPromptTemplates(),
		Size:      len(scripts),
	}

	_, choice, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localizer.Msg("invalid_script_choice"))
		return "", err
	}

	return filepath.Join(s.appConfig.ScriptsDirectory, choice), nil
}

func (s *ScriptManager) selectServer(servers []history.Connection) (history.Connection, error) {
	formattedServers := s.formatConnections(servers)
	prompt := promptui.Select{
		Label:     config.Purple(s.localizer.Msg("select_server_for_script")),
		Items:     formattedServers,
		Templates: s.getPromptTemplates(),
		Size:      len(formattedServers),
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localizer.Msg("invalid_server_choice"))
		return history.Connection{}, err
	}

	return servers[index], nil
}

func (s *ScriptManager) executeScript(server history.Connection, scriptPath string) error {
	cmd := exec.Command(
		"ssh",
		"-i", server.PemFilePath,
		fmt.Sprintf("%s@%s", server.User, server.IP),
		"bash -s",
	)

	scriptFile, err := os.Open(scriptPath)
	if err != nil {
		return fmt.Errorf(s.localizer.Msg("failed_to_open_script"), err)
	}
	defer scriptFile.Close()

	cmd.Stdin = scriptFile
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (s *ScriptManager) formatConnections(servers []history.Connection) []string {
	formatted := make([]string, len(servers))
	for i, server := range servers {
		favMarker := ""
		if server.Favorite {
			favMarker = "⭐ "
		}
		formatted[i] = fmt.Sprintf("%s%s@%s", favMarker, config.Cyan(server.User), config.Yellow(server.IP))
	}
	return formatted
}

func (s *ScriptManager) getPromptTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Active:   "👉 {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "✅ {{ . | green }}",
	}
}
