package lens

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Focus.Type != LensFocus {
		t.Error("Focus type wrong")
	}
	if cfg.Scan.Type != LensScan {
		t.Error("Scan type wrong")
	}
	if cfg.Craft.Type != LensCraft {
		t.Error("Craft type wrong")
	}
	if !cfg.Focus.Enabled {
		t.Error("Focus should be enabled")
	}
}

func TestStaticRouter(t *testing.T) {
	router := NewStaticRouter()

	tests := []struct {
		tool string
		want LensType
	}{
		{"Read", LensScan},
		{"Glob", LensScan},
		{"Grep", LensScan},
		{"Bash", LensCraft},
		{"Edit", LensCraft},
		{"Write", LensCraft},
		{"Agent", LensFocus},
		{"AskUserQuestion", LensFocus},
		{"UnknownTool", LensFocus}, // default
	}
	for _, tt := range tests {
		got := router.RouteToolCall(tt.tool)
		if got != tt.want {
			t.Errorf("RouteToolCall(%q) = %s, want %s", tt.tool, got, tt.want)
		}
	}
}

func TestRouteMessage(t *testing.T) {
	router := NewStaticRouter()

	if router.RouteMessage("find all go files") != LensScan {
		t.Error("'find' should route to Scan")
	}
	if router.RouteMessage("create a new function") != LensCraft {
		t.Error("'create' should route to Craft")
	}
	if router.RouteMessage("what should we do about this") != LensFocus {
		t.Error("ambiguous should route to Focus")
	}
}
