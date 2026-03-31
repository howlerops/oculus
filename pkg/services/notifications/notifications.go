package notifications

import (
	"fmt"
	"os/exec"
	"runtime"
)

type NotificationType string

const (
	NotifyInfo    NotificationType = "info"
	NotifySuccess NotificationType = "success"
	NotifyWarning NotificationType = "warning"
	NotifyError   NotificationType = "error"
)

// SendSystemNotification sends an OS-level notification.
func SendSystemNotification(title, message string, notifType NotificationType) error {
	switch runtime.GOOS {
	case "darwin":
		return sendMacOSNotification(title, message)
	case "linux":
		return sendLinuxNotification(title, message)
	default:
		return fmt.Errorf("notifications not supported on %s", runtime.GOOS)
	}
}

func sendMacOSNotification(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	return exec.Command("osascript", "-e", script).Run()
}

func sendLinuxNotification(title, message string) error {
	return exec.Command("notify-send", title, message).Run()
}

// SendTerminalBell sends a terminal bell character.
func SendTerminalBell() {
	fmt.Print("\a")
}

// SendiTerm2Notification sends an iTerm2-specific notification.
func SendiTerm2Notification(message string) {
	fmt.Printf("\033]9;%s\007", message)
}
