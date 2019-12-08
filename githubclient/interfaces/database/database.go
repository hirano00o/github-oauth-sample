package database

// DBHandler is ...
type DBHandler interface {
	Query(string, ...interface{}) (Row, error)
	Execute(string, ...interface{}) (Result, error)
}

// Result is ...
type Result interface {
	RowAffected() (int64, error)
	LastInsertedId() (int64, error)
}

// Row is ...
type Row interface {
	Scan(...interface{}) error
	Close() error
	Next() bool
}
