package internal

import (
	"context"
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
	_, err := DB.ExecContext(ctx, `MERGE INTO metrics_gauge AS target
USING (VALUES ($1::text, $2::double precision)) AS source (key, value)
ON (target.key = source.key)
WHEN MATCHED THEN
 UPDATE SET value = source.value
WHEN NOT MATCHED THEN
 INSERT (key, value) VALUES (source.key, source.value)`, key, value)
	if err != nil {
		GlobalSugar.Infoln(err)
	}
}

func UpdateCounterSQL(ctx context.Context, key string, value counter) {
	_, err := DB.ExecContext(ctx, `MERGE INTO metrics_counter AS target
USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
ON (target.key = source.key)
WHEN MATCHED THEN
 UPDATE SET value = source.value
WHEN NOT MATCHED THEN
 INSERT (key, value) VALUES (source.key, source.value)`, key, value)
	if err != nil {
		GlobalSugar.Infoln(err)
	}
}

func AddCounterSQL(ctx context.Context, key string, value counter) {
	newValue, ok := GetCounterSQL(ctx, key)
	if !ok {
		GlobalSugar.Infoln("Error Get counter")
		return
	}
	newValue += value
	for i := 1; i < 6; i+=2 {
		_, err := DB.ExecContext(ctx, `MERGE INTO metrics_counter AS target
		USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, newValue)
		if err != nil {
			GlobalSugar.Infoln("Error add counter:",err)
			time.Sleep(time.Second * time.Duration(i))
		} else {
			i = 6
		}
	}
}

func GetCounterSQL(ctx context.Context, key string) (counter, bool) {
	var value counter
	for i := 1; i < 6; i+=2 {
		rows, err := DB.QueryContext(ctx, `SELECT * FROM metrics_counter WHERE key = $1::text`, key)
		if err != nil {
			GlobalSugar.Errorln("Error querying database:", err)
			time.Sleep(time.Second * time.Duration(i))
			continue
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&key, &value)
			if err != nil {
				GlobalSugar.Errorf("Error iterating over rows:: %v", err)
				return 0, false
			}
		}
		if err := rows.Err(); err != nil {
			GlobalSugar.Errorf("Error iterating over rows: %v", err)
			return 0, false
		}
		return value, true
	}
	// Если все попытки исчерпаны и значение не получено, завершаем с ошибкой
	GlobalSugar.Errorf("All retries exhausted, exiting...")
	return 0, false
}

func GetGaugeSQL(ctx context.Context, key string) (gauge, bool) {
	var value gauge
	for i := 1; i < 6; i+=2 {
		rows, err := DB.QueryContext(ctx, `SELECT * FROM metrics_gauge WHERE key = $1::text`, key)
		if err != nil {
			GlobalSugar.Errorln("Error querying database:", err)
			time.Sleep(time.Second * time.Duration(i))
			continue
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&key, &value)
			if err != nil {
				GlobalSugar.Errorf("Error iterating over rows:: %v", err)
				return 0, false
			}
		}
		if err := rows.Err(); err != nil {
			GlobalSugar.Errorf("Error iterating over rows: %v", err)
			return 0, false
		}
		return value, true
	}
	// Если все попытки исчерпаны и значение не получено, завершаем с ошибкой
	GlobalSugar.Errorf("All retries exhausted, exiting...")
	return 0, false
}
