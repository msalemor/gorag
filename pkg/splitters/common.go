package splitters

type SplitterOpts struct {
	Content       string
	LineSize      int
	ParagraphSize int
	Overflows     int
}

type ISplitter interface {
	Split(opts SplitterOpts) []string
}
