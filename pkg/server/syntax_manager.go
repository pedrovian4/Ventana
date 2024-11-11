package server

/**
Em desenvolvimento
esse aquivo deve forncer autocomplete no ssh
além disso syntax highlight
**/
import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"ventana.com/pkg/config"
	"ventana.com/pkg/history"
	"ventana.com/pkg/localization"
)

type SyntaxManager struct {
	appConfig  *config.AppConfig
	historyMgr *history.HistoryManager
	localizer  *localization.Localization
}

func NewSyntaxManager(appConfig *config.AppConfig, historyMgr *history.HistoryManager, localizer *localization.Localization) *SyntaxManager {
	return &SyntaxManager{
		appConfig:  appConfig,
		historyMgr: historyMgr,
		localizer:  localizer,
	}
}

func (s *SyntaxManager) ConfigureAutocomplete() error {
	if _, err := os.Stat(s.appConfig.AutocompleteFile); os.IsNotExist(err) {
		if err := s.generateAutocompleteScript(); err != nil {
			return fmt.Errorf("failed to generate autocomplete script: %v", err)
		}
		fmt.Println(s.localizer.Msg("autocomplete_configured"))
	}
	return s.loadAutocompleteScript()
}

func (s *SyntaxManager) generateAutocompleteScript() error {
	history, err := s.historyMgr.LoadHistory()
	if err != nil {
		return fmt.Errorf("failed to load history for autocomplete: %v", err)
	}

	script := `
	_ssh_autocomplete() {
		local cur prev opts
		COMPREPLY=()
		cur="${COMP_WORDS[COMP_CWORD]}"
		prev="${COMP_WORDS[COMP_CWORD-1]}"

		opts="`

	for _, server := range history {
		serverOption := fmt.Sprintf("%s@%s", server.User, server.IP)
		script += serverOption + " "
	}

	script += `"

		COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
		return 0
	}

	complete -F _ssh_autocomplete ssh
	`

	file, err := os.Create(s.appConfig.AutocompleteFile)
	if err != nil {
		return fmt.Errorf("failed to write autocomplete script: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(script)
	if err != nil {
		return fmt.Errorf("failed to save autocomplete script: %v", err)
	}

	return nil
}

func (s *SyntaxManager) loadAutocompleteScript() error {
	cmd := exec.Command("bash", "-c", "source "+s.appConfig.AutocompleteFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to load autocomplete script: %v", err)
	}
	return nil
}

func (s *SyntaxManager) ConnectWithSyntaxHighlight(server history.Connection) error {
	os.Setenv("TERM", "xterm-256color")

	cmd := exec.Command("ssh", "-t", "-i", server.PemFilePath, fmt.Sprintf("%s@%s", server.User, server.IP))

	sshOutput, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create SSH output pipe: %v", err)
	}

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start SSH command: %v", err)
	}

	if _, err := exec.LookPath("bat"); err == nil {
		batCmd := exec.Command("bat", "-p", "--theme=ansi", "--color=always")
		batCmd.Stdin = sshOutput
		batCmd.Stdout = os.Stdout
		batCmd.Stderr = os.Stderr
		if err := batCmd.Run(); err != nil {
			return fmt.Errorf("error executing bat: %v", err)
		}
	} else {
		fmt.Println(s.localizer.Msg("syntax_highlight_limited"))
		if _, err := io.Copy(os.Stdout, sshOutput); err != nil {
			return fmt.Errorf("error copying SSH output: %v", err)
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for SSH command to finish: %v", err)
	}

	return nil
}
