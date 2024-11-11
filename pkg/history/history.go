package history

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"ventana.com/pkg/config"
)

type Connection struct {
	User         string
	IP           string
	PemFilePath  string
	LastUsed     string
	Favorite     bool
	DelaySeconds int
	RetryCounts  int
}

type HistoryManager struct {
	appConfig *config.AppConfig
}

func NewHistoryManager(appConfig *config.AppConfig) *HistoryManager {
	return &HistoryManager{
		appConfig: appConfig,
	}
}

func (h *HistoryManager) LoadHistory() ([]Connection, error) {
	var history []Connection

	file, err := os.Open(h.appConfig.HistoryFile)
	if os.IsNotExist(err) {
		fmt.Println("History file does not exist, creating a new one.")
		if err := h.SaveHistory(history); err != nil {
			return nil, fmt.Errorf("error creating new history file: %v", err)
		}
		return history, nil
	} else if err != nil {
		return nil, fmt.Errorf("error opening history file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading history file: %v", err)
	}

	for _, record := range records {
		if len(record) < 5 {
			continue
		}
		favorite := record[4] == "true"
		retryCounts, _ := strconv.Atoi(record[5])
		delaySeconds, _ := strconv.Atoi(record[6])
		history = append(history, Connection{
			User:         record[0],
			IP:           record[1],
			PemFilePath:  record[2],
			LastUsed:     record[3],
			Favorite:     favorite,
			DelaySeconds: delaySeconds,
			RetryCounts:  retryCounts,
		})
	}

	return history, nil
}

func (h *HistoryManager) SaveHistory(history []Connection) error {
	file, err := os.Create(h.appConfig.HistoryFile)
	if err != nil {
		return fmt.Errorf("error creating history file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, conn := range history {
		delay := strconv.FormatInt(int64(conn.DelaySeconds), 10)
		retryCounts := strconv.FormatInt(int64(conn.RetryCounts), 10)
		record := []string{
			conn.User,
			conn.IP,
			conn.PemFilePath,
			conn.LastUsed,
			fmt.Sprintf("%v", conn.Favorite),
			delay,
			retryCounts,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing to history file: %v", err)
		}
	}

	return nil
}

func (h *HistoryManager) AddConnection(user, ip, pemFilePath string) error {
	history, err := h.LoadHistory()
	if err != nil {
		return err
	}

	newConn := Connection{
		User:        user,
		IP:          ip,
		PemFilePath: pemFilePath,
		LastUsed:    time.Now().Format(time.RFC3339),
		Favorite:    false,
	}

	history = append(history, newConn)
	return h.SaveHistory(history)
}

func (h *HistoryManager) UpdateLastUsed(ip string) error {
	history, err := h.LoadHistory()
	if err != nil {
		return err
	}

	for i, conn := range history {
		if conn.IP == ip {
			history[i].LastUsed = time.Now().Format(time.RFC3339)
			break
		}
	}

	return h.SaveHistory(history)
}

func (h *HistoryManager) ToggleFavorite(ip string) error {
	history, err := h.LoadHistory()
	if err != nil {
		return err
	}

	for i, conn := range history {
		if conn.IP == ip {
			history[i].Favorite = !history[i].Favorite
			break
		}
	}

	return h.SaveHistory(history)
}
