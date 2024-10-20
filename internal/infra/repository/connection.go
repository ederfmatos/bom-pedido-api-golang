package repository

import (
	"bom-pedido-api/internal/infra/telemetry"
	"bom-pedido-api/internal/infra/tenant"
	"context"
	"database/sql"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"regexp"
	"runtime"
	"strings"
)

type SqlConnection interface {
	Sql(sql string) ConnectionBuilder
	InTransaction(ctx context.Context, handler func(transaction SqlTransaction, ctx context.Context) error) error
}

type SqlTransaction interface {
	Sql(sql string) ConnectionBuilder
}

type RowMapper func(getValues func(dest ...any) error) error

type ConnectionBuilder interface {
	Values(value ...interface{}) ConnectionBuilder
	Update(ctx context.Context) error
	FindOne(ctx context.Context, values ...interface{}) (bool, error)
	Count(ctx context.Context, values ...interface{}) error
	List(ctx context.Context, mapper RowMapper) error
	Exists(ctx context.Context) (bool, error)
}

type Connection interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type DefaultSqlConnection struct {
	database *sql.DB
}

func NewDefaultSqlConnection(database *sql.DB) SqlConnection {
	return &DefaultSqlConnection{database}
}

type DefaultSqlTransaction struct {
	transaction *sql.Tx
}

func (transaction *DefaultSqlTransaction) Sql(sql string) ConnectionBuilder {
	return &DefaultConnectionBuilder{
		sql:        &sql,
		values:     []interface{}{},
		connection: transaction.transaction,
	}
}

func (connection *DefaultSqlConnection) InTransaction(ctx context.Context, handler func(connection SqlTransaction, ctx context.Context) error) error {
	ctx, span := telemetry.StartSpan(ctx, "SqlConnection.InTransaction")
	defer span.End()
	span.AddEvent("Begin transaction")
	tx, err := connection.database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	transaction := &DefaultSqlTransaction{transaction: tx}
	err = handler(transaction, ctx)
	if err == nil {
		span.AddEvent("Commiting transaction")
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	span.AddEvent("Rollback transaction")
	rollbackError := tx.Rollback()
	if rollbackError == nil {
		return err
	}
	return rollbackError
}

func (connection *DefaultSqlConnection) Sql(sql string) ConnectionBuilder {
	return &DefaultConnectionBuilder{
		sql:        &sql,
		values:     []interface{}{},
		connection: connection.database,
	}
}

type DefaultConnectionBuilder struct {
	sql        *string
	values     []interface{}
	connection Connection
}

func (builder *DefaultConnectionBuilder) Values(value ...interface{}) ConnectionBuilder {
	builder.values = value
	return builder
}

func (builder *DefaultConnectionBuilder) Update(ctx context.Context) error {
	span := builder.createSpan(ctx)
	defer span.End()
	statement, err := builder.connection.PrepareContext(ctx, *builder.sql)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.ExecContext(ctx, builder.values...)
	return err
}

func (builder *DefaultConnectionBuilder) Count(ctx context.Context, values ...interface{}) error {
	_, err := builder.FindOne(ctx, values...)
	return err
}

func (builder *DefaultConnectionBuilder) FindOne(ctx context.Context, values ...interface{}) (bool, error) {
	span := builder.createSpan(ctx)
	defer span.End()
	statement, err := builder.connection.PrepareContext(ctx, *builder.sql)
	if err != nil {
		return false, err
	}
	defer statement.Close()
	err = statement.QueryRowContext(ctx, builder.values...).Scan(values...)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (builder *DefaultConnectionBuilder) Exists(ctx context.Context) (bool, error) {
	span := builder.createSpan(ctx)
	defer span.End()
	statement, err := builder.connection.PrepareContext(ctx, *builder.sql)
	if err != nil {
		return false, err
	}
	defer statement.Close()
	rows, err := statement.QueryContext(ctx, builder.values...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (builder *DefaultConnectionBuilder) List(ctx context.Context, mapper RowMapper) error {
	span := builder.createSpan(ctx)
	defer span.End()
	statement, err := builder.connection.PrepareContext(ctx, *builder.sql)
	if err != nil {
		return err
	}
	defer statement.Close()
	rows, err := statement.QueryContext(ctx, builder.values...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = mapper(rows.Scan)
		if err != nil {
			return err
		}
	}
	return nil
}

var spanRegex = regexp.MustCompile(`\(\*|repository\.|\)`)

func (builder *DefaultConnectionBuilder) createSpan(ctx context.Context) trace.Span {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	function := strings.Split(frame.Function, "/")
	functionName := spanRegex.ReplaceAllString(function[len(function)-1], "")
	var tenantId string
	if tenantValue := ctx.Value(tenant.Id); tenantValue != nil {
		tenantId = tenantValue.(string)
	}
	_, span := telemetry.StartSpan(ctx, functionName, "sql", *builder.sql, "tenant.id", tenantId)
	return span
}
