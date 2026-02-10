package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

type WarehouseClient struct {
	db *sql.DB
}

func NewWarehouseClient() *WarehouseClient {
	dsn := os.Getenv("WAREHOUSE_DSN")
	if dsn == "" {
		dsn = "file:./dev.db?cache=shared&_pragma=busy_timeout=5000&_pragma=journal_mode=WAL"
	}

	dbConn, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	client := &WarehouseClient{db: dbConn}
	if err := client.Migrate(context.Background()); err != nil {
		panic(err)
	}
	if err := client.Seed(context.Background()); err != nil {
		panic(err)
	}

	return client
}

func (w *WarehouseClient) Close() error {
	if w.db == nil {
		return nil
	}
	return w.db.Close()
}

func (w *WarehouseClient) Migrate(ctx context.Context) error {
	queries := []string{
		`create table if not exists fact_orders (
			order_id text primary key,
			order_date date not null,
			account_id text not null,
			net_amount numeric not null
		);`,
		`create table if not exists fact_sessions (
			session_id text primary key,
			session_date date not null,
			account_id text not null,
			had_conversion integer not null
		);`,
	}

	for _, q := range queries {
		if _, err := w.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	return nil
}

func (w *WarehouseClient) Seed(ctx context.Context) error {
	var count int
	if err := w.db.QueryRowContext(ctx, "select count(*) from fact_orders").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	insertOrder := `insert into fact_orders (order_id, order_date, account_id, net_amount) values (?, ?, ?, ?);`
	insertSession := `insert into fact_sessions (session_id, session_date, account_id, had_conversion) values (?, ?, ?, ?);`

	now := time.Now().UTC()
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		accountID := "acct_001"
		orderID := fmt.Sprintf("order_%02d", i+1)
		amount := float64(1000 + i*25)

		if _, err := w.db.ExecContext(ctx, insertOrder, orderID, date, accountID, amount); err != nil {
			return err
		}

		for s := 0; s < 40; s++ {
			sessionID := fmt.Sprintf("sess_%02d_%02d", i+1, s+1)
			hadConversion := 0
			if s%10 == 0 {
				hadConversion = 1
			}
			if _, err := w.db.ExecContext(ctx, insertSession, sessionID, date, accountID, hadConversion); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *WarehouseClient) GetRevenue(ctx context.Context, startDate, endDate, accountID string) string {
	query := "select coalesce(sum(net_amount), 0) from fact_orders where order_date between ? and ?"
	args := []interface{}{startDate, endDate}
	if accountID != "" {
		query += " and account_id = ?"
		args = append(args, accountID)
	}

	var value float64
	if err := w.db.QueryRowContext(ctx, query, args...).Scan(&value); err != nil {
		return "0"
	}

	return strconv.FormatFloat(value, 'f', 2, 64)
}

func (w *WarehouseClient) GetConversionRate(ctx context.Context, startDate, endDate, accountID string) string {
	query := "select count(*) as sessions, sum(had_conversion) as conversions from fact_sessions where session_date between ? and ?"
	args := []interface{}{startDate, endDate}
	if accountID != "" {
		query += " and account_id = ?"
		args = append(args, accountID)
	}

	var sessions, conversions int
	if err := w.db.QueryRowContext(ctx, query, args...).Scan(&sessions, &conversions); err != nil {
		return "0%"
	}
	if sessions == 0 {
		return "0%"
	}

	rate := (float64(conversions) / float64(sessions)) * 100
	return fmt.Sprintf("%.2f%%", rate)
}
