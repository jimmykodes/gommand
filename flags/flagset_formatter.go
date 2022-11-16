package flags

import (
	"sort"
	"strings"
)

type FlagSetFormatter struct {
	fs *FlagSet
}

func NewFlagSetFormatter(fs *FlagSet) *FlagSetFormatter {
	return &FlagSetFormatter{fs: fs}
}

func (fsf *FlagSetFormatter) Empty() bool {
	return len(fsf.fs.flags) == 0
}

func (fsf *FlagSetFormatter) Format() string {
	names := make([]string, 0, len(fsf.fs.flags))
	maxLen := 0
	for n, flag := range fsf.fs.flags {
		names = append(names, n)
		if l := len(flag.Name()); l > maxLen {
			maxLen = l
		}
	}
	sort.Strings(names)
	strs := make([]string, len(names))
	for i, name := range names {
		strs[i] = Stringer(fsf.fs.flags[name], maxLen)
	}
	return strings.Join(strs, "\n")
}
