package tool

import (
	"os"
	"strings"
	"time"
)

var logs []string

func WriteLog(msg string) {
	msg = logWithTime(msg)
	logs = append(logs, msg)
	LogWidget.SetText(strings.Join(logs, "\n"))
	LogContainer.ScrollToBottom()
}

func WriteLogFail(msg string) {
	WriteLog(msg)
	os.Exit(1)
}

func logWithTime(msg string) string {
	now := time.Now()
	logMsg := now.Format("2006-01-02 15:04:05") + " " + msg
	return logMsg
}
