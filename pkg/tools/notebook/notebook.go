package notebook

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

// NotebookCell represents a single cell in a Jupyter notebook.
type NotebookCell struct {
	CellType       string        `json:"cell_type"`
	ID             string        `json:"id,omitempty"`
	Source         interface{}   `json:"source"` // string or []string
	Metadata       interface{}   `json:"metadata"`
	Outputs        []interface{} `json:"outputs,omitempty"`
	ExecutionCount *int          `json:"execution_count,omitempty"`
}

// Notebook is the top-level structure of a .ipynb file.
type Notebook struct {
	Cells         []NotebookCell         `json:"cells"`
	Metadata      map[string]interface{} `json:"metadata"`
	NBFormat      int                    `json:"nbformat"`
	NBFormatMinor int                    `json:"nbformat_minor"`
}

// NotebookEditTool edits a Jupyter notebook cell.
type NotebookEditTool struct {
	tool.BaseTool
}

// NewNotebookEditTool creates a new NotebookEditTool.
func NewNotebookEditTool() *NotebookEditTool {
	return &NotebookEditTool{
		BaseTool: tool.BaseTool{
			ToolName:       "NotebookEdit",
			ToolSearchHint: "jupyter notebook ipynb cell edit",
		},
	}
}

func (t *NotebookEditTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"notebook_path": map[string]interface{}{
				"type":        "string",
				"description": "Path to the .ipynb notebook file",
			},
			"cell_number": map[string]interface{}{
				"type":        "number",
				"description": "0-based cell index (used when cell_id is not provided)",
			},
			"cell_id": map[string]interface{}{
				"type":        "string",
				"description": "Cell ID or cell-N index string (e.g. cell-0)",
			},
			"new_source": map[string]interface{}{
				"type":        "string",
				"description": "New source content for the cell",
			},
			"cell_type": map[string]interface{}{
				"type":        "string",
				"description": "Cell type: code or markdown",
			},
			"edit_mode": map[string]interface{}{
				"type":        "string",
				"description": "Edit mode: replace (default), insert, or delete",
			},
		},
		Required: []string{"notebook_path"},
	}
}

func (t *NotebookEditTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Edit cells in a Jupyter notebook (.ipynb). Supports replace, insert, and delete modes.", nil
}

func (t *NotebookEditTool) Prompt(_ context.Context) (string, error) {
	return "Edit Jupyter notebook (.ipynb) cells.\n- notebook_path must be absolute\n- cell_number is 0-indexed\n- edit_mode: replace (default), insert, delete", nil
}

func (t *NotebookEditTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	path, _ := input["notebook_path"].(string)
	if path == "" {
		return &tool.Result{Data: "Error: notebook_path is required"}, nil
	}
	if !strings.HasSuffix(path, ".ipynb") {
		return &tool.Result{Data: "Error: file must have .ipynb extension"}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error reading notebook: %v", err)}, nil
	}

	var nb Notebook
	if err := json.Unmarshal(data, &nb); err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error parsing notebook: %v", err)}, nil
	}

	editMode, _ := input["edit_mode"].(string)
	if editMode == "" {
		editMode = "replace"
	}
	newSource, _ := input["new_source"].(string)
	cellType, _ := input["cell_type"].(string)

	// Resolve cell index from cell_id or cell_number.
	idx := -1
	if cellID, ok := input["cell_id"].(string); ok && cellID != "" {
		if strings.HasPrefix(cellID, "cell-") {
			fmt.Sscanf(cellID, "cell-%d", &idx)
		} else {
			for i, c := range nb.Cells {
				if c.ID == cellID {
					idx = i
					break
				}
			}
		}
	} else if n, ok := input["cell_number"].(float64); ok {
		idx = int(n)
	}

	switch editMode {
	case "delete":
		if idx < 0 || idx >= len(nb.Cells) {
			return &tool.Result{Data: fmt.Sprintf("Error: cell index %d out of range (notebook has %d cells)", idx, len(nb.Cells))}, nil
		}
		nb.Cells = append(nb.Cells[:idx], nb.Cells[idx+1:]...)

	case "insert":
		if cellType == "" {
			cellType = "code"
		}
		newCell := NotebookCell{
			CellType: cellType,
			Source:   newSource,
			Metadata: map[string]interface{}{},
			Outputs:  []interface{}{},
		}
		insertIdx := idx + 1
		if insertIdx < 0 {
			insertIdx = 0
		}
		if insertIdx > len(nb.Cells) {
			insertIdx = len(nb.Cells)
		}
		nb.Cells = append(nb.Cells[:insertIdx], append([]NotebookCell{newCell}, nb.Cells[insertIdx:]...)...)

	case "replace":
		if idx < 0 || idx >= len(nb.Cells) {
			return &tool.Result{Data: fmt.Sprintf("Error: cell index %d out of range (notebook has %d cells)", idx, len(nb.Cells))}, nil
		}
		nb.Cells[idx].Source = newSource
		if cellType != "" {
			nb.Cells[idx].CellType = cellType
		}
		nb.Cells[idx].Outputs = []interface{}{}
		nb.Cells[idx].ExecutionCount = nil

	default:
		return &tool.Result{Data: fmt.Sprintf("Error: unknown edit_mode %q (use replace, insert, or delete)", editMode)}, nil
	}

	out, err := json.MarshalIndent(nb, "", " ")
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error serializing notebook: %v", err)}, nil
	}
	if err := os.WriteFile(path, out, 0644); err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error writing notebook: %v", err)}, nil
	}

	return &tool.Result{Data: fmt.Sprintf("Notebook %s updated: %s cell at index %d", path, editMode, idx)}, nil
}
