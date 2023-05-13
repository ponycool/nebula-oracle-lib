package oradb

import (
	"database/sql"
)

type IOra interface {
	Query(psql string, param ...any) (rows *sql.Rows, err error)
}

// Exec 执行
func (ora *Ora) Exec(psql string, param ...any) (result sql.Result, err error) {
	stmt, err := oracle.Prepare(psql)
	if err != nil {
		return nil, err
	}

	err = oracle.Ping()
	if err != nil {
		return nil, err
	}

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	result, err = stmt.Exec(param...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Query 执行查询
func (ora *Ora) Query(psql string, param ...any) (rows *sql.Rows, err error) {
	stmt, err := oracle.Prepare(psql)
	if err != nil {
		return nil, err
	}

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	rows, err = stmt.Query(param...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// QueryRow 执行查询返回一行数据
func (ora *Ora) QueryRow(psql string, param ...any) (row *sql.Row, err error) {
	stmt, err := oracle.Prepare(psql)
	if err != nil {
		return nil, err
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	row = stmt.QueryRow(param...)
	return row, nil
}
