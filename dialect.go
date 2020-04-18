package pop

import (
	"io"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/pop/columns"
	"github.com/pkg/errors"
)

var ErrNotImplemented = errors.New("sql: not implemented")

type crudable interface {
	SelectOne(store, *Model, Query) error
	SelectMany(store, *Model, Query) error
	Create(store, *Model, columns.Columns) error
	Update(store, *Model, columns.Columns) error
	Upsert(store, *Model, columns.Columns, string) error
	Destroy(store, *Model) error
}

type fizzable interface {
	FizzTranslator() fizz.Translator
}

type quotable interface {
	Quote(key string) string
}

type dialect interface {
	crudable
	fizzable
	quotable
	Name() string
	URL() string
	MigrationURL() string
	Details() *ConnectionDetails
	TranslateSQL(string) string
	CreateDB() error
	DropDB() error
	DumpSchema(io.Writer) error
	LoadSchema(io.Reader) error
	Lock(func() error) error
	TruncateAll(*Connection) error
}

type afterOpenable interface {
	AfterOpen(*Connection) error
}
