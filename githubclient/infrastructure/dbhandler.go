package infrastructure

import (
	"backend/interfaces/database"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // mysql package

	"go.uber.org/zap"
)

// DBHandler is DB Connection.
type DBHandler struct {
	Conn *sql.DB
}

// NewDB returns new connection.
func NewDB(db, dsn string) database.DBHandler {
	conn, err := sql.Open(db, dsn+"?parseTime=true")
	if err != nil {
		zap.S().Errorw(err.Error())
	}
	handler := new(DBHandler)
	handler.Conn = conn
	return handler
}

// Query executes select statement and returns result.
func (d *DBHandler) Query(statement string, args ...interface{}) (database.Row, error) {
	row := new(SqlRow)

	stmt, err := d.Conn.Prepare(statement)
	defer stmt.Close()
	if err != nil {
		return row, err
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return row, err
	}

	row.Rows = rows
	return row, nil
}

// Execute executes insert or update or delete statement and returns result.
func (d *DBHandler) Execute(statement string, args ...interface{}) (database.Result, error) {
	result := new(SqlResult)
	stmt, err := d.Conn.Prepare(statement)
	defer stmt.Close()
	if err != nil {
		return result, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		return result, err
	}

	result.Result = res
	return result, nil
}

// SqlResult is executed result.
type SqlResult struct {
	Result sql.Result
}

// LastInsertedId returns last inserted id.
func (r SqlResult) LastInsertedId() (int64, error) {
	return r.Result.LastInsertId()
}

// RowAffected returns number of affected.
func (r SqlResult) RowAffected() (int64, error) {
	return r.Result.RowsAffected()
}

// SqlRow is rows of result.
type SqlRow struct {
	Rows *sql.Rows
}

// Scan assigned to dest.
func (r SqlRow) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest)
}

// Close closed rows.
func (r SqlRow) Close() error {
	return r.Rows.Close()
}

// Next nexts rows.
func (r SqlRow) Next() bool {
	return r.Rows.Next()
}
