package splitters

import "strings"

// IParagraphSplitter is an interface for splitting text into paragraphs.
type WordCoundSplitter struct{}

// Split splits text into paragraphs.
func (s WordCoundSplitter) Split(opts SplitterOpts) []string {

	return strings.Split(opts.Content, " ")
}
