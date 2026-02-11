package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"

	_ "modernc.org/sqlite"
)

type WarehouseClient struct {
	mode    string
	db      *sql.DB
	bq      *bigquery.Client
	project string
	dataset string
}

func NewWarehouseClient() *WarehouseClient {
	mode := os.Getenv("WAREHOUSE_DRIVER")
	if mode == "" {
		mode = "sqlite"
	}

	client := &WarehouseClient{mode: mode}
	if mode == "bigquery" {
		project := os.Getenv("WAREHOUSE_PROJECT")
		if project == "" {
			panic(errors.New("WAREHOUSE_PROJECT is required for bigquery"))
		}
		dataset := os.Getenv("WAREHOUSE_DATASET")
		if dataset == "" {
			panic(errors.New("WAREHOUSE_DATASET is required for bigquery"))
		}

		ctx := context.Background()
		creds := os.Getenv("WAREHOUSE_CREDENTIALS")
		var bqClient *bigquery.Client
		var err error
		if creds != "" {
			bqClient, err = bigquery.NewClient(ctx, project, option.WithCredentialsFile(creds))
		} else {
			bqClient, err = bigquery.NewClient(ctx, project)
		}
		if err != nil {
			panic(err)
		}

		client.bq = bqClient
		client.project = project
		client.dataset = dataset
		return client
	}

	dsn := os.Getenv("WAREHOUSE_DSN")
	if dsn == "" {
		dsn = "file:./dev.db?cache=shared&_pragma=busy_timeout=5000&_pragma=journal_mode=WAL"
	}

	dbConn, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	client.db = dbConn
	if err := client.Migrate(context.Background()); err != nil {
		panic(err)
	}
	if err := client.Seed(context.Background()); err != nil {
		panic(err)
	}

	return client
}

