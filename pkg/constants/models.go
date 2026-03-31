package constants

// Model IDs
const (
	ModelOpus4   = "claude-opus-4-20250514"
	ModelSonnet4 = "claude-sonnet-4-20250514"
	ModelHaiku4  = "claude-haiku-4-20250506"
	DefaultModel = ModelSonnet4
)

// Context windows
const (
	ContextWindow200k = 200000
	ContextWindow1m   = 1000000
)

// ModelAliases maps short names to full model IDs
var ModelAliases = map[string]string{
	"opus":   ModelOpus4,
	"sonnet": ModelSonnet4,
	"haiku":  ModelHaiku4,
	"o":      ModelOpus4,
	"s":      ModelSonnet4,
	"h":      ModelHaiku4,
}

func ResolveModelAlias(name string) string {
	if full, ok := ModelAliases[name]; ok {
		return full
	}
	return name
}
