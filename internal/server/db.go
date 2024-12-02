package internal

import (
	"context"
	"database/sql"
	"time"
)

func CreateTables(ctx context.Context) {
	_, err := DB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS metrics_gauge (
		"key" TEXT,
		"value" DOUBLE PRECISION
	)`)
	if err != nil {
		GlobalSugar.Fatal(err)
	}

	_, err = DB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS metrics_counter (
		"key" TEXT,
		"value" BIGINT
	)`)
	if err != nil {
		GlobalSugar.Fatal(err)
	}
}

func UpdateGaugeSQL(ctx context.Context, key string, value gauge) {
	var countRetry = 1	
	for i := 1; i < 6; i += 2 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		DBMutex.Lock()
		defer DBMutex.Unlock()
		_, err := DB.ExecContext(ctxWithTimeout, `MERGE INTO metrics_gauge AS target
		USING (VALUES ($1::text, $2::double precision)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, value)
		cancel()
		if err != nil {
			GlobalSugar.Infoln("Error update gauge:", err)
			if i == 5 {
				GlobalSugar.Errorln("All retries exhausted, exiting...")
				break
			}
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			continue
		}
	}
}

func UpdateCounterSQL(ctx context.Context, key string, value counter) {
	var countRetry = 1	
	for i := 1; i < 6; i += 2 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		DBMutex.Lock()
		defer DBMutex.Unlock()
		_, err := DB.ExecContext(ctxWithTimeout, `MERGE INTO metrics_counter AS target
		USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, value)
		cancel()
		if err != nil {
			GlobalSugar.Infoln("Error update counter:", err)
			if i == 5 {
				GlobalSugar.Errorln("All retries exhausted, exiting...")
				break
			}
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			continue
		}
	}
}

func AddCounterSQL(ctx context.Context, key string, value counter) {
	var countRetry = 1	
	for i := 1; i < 6; i += 2 {
		newValue, ok := GetCounterSQL(ctx, key)
		if !ok {
			GlobalSugar.Infoln("Error Get counter")
			break
		}
		newValue += value
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		DBMutex.Lock()
		defer DBMutex.Unlock()
		_, err := DB.ExecContext(ctxWithTimeout, `MERGE INTO metrics_counter AS target
		USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, newValue)
		cancel()
		if err != nil {
			GlobalSugar.Infoln("Error add counter:", err)
			if i == 5 {
				GlobalSugar.Errorln("All retries exhausted, exiting...")
				break
			}
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			continue
		}
	}
}

func GetCounterSQL(ctx context.Context, key string) (counter, bool) {
	var value counter
	var Rows *sql.Rows
	var err error
	var countRetry = 1	
	for i := 1; i < 6; i += 2 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		Rows, err = DB.QueryContext(ctxWithTimeout, `SELECT * FROM metrics_counter WHERE key = $1::text`, key)
		cancel()
		if err != nil {
			GlobalSugar.Infoln("Error get counter:", err)
			if i == 5 {
				GlobalSugar.Errorln("All retries exhausted, exiting...")
				return 0, false
			}
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			continue
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

func GetGaugeSQL(ctx context.Context, key string) (gauge, bool) {
	var value gauge
	var Rows *sql.Rows
	var err error
	var countRetry = 1
	for i := 1; i < 6; i += 2 {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(i)*time.Second)
		defer cancel()
		Rows, err = DB.QueryContext(ctxWithTimeout, `SELECT * FROM metrics_gauge WHERE key = $1::text`, key)
		cancel()
		if err != nil {
			GlobalSugar.Infoln("Error get counter:", err)
			if i == 5 {
				GlobalSugar.Errorln("All retries exhausted, exiting...")
				return 0, false
			}
			GlobalSugar.Infof("Retry %v...", countRetry)
			countRetry++
			continue
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
