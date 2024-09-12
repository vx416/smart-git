package git

import (
	"fmt"
	"strings"

	"github.com/sourcegraph/go-diff/diff"
)

// ParseDiffContentToMultiFileDiff parses the given diff content and returns a slice of FileDiff.
func ParseDiffContentToMultiFileDiff(diffContent []byte) ([]*diff.FileDiff, error) {
	return diff.ParseMultiFileDiff(diffContent)
}

// FileDiffToString converts a FileDiff to a string.
func FileDiffToString(diffFile *diff.FileDiff) string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("--- %s\n", diffFile.OrigName))
	str.WriteString(fmt.Sprintf("+++ %s\n", diffFile.NewName))
	for _, hunk := range diffFile.Hunks {
		str.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@ %s\n", hunk.OrigStartLine, hunk.OrigLines, hunk.NewStartLine, hunk.NewLines, hunk.Section))
		str.Write(hunk.Body)
	}
	return str.String()
}

// MultiFileDiffToPayloadSlice converts a slice of FileDiff to a slice of string payloads.
func MultiFileDiffToPayloadSlice(diffFiles []*diff.FileDiff, maxPayloadSize int) []string {
	payloads := make([]string, 0, len(diffFiles))
	for _, diffFile := range diffFiles {
		strBuilder := strings.Builder{}
		for _, hunk := range diffFile.Hunks {
			if strBuilder.Len() == 0 {
				strBuilder.WriteString(fmt.Sprintf("--- %s\n", diffFile.OrigName))
				strBuilder.WriteString(fmt.Sprintf("+++ %s\n", diffFile.NewName))
			}
			strBuilder.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@ %s\n", hunk.OrigStartLine, hunk.OrigLines, hunk.NewStartLine, hunk.NewLines, hunk.Section))
			strBuilder.Write(hunk.Body)
			if maxPayloadSize > 0 && strBuilder.Len() > maxPayloadSize {
				payloads = append(payloads, strBuilder.String())
				strBuilder.Reset()
			}
		}
		if strBuilder.Len() > 0 {
			payloads = append(payloads, strBuilder.String())
		}
	}
	return payloads
}
