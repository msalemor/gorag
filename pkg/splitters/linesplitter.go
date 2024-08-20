package splitters

import "strings"

// IParagraphSplitter is an interface for splitting text into paragraphs.
type LineSplitter struct{}

// Split splits text into paragraphs.
func (s LineSplitter) Split(opts SplitterOpts) []string {
	return strings.Split(opts.Content, "\n")
}
