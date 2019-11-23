package generators

import (
	"k8s.io/gengo/types"
)

func extractTags(key string, lines []string) []string {
	values := types.ExtractCommentTags("+", lines)[key]
	return values
}
