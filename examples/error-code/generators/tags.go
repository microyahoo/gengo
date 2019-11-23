package generators

import (
	// "fmt"

	"k8s.io/gengo/types"
)

func extractTags(key string, lines []string) []string {
	// fmt.Printf("**lines = %+v\n", lines)
	values := types.ExtractCommentTags("+", lines)[key]
	return values
}
