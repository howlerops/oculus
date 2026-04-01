package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/howlerops/oculus/pkg/types"
)

// TUIStreamHandler implements query.StreamHandler by sending bubbletea messages
type TUIStreamHandler struct {
	Program *tea.Program
}

func (h *TUIStreamHandler) OnText(text string) {
	if h.Program != nil {
		h.Program.Send(StreamTextMsg{Text: text})
	}
}

func (h *TUIStreamHandler) OnToolUseStart(id, name string) {
	if h.Program != nil {
		h.Program.Send(ToolStartMsg{ToolID: id, ToolName: name})
	}
}

func (h *TUIStreamHandler) OnToolUseResult(id string, result interface{}) {
	if h.Program != nil {
		resultStr := ""
		switch v := result.(type) {
		case string:
			resultStr = v
		default:
			resultStr = ""
		}
		h.Program.Send(ToolResultMsg{ToolID: id, Result: resultStr})
	}
}

func (h *TUIStreamHandler) OnThinking(text string) {
	if h.Program != nil {
		h.Program.Send(StreamThinkingMsg{Text: text})
	}
}

func (h *TUIStreamHandler) OnComplete(stopReason types.StopReason, usage *types.Usage) {
	// Completion is handled by ResponseMsg from the query goroutine
}

func (h *TUIStreamHandler) OnError(err error) {
	if h.Program != nil {
		h.Program.Send(ErrorMsg{Err: err})
	}
}
