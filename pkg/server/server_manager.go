// server/server.go
package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
	"ventana.com/pkg/config"
	"ventana.com/pkg/history"
	"ventana.com/pkg/localization"
)

type ServerManager struct {
	localization *localization.Localization
	appConfig    *config.AppConfig
	historyMgr   *history.HistoryManager
	syntaxMgr    *SyntaxManager
}

func NewServerManager(localization *localization.Localization, appConfig *config.AppConfig) *ServerManager {
	historyMgr := history.NewHistoryManager(appConfig)
	server := &ServerManager{
		localization: localization,
		appConfig:    appConfig,
		historyMgr:   historyMgr,
		syntaxMgr:    NewSyntaxManager(appConfig, historyMgr, localization),
	}
	return server
}

func (s *ServerManager) ConnectToNewServer() {
	ip := s.promptUserInput(s.localization.Msg("enter_ip"))
	if ip == "" {
		return
	}

	user := s.promptUserInput(s.localization.Msg("enter_user"))
	if user == "" {
		return
	}

	pemFile := s.choosePemFile()
	if pemFile == "" {
		return
	}

	attempts := s.promptIntInput(s.localization.Msg("retries_prompt"))
	delay := 0
	if attempts > 0 {
		delay = s.promptIntInput(s.localization.Msg("delay_prompt"))
	}

	server := history.Connection{
		User:         user,
		IP:           ip,
		PemFilePath:  pemFile,
		LastUsed:     time.Now().Format(time.RFC3339),
		Favorite:     false,
		DelaySeconds: delay,
		RetryCounts:  attempts,
	}

	historyList, _ := s.historyMgr.LoadHistory()
	if err := s.historyMgr.SaveHistory(append(historyList, server)); err != nil {
		fmt.Println(s.localization.Msg("error_saving_history"), err)
		return
	}

	s.connectWithRetries(server)
}

func (s *ServerManager) connectWithRetries(server history.Connection) {
	for i := 0; i <= server.RetryCounts; i++ {
		fmt.Printf(s.localization.Msg("connecting"), i+1, server.User, server.IP, filepath.Base(server.PemFilePath))
		if s.connect(server) == nil {
			fmt.Println(s.localization.Msg("connected"))
			return
		}
		if i < server.RetryCounts {
			fmt.Printf(s.localization.Msg("retrying"), server.DelaySeconds)
			time.Sleep(time.Duration(server.DelaySeconds) * time.Second)
		} else {
			fmt.Println(s.localization.Msg("failed"))
		}
	}
}

func (s *ServerManager) connect(server history.Connection) error {
	os.Setenv("TERM", "xterm")
	cmd := exec.Command("ssh", "-i", server.PemFilePath, fmt.Sprintf("%s@%s", server.User, server.IP))
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func (s *ServerManager) SelectServerFromHistory() {
	historyList, err := s.historyMgr.LoadHistory()
	if err != nil || len(historyList) == 0 {
		fmt.Println(s.localization.Msg("no_saved_servers"))
		return
	}

	prompt := promptui.Select{
		Label:     s.localization.Msg("choose_server"),
		Items:     s.formatConnections(historyList),
		Templates: s.getPromptTemplates(),
		Size:      len(historyList),
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localization.Msg("invalid_choice"))
		return
	}

	server := historyList[index]
	s.connectWithRetries(server)
}

func (s *ServerManager) RemoveServerFromHistory() {
	historyList, err := s.historyMgr.LoadHistory()
	if err != nil || len(historyList) == 0 {
		fmt.Println(s.localization.Msg("no_saved_servers"))
		return
	}

	prompt := promptui.Select{
		Label:     s.localization.Msg("choose_remove_server"),
		Items:     s.formatConnections(historyList),
		Templates: s.getPromptTemplates(),
		Size:      len(historyList),
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localization.Msg("invalid_choice"))
		return
	}

	historyList = append(historyList[:index], historyList[index+1:]...)
	if err := s.historyMgr.SaveHistory(historyList); err != nil {
		fmt.Println(s.localization.Msg("error_saving_history"), err)
		return
	}
	fmt.Println(s.localization.Msg("server_removed"))
}

