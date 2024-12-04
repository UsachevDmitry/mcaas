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
	poolConfig.MaxConns = 10 // Максимальное количество соединений в пуле
	poolConfig.MinConns = 5 // Минимальное количество поддерживаемых соединений

	p.db, err = pgxpool.New(context.Background(), *DatabaseDsn)
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
	_, err := p.Exec(ctx, `CREATE TABLE IF NOT EXISTS metrics_gauge (
		"key" TEXT,
		"value" DOUBLE PRECISION
	)`)
	if err != nil {
		GlobalSugar.Fatal(err)
	}
}
func (p *PostgresStorage) CreateTableCounter(ctx context.Context) {
	_, err := p.Exec(ctx,`CREATE TABLE IF NOT EXISTS metrics_counter (
		"key" TEXT,
		"value" BIGINT
	)`)
	if err != nil {
		GlobalSugar.Fatal(err)
	}
}


func (p *PostgresStorage) UpdateGauge(ctx context.Context, key string, value gauge) {
	for i := 1; i < 6; i += 2 {
		p.Connect()
		_, err := p.Exec(ctx, `MERGE INTO metrics_gauge AS target
		USING (VALUES ($1::text, $2::double precision)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, value)
		p.Close()
		if err != nil {
			GlobalSugar.Infoln("Error update gauge:", err)
			GlobalSugar.Infof("Retry after %v second", i)
			time.Sleep(time.Second * time.Duration(i))
			continue
		} else {
			break
		}
	}
}

func (p *PostgresStorage) UpdateCounter(ctx context.Context, key string, value counter) {
	for i := 1; i < 6; i += 2 {
		p.Connect()
		_, err := p.Exec(ctx, `MERGE INTO metrics_counter AS target
		USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, value)
		p.Close()
		if err != nil {
			GlobalSugar.Infoln("Error update counter:", err)
			GlobalSugar.Infof("Retry after %v second", i)
			time.Sleep(time.Second * time.Duration(i))
			continue
		} else {
			break
		}
	}
}

func (p *PostgresStorage) AddCounter(ctx context.Context, key string, value counter) {	
	for i := 1; i < 6; i += 2 {
		newValue, ok := p.GetCounter(ctx, key)
		if !ok {
			GlobalSugar.Infoln("Error Get counter")
			break
		}
		newValue += value
		_, err := p.Exec(ctx, `MERGE INTO metrics_counter AS target
		USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, newValue)
		if err != nil {
			GlobalSugar.Infoln("Error add counter:", err)
			GlobalSugar.Infof("Retry after %v second", i)
			time.Sleep(time.Second * time.Duration(i))
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
	for i := 1; i < 6; i += 2 {
		Rows, err = p.Query(ctx, `SELECT * FROM metrics_counter WHERE key = $1::text`, key)
		if err != nil {
			GlobalSugar.Infoln("Error get counter:", err)
			if i == 5 {
				GlobalSugar.Errorln("All retries exhausted, exiting...")
				return 0, false
			}
			GlobalSugar.Infof("Retry after %v second", i)
			time.Sleep(time.Second * time.Duration(i))
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
	for i := 1; i < 6; i += 2 {
		Rows, err = p.Query(ctx, `SELECT * FROM metrics_gauge WHERE key = $1::text`, key)
		if err != nil {
			GlobalSugar.Infoln("Error get gauge:", err)
			if i == 5 {
				GlobalSugar.Errorln("All retries exhausted, exiting...")
				return 0, false
			}
			GlobalSugar.Infof("Retry after %v second", i)
			time.Sleep(time.Second * time.Duration(i))
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