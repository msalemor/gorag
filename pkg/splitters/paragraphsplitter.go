package splitters

import "strings"

// IParagraphSplitter is an interface for splitting text into paragraphs.
type ParagraphSplitter struct{}

// Split splits text into paragraphs.
func (s ParagraphSplitter) Split(opts SplitterOpts) []string {
	return strings.Split(opts.Content, "\n\n")
}
