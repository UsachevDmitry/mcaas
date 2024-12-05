package internal

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)
func (p *PostgresStorage) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return p.db.Exec(ctx, query, args...)
}

func (p *PostgresStorage) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return p.db.Query(ctx, query, args...)
}

func (p *PostgresStorage) Connect() error {
	var err error
	var poolConfig *pgxpool.Config

	poolConfig, err = pgxpool.ParseConfig(*DatabaseDsn)
	if err != nil {
		GlobalSugar.Fatalf("Не удалось инициализировать пул: %v", err)
	}
	poolConfig.MaxConns = 1 // Максимальное количество соединений в пуле
	poolConfig.MinConns = 1 // Минимальное количество поддерживаемых соединений
	poolConfig.ConnConfig.TLSConfig = nil
	//poolConfig.ConnConfig.ConnectTimeout = 500 * time.Millisecond
	p.db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		GlobalSugar.Fatalf("QueryRow failed: %v\n", err)
		return err
	}
	
	return err
} 

func (p *PostgresStorage) Close() {
	p.db.Close()
}

func (p *PostgresStorage) CreateTableGauge(ctx context.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Second)
	defer cancel()
	_, err := p.Exec(ctxWithTimeout, `CREATE TABLE IF NOT EXISTS metrics_gauge (
		"key" TEXT,
		"value" DOUBLE PRECISION
	)`)
	if err != nil {
		GlobalSugar.Fatal(err)
	}
}
func (p *PostgresStorage) CreateTableCounter(ctx context.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Second)
	defer cancel()
	_, err := p.Exec(ctxWithTimeout,`CREATE TABLE IF NOT EXISTS metrics_counter (
		"key" TEXT,
		"value" BIGINT
	)`)
	if err != nil {
		GlobalSugar.Fatal(err)
	}
}


func (p *PostgresStorage) UpdateGauge(ctx context.Context, key string, value gauge) {
	var countRetry = 1
	for i := 1; i < 6; i += 2 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		_, err := p.Exec(ctxWithTimeout, `MERGE INTO metrics_gauge AS target
		USING (VALUES ($1::text, $2::double precision)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, value)
		cancel()
		if err != nil {
			GlobalSugar.Infoln("Error update gauge:", err)
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			continue
		} else {
			break
		}
	}
}

func (p *PostgresStorage) UpdateCounter(ctx context.Context, key string, value counter) {
	var countRetry = 1
	for i := 1; i < 6; i += 2 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		_, err := p.Exec(ctxWithTimeout, `MERGE INTO metrics_counter AS target
		USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, value)
		cancel()
		if err != nil {
			GlobalSugar.Infoln("Error update counter:", err)
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			continue
		} else {
			break
		}
	}
}

func (p *PostgresStorage) AddCounter(ctx context.Context, key string, value counter) {
	var countRetry = 1	
	for i := 1; i < 6; i += 2 {
		newValue, ok := p.GetCounter(ctx, key)
		if !ok {
			GlobalSugar.Infoln("Error Get counter")
			break
		}
		newValue += value
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		_, err := p.Exec(ctxWithTimeout, `MERGE INTO metrics_counter AS target
		USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, newValue)
		cancel()
		if err != nil {
			GlobalSugar.Infoln("Error add counter:", err)
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			continue
		} else {
			break
		}
	}
}

func (p *PostgresStorage) GetCounter(ctx context.Context, key string) (counter, bool) {
	var value counter
	var Rows pgx.Rows
	var err error
	var countRetry = 1
	var cancel context.CancelFunc
	var ctxWithTimeout context.Context
	for i := 1; i < 6; i += 2 {
		ctxWithTimeout, cancel = context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		Rows, err = p.Query(ctxWithTimeout, `SELECT * FROM metrics_counter WHERE key = $1::text`, key)
		if err != nil {
			GlobalSugar.Infoln("Error get counter:", err)
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			cancel()
			continue
		} else {
			break
		}
	}
	defer Rows.Close()
	for Rows.Next() {
		err = Rows.Scan(&key, &value)
		if err != nil {
			GlobalSugar.Errorf("Error iterating over rows:: %v", err)
			return 0, false
		}
	}
	if err := Rows.Err(); err != nil {
		GlobalSugar.Errorf("Error iterating over rows: %v", err)
		return 0, false
	}
	return value, true
}

func (p *PostgresStorage) GetGauge(ctx context.Context, key string) (gauge, bool) {
	var value gauge
	var Rows pgx.Rows
	var err error
	var countRetry = 1
	var cancel context.CancelFunc
	var ctxWithTimeout context.Context
	for i := 1; i < 6; i += 2 {
		ctxWithTimeout, cancel = context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		Rows, err = p.Query(ctxWithTimeout, `SELECT * FROM metrics_gauge WHERE key = $1::text`, key)
		if err != nil {
			GlobalSugar.Infoln("Error get counter:", err)
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			cancel()
			continue
		} else {
			break
		}
	}
	defer Rows.Close()
	for Rows.Next() {
		err = Rows.Scan(&key, &value)
		if err != nil {
			GlobalSugar.Errorf("Error iterating over rows:: %v", err)
			return 0, false
		}
	}
	if err := Rows.Err(); err != nil {
		GlobalSugar.Errorf("Error iterating over rows: %v", err)
		return 0, false
	}
	return value, true
}

func (p *PostgresStorage) Ping(ctx context.Context) error {
	err := p.db.Ping(ctx)
	return err
}