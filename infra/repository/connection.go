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
	FindOne(ctx context.Context, values ...interface{}) error
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
	statement, err := builder.database.Prepare(*builder.sql)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.ExecContext(ctx, builder.values)
	return err
}

func (builder *DefaultConnectionBuilder) FindOne(ctx context.Context, values ...interface{}) error {
	statement, err := builder.database.Prepare(*builder.sql)
	if err != nil {
		return err
	}
	defer statement.Close()
	err = statement.QueryRowContext(ctx, builder.values).Scan(values)
	return err
}
