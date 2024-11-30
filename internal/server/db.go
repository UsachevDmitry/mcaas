package internal

import (
	"context"
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
	newValue, _ := GetCounterSQL(ctx, key)
	newValue += value
	_, err := DB.ExecContext(ctx, `MERGE INTO metrics_counter AS target
USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
ON (target.key = source.key)
WHEN MATCHED THEN
 UPDATE SET value = source.value
WHEN NOT MATCHED THEN
 INSERT (key, value) VALUES (source.key, source.value)`, key, newValue)
	if err != nil {
		GlobalSugar.Infoln(err)
	}
}

func GetCounterSQL(ctx context.Context, key string) (counter, bool) {
	var value counter
	rows, err := DB.QueryContext(ctx, `SELECT * FROM metrics_counter WHERE key = $1::text`, key)
	if err != nil {
		GlobalSugar.Panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&key, &value)
		if err != nil {
			return 0, false
		}
	}
	if err := rows.Err(); err != nil {
		GlobalSugar.Errorf("Ошибка итерации по строкам: %v", err)
	}
	return value, true
}

func GetGaugeSQL(ctx context.Context, key string) (gauge, bool) {
	var value gauge
	rows, err := DB.QueryContext(ctx, `SELECT * FROM metrics_gauge WHERE key = $1::text`, key)
	if err != nil {
		GlobalSugar.Panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&key, &value)
		if err != nil {
			return 0, false
		}
	}
	if err := rows.Err(); err != nil {
		GlobalSugar.Errorf("Ошибка итерации по строкам: %v", err)
	}
	return value, true
}
