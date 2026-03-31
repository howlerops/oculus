package constants

// XML tags used in system prompts
const (
	CommandMessageTag     = "command-message"
	CommandNameTag        = "command-name"
	CommandArgsTag        = "command-args"
	LocalCommandStdoutTag = "local-command-stdout"
	LocalCommandCaveatTag = "local-command-caveat"
	SystemReminderTag     = "system-reminder"
)

// Message constants
const (
	NoContentMessage = "(no content)"
)

// API retry delays
const (
	RetryBaseDelayMs = 1000
	RetryMaxDelayMs  = 60000
)

// File paths
const (
	OculusConfigDir = ".oculus"
	SettingsFile    = "settings.json"
	ConfigFile      = "config.json"
	OculusMdFile    = "OCULUS.md"
)