func (s *ServerManager) UpdateServerInHistory() {
	historyList, err := s.historyMgr.LoadHistory()
	if err != nil || len(historyList) == 0 {
		fmt.Println(s.localization.Msg("no_saved_servers"))
		return
	}

	prompt := promptui.Select{
		Label:     s.localization.Msg("choose_update_server"),
		Items:     s.formatConnections(historyList),
		Templates: s.getPromptTemplates(),
		Size:      len(historyList),
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localization.Msg("invalid_choice"))
		return
	}

	server := &historyList[index]
	server.IP = s.promptOptionalInput(s.localization.Msg("update_ip"), server.IP)
	server.User = s.promptOptionalInput(s.localization.Msg("update_user"), server.User)
	server.PemFilePath = s.promptOptionalInput(s.localization.Msg("update_pem"), server.PemFilePath)

	if err := s.historyMgr.SaveHistory(historyList); err != nil {
		fmt.Println(s.localization.Msg("error_saving_history"), err)
		return
	}
	fmt.Println(s.localization.Msg("server_updated"))
}

func (s *ServerManager) ToggleFavoriteServer() {
	historyList, err := s.historyMgr.LoadHistory()
	if err != nil || len(historyList) == 0 {
		fmt.Println(s.localization.Msg("no_saved_servers"))
		return
	}

	prompt := promptui.Select{
		Label:     s.localization.Msg("choose_toggle_favorite"),
		Items:     s.formatConnections(historyList),
		Templates: s.getPromptTemplates(),
		Size:      len(historyList),
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localization.Msg("invalid_choice"))
		return
	}

	historyList[index].Favorite = !historyList[index].Favorite
	if err := s.historyMgr.SaveHistory(historyList); err != nil {
		fmt.Println(s.localization.Msg("error_saving_history"), err)
		return
	}
	fmt.Println(s.localization.Msg("favorite_toggled"))
}

func (s *ServerManager) promptUserInput(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}

	input, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localization.Msg("invalid_input"))
		return ""
	}

	return input
}

func (s *ServerManager) promptIntInput(label string) int {
	for {
		prompt := promptui.Prompt{
			Label: label,
		}

		input, err := prompt.Run()
		if err != nil {
			fmt.Println(s.localization.Msg("invalid_input"))
			return 0
		}

		intValue, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(s.localization.Msg("invalid_integer"))
			continue
		}

		return intValue
	}
}

func (s *ServerManager) promptOptionalInput(label, currentValue string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: currentValue,
	}

	input, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localization.Msg("invalid_input"))
		return currentValue
	}

	if input == "" {
		return currentValue
	}

	return input
}

func (s *ServerManager) choosePemFile() string {
	pemFiles := s.listPemFiles()
	if len(pemFiles) == 0 {
		fmt.Println(s.localization.Msg("no_pem_files"))
		return ""
	}

	prompt := promptui.Select{
		Label:     s.localization.Msg("select_pem_file"),
		Items:     s.formatPemFiles(pemFiles),
		Templates: s.getPromptTemplates(),
		Size:      len(pemFiles),
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println(s.localization.Msg("invalid_choice"))
		return ""
	}

	return pemFiles[index]
}

func (s *ServerManager) listPemFiles() []string {
	pemDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	files, err := filepath.Glob(filepath.Join(pemDir, "*.pem"))
	if err != nil {
		fmt.Println(s.localization.Msg("pem_error"), err)
		return []string{}
	}
	return files
}

func (s *ServerManager) formatConnections(connections []history.Connection) []string {
	formatted := make([]string, len(connections))
	for i, conn := range connections {
		favorite := ""
		if conn.Favorite {
			favorite = "* "
		}
		formatted[i] = fmt.Sprintf("%s%s@%s", favorite, conn.User, conn.IP)
	}
	return formatted
}

func (s *ServerManager) formatPemFiles(pemFiles []string) []string {
	formatted := make([]string, len(pemFiles))
	for i, file := range pemFiles {
		formatted[i] = filepath.Base(file)
	}
	return formatted
}

func (s *ServerManager) getPromptTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Active:   "👉 {{ . | cyan }}",
		Inactive: "  {{ . }}",
		Selected: "✅ {{ . | green }}",
	}
}
