package columns

import (
	"sort"
	"strings"
)

// WriteableColumns represents a list of columns Pop is allowed to write.
type WriteableColumns struct {
	Columns
}

// UpdateString returns the SQL column list part of the UPDATE query.
func (c WriteableColumns) UpdateString() string {
	var xs []string
	for _, t := range c.Cols {
		xs = append(xs, t.UpdateString())
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

// QuotedUpdateString returns the quoted SQL column list part of the UPDATE query.
// excludeColumns will be excluded
func (c Columns) QuotedUpdateString(quoter quoter, excludeColumns ...string) string {
	var xs []string
	outer:
	for _, t := range c.Cols {
		for _, name := range excludeColumns {
			if t.Name == name {
				continue outer
			}
		}
		xs = append(xs, t.QuotedUpdateString(quoter))
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}
