package infrastructure

import (
	"backend/interfaces/database"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"go.uber.org/zap"
)

type DBHandler struct {
	Conn *sql.DB
}

func NewDB(db, dsn string) database.DBHandler {
	conn, err := sql.Open(db, dsn+"?parseTime=true")
	if err != nil {
		zap.S().Errorw(err.Error())
	}
	handler := new(DBHandler)
	handler.Conn = conn
	return handler
}

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

type SqlResult struct {
	Result sql.Result
}

func (r SqlResult) LastInsertedId() (int64, error) {
	return r.Result.LastInsertId()
}

func (r SqlResult) RowAffected() (int64, error) {
	return r.Result.RowsAffected()
}

type SqlRow struct {
	Rows *sql.Rows
}

func (r SqlRow) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest)
}

func (r SqlRow) Close() error {
	return r.Rows.Close()
}

func (r SqlRow) Next() bool {
	return r.Rows.Next()
}
