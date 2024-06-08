package repository

import (
	"context"
	"database/sql"
)

type SqlConnection interface {
	Sql(sql string) ConnectionBuilder
}

type ConnectionBuilder interface {
	Values(value ...interface{}) ConnectionBuilder
	Update(ctx context.Context) error
	FindOne(ctx context.Context, values ...interface{}) (bool, error)
	Exists(ctx context.Context) (bool, error)
}

type DefaultSqlConnection struct {
	database *sql.DB
}

func NewDefaultSqlConnection(database *sql.DB) SqlConnection {
	return &DefaultSqlConnection{database: database}
}

func (connection *DefaultSqlConnection) Sql(sql string) ConnectionBuilder {
	return &DefaultConnectionBuilder{
		sql:      &sql,
		database: connection.database,
		values:   []interface{}{},
	}
}

type DefaultConnectionBuilder struct {
	sql      *string
	database *sql.DB
	values   []interface{}
}

func (builder *DefaultConnectionBuilder) Values(value ...interface{}) ConnectionBuilder {
	builder.values = value
	return builder
}

func (builder *DefaultConnectionBuilder) Update(ctx context.Context) error {
	statement, err := builder.database.PrepareContext(ctx, *builder.sql)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.ExecContext(ctx, builder.values...)
	return err
}

func (builder *DefaultConnectionBuilder) FindOne(ctx context.Context, values ...interface{}) (bool, error) {
	statement, err := builder.database.PrepareContext(ctx, *builder.sql)
	if err != nil {
		return false, err
	}
	defer statement.Close()
	err = statement.QueryRowContext(ctx, builder.values...).Scan(values...)
	return true, err
}

func (builder *DefaultConnectionBuilder) Exists(ctx context.Context) (bool, error) {
	statement, err := builder.database.PrepareContext(ctx, *builder.sql)
	if err != nil {
		return false, err
	}
	defer statement.Close()
	rows, err := statement.QueryContext(ctx, builder.values...)
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}
