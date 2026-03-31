package tui

// VimMode represents the current vim editing mode
type VimMode string

const (
	VimNormal  VimMode = "normal"
	VimInsert  VimMode = "insert"
	VimVisual  VimMode = "visual"
	VimCommand VimMode = "command"
)

// VimActionType identifies the editing operation to perform
type VimActionType string

const (
	VimActionNone            VimActionType = "none"
	VimActionInsert          VimActionType = "insert"
	VimActionModeChange      VimActionType = "mode_change"
	VimActionMoveLeft        VimActionType = "move_left"
	VimActionMoveRight       VimActionType = "move_right"
	VimActionMoveUp          VimActionType = "move_up"
	VimActionMoveDown        VimActionType = "move_down"
	VimActionMoveWordForward VimActionType = "move_word_forward"
	VimActionMoveWordBack    VimActionType = "move_word_back"
	VimActionMoveHome        VimActionType = "move_home"
	VimActionMoveEnd         VimActionType = "move_end"
	VimActionDeleteChar      VimActionType = "delete_char"
	VimActionDeleteLine      VimActionType = "delete_line"
	VimActionUndo            VimActionType = "undo"
	VimActionCommand         VimActionType = "command"
	VimActionNewLineBelow    VimActionType = "new_line_below"
	VimActionNewLineAbove    VimActionType = "new_line_above"
)

// VimAction describes the result of processing a keystroke in vim mode
type VimAction struct {
	Type    VimActionType
	Char    string
	Command string
}

// VimState tracks the current vim editing state
type VimState struct {
	Mode            VimMode
	Enabled         bool
	PendingOperator string // d, c, y, etc.
	Count           int
	Register        string
	CommandBuffer   string
	LastSearch      string
}

// NewVimState creates a VimState in normal mode (disabled by default)
func NewVimState() *VimState {
	return &VimState{Mode: VimNormal}
}

// Enable activates vim mode and switches to normal mode
func (v *VimState) Enable() { v.Enabled = true; v.Mode = VimNormal }

// Disable deactivates vim mode and switches to insert mode
func (v *VimState) Disable() { v.Enabled = false; v.Mode = VimInsert }

// Toggle flips vim mode on or off
func (v *VimState) Toggle() {
	if v.Enabled {
		v.Disable()
	} else {
		v.Enable()
	}
}

// HandleKey processes a keystroke and returns the action to take
func (v *VimState) HandleKey(key string) VimAction {
	if !v.Enabled {
		return VimAction{Type: VimActionInsert, Char: key}
	}

	switch v.Mode {
	case VimNormal:
		return v.handleNormal(key)
	case VimInsert:
		if key == "escape" {
			v.Mode = VimNormal
			return VimAction{Type: VimActionModeChange}
		}
		return VimAction{Type: VimActionInsert, Char: key}
	case VimCommand:
		if key == "enter" {
			cmd := v.CommandBuffer
			v.CommandBuffer = ""
			v.Mode = VimNormal
			return VimAction{Type: VimActionCommand, Command: cmd}
		}
		if key == "escape" {
			v.CommandBuffer = ""
			v.Mode = VimNormal
			return VimAction{Type: VimActionModeChange}
		}
		v.CommandBuffer += key
		return VimAction{Type: VimActionNone}
	}
	return VimAction{Type: VimActionNone}
}

func (v *VimState) handleNormal(key string) VimAction {
	switch key {
	case "i":
		v.Mode = VimInsert
		return VimAction{Type: VimActionModeChange}
	case "a":
		v.Mode = VimInsert
		return VimAction{Type: VimActionMoveRight}
	case "A":
		v.Mode = VimInsert
		return VimAction{Type: VimActionMoveEnd}
	case "I":
		v.Mode = VimInsert
		return VimAction{Type: VimActionMoveHome}
	case "h":
		return VimAction{Type: VimActionMoveLeft}
	case "l":
		return VimAction{Type: VimActionMoveRight}
	case "j":
		return VimAction{Type: VimActionMoveDown}
	case "k":
		return VimAction{Type: VimActionMoveUp}
	case "w":
		return VimAction{Type: VimActionMoveWordForward}
	case "b":
		return VimAction{Type: VimActionMoveWordBack}
	case "0":
		return VimAction{Type: VimActionMoveHome}
	case "$":
		return VimAction{Type: VimActionMoveEnd}
	case "x":
		return VimAction{Type: VimActionDeleteChar}
	case "dd":
		return VimAction{Type: VimActionDeleteLine}
	case "u":
		return VimAction{Type: VimActionUndo}
	case ":":
		v.Mode = VimCommand
		v.CommandBuffer = ""
		return VimAction{Type: VimActionModeChange}
	case "o":
		v.Mode = VimInsert
		return VimAction{Type: VimActionNewLineBelow}
	case "O":
		v.Mode = VimInsert
		return VimAction{Type: VimActionNewLineAbove}
	}
	return VimAction{Type: VimActionNone}
}
