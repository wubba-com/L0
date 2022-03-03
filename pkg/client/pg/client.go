package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/wubba-com/L0/internal/config"
	"github.com/wubba-com/L0/pkg/utils"
	"log"
	"time"
)

// NewClient - пул бд
func NewClient(ctx context.Context, config *config.Config, maxAttempts int) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			log.Printf("failed to connect postgresql: %s", err.Error())
			return err
		}
		return nil
	}, maxAttempts, 3*time.Second)
	if err != nil {
		log.Fatalf("error do with tries postgresql: %s", err.Error())
	}
	//defer pool.Close()
	return pool, nil
}

// Client - для работы с postgres
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
}