func (w *WarehouseClient) Close() error {
	if w.mode == "bigquery" {
		if w.bq == nil {
			return nil
		}
		return w.bq.Close()
	}
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
		`create table if not exists fact_active_users (
			user_id text not null,
			activity_date date not null,
			account_id text not null,
			primary key (user_id, activity_date)
		);`,
		`create table if not exists fact_subscriptions (
			subscription_id text primary key,
			account_id text not null,
			mrr numeric not null,
			is_active integer not null
		);`,
		`create table if not exists fact_mrr_snapshots (
			snapshot_date date not null,
			account_id text not null,
			mrr numeric not null,
			primary key (snapshot_date, account_id)
		);`,
		`create table if not exists fact_customer_snapshots (
			snapshot_date date not null,
			account_id text not null,
			active_customers integer not null,
			primary key (snapshot_date, account_id)
		);`,
		`create table if not exists fact_marketing_spend (
			spend_date date not null,
			account_id text not null,
			amount numeric not null,
			primary key (spend_date, account_id)
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
	insertActiveUser := `insert into fact_active_users (user_id, activity_date, account_id) values (?, ?, ?);`
	insertSubscription := `insert into fact_subscriptions (subscription_id, account_id, mrr, is_active) values (?, ?, ?, ?);`
	insertMRRSnapshot := `insert into fact_mrr_snapshots (snapshot_date, account_id, mrr) values (?, ?, ?);`
	insertCustomerSnapshot := `insert into fact_customer_snapshots (snapshot_date, account_id, active_customers) values (?, ?, ?);`
	insertMarketingSpend := `insert into fact_marketing_spend (spend_date, account_id, amount) values (?, ?, ?);`

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

		for u := 0; u < 20; u++ {
			userID := fmt.Sprintf("user_%02d", u+1)
			if _, err := w.db.ExecContext(ctx, insertActiveUser, userID, date, accountID); err != nil {
				return err
			}
		}

		spend := float64(300 + i*5)
		if _, err := w.db.ExecContext(ctx, insertMarketingSpend, date, accountID, spend); err != nil {
			return err
		}
	}

	for i := 0; i < 5; i++ {
		subscriptionID := fmt.Sprintf("sub_%02d", i+1)
		accountID := "acct_001"
		mrr := float64(500 + i*100)
		if _, err := w.db.ExecContext(ctx, insertSubscription, subscriptionID, accountID, mrr, 1); err != nil {
			return err
		}
	}

	startDate := now.AddDate(0, 0, -30).Format("2006-01-02")
	endDate := now.Format("2006-01-02")
	if _, err := w.db.ExecContext(ctx, insertMRRSnapshot, startDate, "acct_001", 2500); err != nil {
		return err
	}
	if _, err := w.db.ExecContext(ctx, insertMRRSnapshot, endDate, "acct_001", 2800); err != nil {
		return err
	}
	if _, err := w.db.ExecContext(ctx, insertCustomerSnapshot, startDate, "acct_001", 120); err != nil {
		return err
	}
	if _, err := w.db.ExecContext(ctx, insertCustomerSnapshot, endDate, "acct_001", 110); err != nil {
		return err
	}

	return nil
}

func (w *WarehouseClient) GetRevenue(ctx context.Context, startDate, endDate, accountID string) string {
	if w.mode == "bigquery" {
		query := w.bqQuery(`
			select coalesce(sum(net_amount), 0) as value
			from {{dataset}}.fact_orders
			where order_date between @start_date and @end_date
			{{account_filter}}
		`)
		params := []bigquery.QueryParameter{
			{Name: "start_date", Value: startDate},
			{Name: "end_date", Value: endDate},
		}
		if accountID != "" {
			query = w.withAccountFilter(query)
			params = append(params, bigquery.QueryParameter{Name: "account_id", Value: accountID})
		}
		value, err := w.runBigQueryFloat(ctx, query, params)
		if err != nil {
			return "0"
		}
		return strconv.FormatFloat(value, 'f', 2, 64)
	}

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
	if w.mode == "bigquery" {
		query := w.bqQuery(`
			select count(*) as sessions, sum(had_conversion) as conversions
			from {{dataset}}.fact_sessions
			where session_date between @start_date and @end_date
			{{account_filter}}
		`)
		params := []bigquery.QueryParameter{
			{Name: "start_date", Value: startDate},
			{Name: "end_date", Value: endDate},
		}
		if accountID != "" {
			query = w.withAccountFilter(query)
			params = append(params, bigquery.QueryParameter{Name: "account_id", Value: accountID})
		}
		sessions, conversions, err := w.runBigQueryCounts(ctx, query, params)
		if err != nil || sessions == 0 {
			return "0%"
		}
		rate := (float64(conversions) / float64(sessions)) * 100
		return fmt.Sprintf("%.2f%%", rate)
	}

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

func (w *WarehouseClient) GetARPU(ctx context.Context, startDate, endDate, accountID string) string {
	if w.mode == "bigquery" {
		revenueQuery := w.bqQuery(`
			select coalesce(sum(net_amount), 0) as revenue
			from {{dataset}}.fact_orders
			where order_date between @start_date and @end_date
			{{account_filter}}
		`)
		usersQuery := w.bqQuery(`
			select count(distinct user_id) as users
			from {{dataset}}.fact_active_users
			where activity_date between @start_date and @end_date
			{{account_filter}}
		`)
		params := []bigquery.QueryParameter{
			{Name: "start_date", Value: startDate},
			{Name: "end_date", Value: endDate},
		}
		if accountID != "" {
			revenueQuery = w.withAccountFilter(revenueQuery)
			usersQuery = w.withAccountFilter(usersQuery)
			params = append(params, bigquery.QueryParameter{Name: "account_id", Value: accountID})
		}
		revenue, err := w.runBigQueryFloat(ctx, revenueQuery, params)
		if err != nil {
			return "0"
		}
		users, err := w.runBigQueryInt(ctx, usersQuery, params)
		if err != nil || users == 0 {
			return "0"
		}
		arpu := revenue / float64(users)
		return strconv.FormatFloat(arpu, 'f', 2, 64)
	}

	revenueQuery := "select coalesce(sum(net_amount), 0) from fact_orders where order_date between ? and ?"
	usersQuery := "select count(distinct user_id) from fact_active_users where activity_date between ? and ?"
	args := []interface{}{startDate, endDate}
	userArgs := []interface{}{startDate, endDate}
	if accountID != "" {
		revenueQuery += " and account_id = ?"
		usersQuery += " and account_id = ?"
		args = append(args, accountID)
		userArgs = append(userArgs, accountID)
	}

	var revenue float64
	if err := w.db.QueryRowContext(ctx, revenueQuery, args...).Scan(&revenue); err != nil {
		return "0"
	}

	var users int
	if err := w.db.QueryRowContext(ctx, usersQuery, userArgs...).Scan(&users); err != nil {
		return "0"
	}
	if users == 0 {
		return "0"
	}

	arpu := revenue / float64(users)
	return strconv.FormatFloat(arpu, 'f', 2, 64)
}

func (w *WarehouseClient) GetMRR(ctx context.Context, accountID string) string {
	if w.mode == "bigquery" {
		query := w.bqQuery(`
			select coalesce(sum(mrr), 0) as value
			from {{dataset}}.fact_subscriptions
			where is_active = 1
			{{account_filter}}
		`)
		params := []bigquery.QueryParameter{}
		if accountID != "" {
			query = w.withAccountFilter(query)
			params = append(params, bigquery.QueryParameter{Name: "account_id", Value: accountID})
		}
		value, err := w.runBigQueryFloat(ctx, query, params)
		if err != nil {
			return "0"
		}
		return strconv.FormatFloat(value, 'f', 2, 64)
	}

	query := "select coalesce(sum(mrr), 0) from fact_subscriptions where is_active = 1"
	args := []interface{}{}
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

func (w *WarehouseClient) GetNRR(ctx context.Context, startDate, endDate, accountID string) string {
	if w.mode == "bigquery" {
		startQuery := w.bqQuery(`
			select mrr from {{dataset}}.fact_mrr_snapshots
			where snapshot_date = @start_date
			{{account_filter}}
		`)
		endQuery := w.bqQuery(`
			select mrr from {{dataset}}.fact_mrr_snapshots
			where snapshot_date = @end_date
			{{account_filter}}
		`)
		params := []bigquery.QueryParameter{
			{Name: "start_date", Value: startDate},
			{Name: "end_date", Value: endDate},
		}
		if accountID != "" {
			startQuery = w.withAccountFilter(startQuery)
			endQuery = w.withAccountFilter(endQuery)
			params = append(params, bigquery.QueryParameter{Name: "account_id", Value: accountID})
		}
		startMRR, err := w.runBigQueryFloat(ctx, startQuery, params)
		if err != nil || startMRR == 0 {
			return "0%"
		}
		endMRR, err := w.runBigQueryFloat(ctx, endQuery, params)
		if err != nil {
			return "0%"
		}
		nrr := (endMRR / startMRR) * 100
		return fmt.Sprintf("%.2f%%", nrr)
	}

	query := "select mrr from fact_mrr_snapshots where snapshot_date = ?"
	args := []interface{}{startDate}
	if accountID != "" {
		query += " and account_id = ?"
		args = append(args, accountID)
	}

	var startMRR float64
	if err := w.db.QueryRowContext(ctx, query, args...).Scan(&startMRR); err != nil {
		return "0%"
	}

	endQuery := "select mrr from fact_mrr_snapshots where snapshot_date = ?"
	endArgs := []interface{}{endDate}
	if accountID != "" {
		endQuery += " and account_id = ?"
		endArgs = append(endArgs, accountID)
	}

	var endMRR float64
	if err := w.db.QueryRowContext(ctx, endQuery, endArgs...).Scan(&endMRR); err != nil {
		return "0%"
	}
	if startMRR == 0 {
		return "0%"
	}

	nrr := (endMRR / startMRR) * 100
	return fmt.Sprintf("%.2f%%", nrr)
}

func (w *WarehouseClient) GetChurnRate(ctx context.Context, startDate, endDate, accountID string) string {
	if w.mode == "bigquery" {
		startQuery := w.bqQuery(`
			select active_customers from {{dataset}}.fact_customer_snapshots
			where snapshot_date = @start_date
			{{account_filter}}
		`)
		endQuery := w.bqQuery(`
			select active_customers from {{dataset}}.fact_customer_snapshots
			where snapshot_date = @end_date
			{{account_filter}}
		`)
		params := []bigquery.QueryParameter{
			{Name: "start_date", Value: startDate},
			{Name: "end_date", Value: endDate},
		}
		if accountID != "" {
			startQuery = w.withAccountFilter(startQuery)
			endQuery = w.withAccountFilter(endQuery)
			params = append(params, bigquery.QueryParameter{Name: "account_id", Value: accountID})
		}
		startCustomers, err := w.runBigQueryInt(ctx, startQuery, params)
		if err != nil || startCustomers == 0 {
			return "0%"
		}
		endCustomers, err := w.runBigQueryInt(ctx, endQuery, params)
		if err != nil {
			return "0%"
		}
		lost := startCustomers - endCustomers
		if lost < 0 {
			lost = 0
		}
		churn := (float64(lost) / float64(startCustomers)) * 100
		return fmt.Sprintf("%.2f%%", churn)
	}

	query := "select active_customers from fact_customer_snapshots where snapshot_date = ?"
	args := []interface{}{startDate}
	if accountID != "" {
		query += " and account_id = ?"
		args = append(args, accountID)
	}

	var startCustomers int
	if err := w.db.QueryRowContext(ctx, query, args...).Scan(&startCustomers); err != nil {
		return "0%"
	}

	endQuery := "select active_customers from fact_customer_snapshots where snapshot_date = ?"
	endArgs := []interface{}{endDate}
	if accountID != "" {
		endQuery += " and account_id = ?"
		endArgs = append(endArgs, accountID)
	}

	var endCustomers int
	if err := w.db.QueryRowContext(ctx, endQuery, endArgs...).Scan(&endCustomers); err != nil {
		return "0%"
	}
	if startCustomers == 0 {
		return "0%"
	}

	lost := startCustomers - endCustomers
	if lost < 0 {
		lost = 0
	}
	churn := (float64(lost) / float64(startCustomers)) * 100
	return fmt.Sprintf("%.2f%%", churn)
}

func (w *WarehouseClient) GetLTV(ctx context.Context, startDate, endDate, accountID string) string {
	if w.mode == "bigquery" {
		arpuValue := w.GetARPU(ctx, startDate, endDate, accountID)
		arpu, err := strconv.ParseFloat(arpuValue, 64)
		if err != nil {
			return "0"
		}
		churnValue := w.GetChurnRate(ctx, startDate, endDate, accountID)
		churn, err := strconv.ParseFloat(strings.TrimSuffix(churnValue, "%"), 64)
		if err != nil || churn <= 0 {
			return "0"
		}
		ltv := arpu / (churn / 100)
		return strconv.FormatFloat(ltv, 'f', 2, 64)
	}

	arpuValue := w.GetARPU(ctx, startDate, endDate, accountID)
	arpu, err := strconv.ParseFloat(arpuValue, 64)
	if err != nil {
		return "0"
	}

	query := "select active_customers from fact_customer_snapshots where snapshot_date = ?"
	args := []interface{}{startDate}
	if accountID != "" {
		query += " and account_id = ?"
		args = append(args, accountID)
	}

	var startCustomers int
	if err := w.db.QueryRowContext(ctx, query, args...).Scan(&startCustomers); err != nil {
		return "0"
	}

	endQuery := "select active_customers from fact_customer_snapshots where snapshot_date = ?"
	endArgs := []interface{}{endDate}
	if accountID != "" {
		endQuery += " and account_id = ?"
		endArgs = append(endArgs, accountID)
	}

	var endCustomers int
	if err := w.db.QueryRowContext(ctx, endQuery, endArgs...).Scan(&endCustomers); err != nil {
		return "0"
	}
	if startCustomers == 0 {
		return "0"
	}

	lost := startCustomers - endCustomers
	if lost < 0 {
		lost = 0
	}
	churnRate := float64(lost) / float64(startCustomers)
	if churnRate <= 0 {
		return "0"
	}

	ltv := arpu / churnRate
	return strconv.FormatFloat(ltv, 'f', 2, 64)
}

func (w *WarehouseClient) GetCAC(ctx context.Context, startDate, endDate, accountID string) string {
	if w.mode == "bigquery" {
		spendQuery := w.bqQuery(`
			select coalesce(sum(amount), 0) as spend
			from {{dataset}}.fact_marketing_spend
			where spend_date between @start_date and @end_date
			{{account_filter}}
		`)
		newCustomersQuery := w.bqQuery(`
			select count(distinct account_id) as customers
			from {{dataset}}.fact_orders
			where order_date between @start_date and @end_date
			{{account_filter}}
		`)
		params := []bigquery.QueryParameter{
			{Name: "start_date", Value: startDate},
			{Name: "end_date", Value: endDate},
		}
		if accountID != "" {
			spendQuery = w.withAccountFilter(spendQuery)
			newCustomersQuery = w.withAccountFilter(newCustomersQuery)
			params = append(params, bigquery.QueryParameter{Name: "account_id", Value: accountID})
		}
		spend, err := w.runBigQueryFloat(ctx, spendQuery, params)
		if err != nil {
			return "0"
		}
		newCustomers, err := w.runBigQueryInt(ctx, newCustomersQuery, params)
		if err != nil || newCustomers == 0 {
			return "0"
		}
		cac := spend / float64(newCustomers)
		return strconv.FormatFloat(cac, 'f', 2, 64)
	}

	spendQuery := "select coalesce(sum(amount), 0) from fact_marketing_spend where spend_date between ? and ?"
	args := []interface{}{startDate, endDate}
	if accountID != "" {
		spendQuery += " and account_id = ?"
		args = append(args, accountID)
	}

	var spend float64
	if err := w.db.QueryRowContext(ctx, spendQuery, args...).Scan(&spend); err != nil {
		return "0"
	}

	newCustomersQuery := "select count(distinct account_id) from fact_orders where order_date between ? and ?"
	newArgs := []interface{}{startDate, endDate}
	if accountID != "" {
		newCustomersQuery += " and account_id = ?"
		newArgs = append(newArgs, accountID)
	}

	var newCustomers int
	if err := w.db.QueryRowContext(ctx, newCustomersQuery, newArgs...).Scan(&newCustomers); err != nil {
		return "0"
	}
	if newCustomers == 0 {
		return "0"
	}

	cac := spend / float64(newCustomers)
	return strconv.FormatFloat(cac, 'f', 2, 64)
}

func (w *WarehouseClient) bqQuery(sqlText string) string {
	query := strings.ReplaceAll(sqlText, "{{dataset}}", fmt.Sprintf("`%s.%s`", w.project, w.dataset))
	query = strings.ReplaceAll(query, "{{account_filter}}", "")
	return query
}

func (w *WarehouseClient) withAccountFilter(sqlText string) string {
	return strings.ReplaceAll(sqlText, "{{account_filter}}", "and account_id = @account_id")
}

func (w *WarehouseClient) runBigQueryFloat(ctx context.Context, sqlText string, params []bigquery.QueryParameter) (float64, error) {
	query := w.bq.Query(sqlText)
	query.Parameters = params
	iter, err := query.Read(ctx)
	if err != nil {
		return 0, err
	}
	var row struct {
		Value float64 `bigquery:"value"`
		Spend float64 `bigquery:"spend"`
		Revenue float64 `bigquery:"revenue"`
		MRR float64 `bigquery:"mrr"`
	}
	if err := iter.Next(&row); err != nil {
		return 0, err
	}
	if row.Value != 0 {
		return row.Value, nil
	}
	if row.Spend != 0 {
		return row.Spend, nil
	}
	if row.Revenue != 0 {
		return row.Revenue, nil
	}
	if row.MRR != 0 {
		return row.MRR, nil
	}
	return 0, nil
}

func (w *WarehouseClient) runBigQueryInt(ctx context.Context, sqlText string, params []bigquery.QueryParameter) (int, error) {
	query := w.bq.Query(sqlText)
	query.Parameters = params
	iter, err := query.Read(ctx)
	if err != nil {
		return 0, err
	}
	var row struct {
		Users int64 `bigquery:"users"`
		Customers int64 `bigquery:"customers"`
		ActiveCustomers int64 `bigquery:"active_customers"`
	}
	if err := iter.Next(&row); err != nil {
		return 0, err
	}
	if row.Users != 0 {
		return int(row.Users), nil
	}
	if row.Customers != 0 {
		return int(row.Customers), nil
	}
	if row.ActiveCustomers != 0 {
		return int(row.ActiveCustomers), nil
	}
	return 0, nil
}

func (w *WarehouseClient) runBigQueryCounts(ctx context.Context, sqlText string, params []bigquery.QueryParameter) (int, int, error) {
	query := w.bq.Query(sqlText)
	query.Parameters = params
	iter, err := query.Read(ctx)
	if err != nil {
		return 0, 0, err
	}
	var row struct {
		Sessions int64 `bigquery:"sessions"`
		Conversions int64 `bigquery:"conversions"`
	}
	if err := iter.Next(&row); err != nil {
		return 0, 0, err
	}
	return int(row.Sessions), int(row.Conversions), nil
}
