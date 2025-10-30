package type_converter

import (
	"testing"

	"github.com/mcpunzo/gomcp/types"
)

func TestMapValueToArray(t *testing.T) {
	table := []struct {
		input    map[string]*types.Tool
		expected []types.Tool
	}{
		{
			map[string]*types.Tool{
				"tool1": types.NewTool("tool1", "Tool 1", nil, nil),
				"tool2": types.NewTool("tool2", "Tool 2", nil, nil),
			},
			[]types.Tool{
				*types.NewTool("tool1", "Tool 1", nil, nil),
				*types.NewTool("tool2", "Tool 2", nil, nil),
			},
		},
		{
			nil,
			nil,
		},
	}
	for _, entry := range table {
		result := MapValueToArray(entry.input)
		if len(result) != len(entry.expected) {
			t.Errorf("Expected length %d but got %d", len(entry.expected), len(result))
		}
		for _, tool := range entry.expected {
			found := false
			for _, resTool := range result {
				if tool.Name == resTool.Name && tool.Description == resTool.Description {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected tool %v not found in result %v", tool, result)
			}
		}
	}
}
