package internal

import (
	"database/sql"
	"time"
)
func (p *PostgresStorage) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(query, args...)
}

func (p *PostgresStorage) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.Query(query, args...)
}

// func (p *PostgresStorage) Open(driverName, dataSourceName string) (*sql.DB, error) {
// 	return p.db.Open(driverName, dataSourceName)
// }

func (p *PostgresStorage) Connect() error {
	var err error
	p.db, err = sql.Open("pgx", *DatabaseDsn)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) Close() error {
	return p.db.Close()
}

func (p *PostgresStorage) CreateTables() {
	_, err := p.Exec(`CREATE TABLE IF NOT EXISTS metrics_gauge (
		"key" TEXT,
		"value" DOUBLE PRECISION
	)`)
	if err != nil {
		GlobalSugar.Fatal(err)
	}

	_, err = p.Exec(`CREATE TABLE IF NOT EXISTS metrics_counter (
		"key" TEXT,
		"value" BIGINT
	)`)
	if err != nil {
		GlobalSugar.Fatal(err)
	}
}

func (p *PostgresStorage) UpdateGauge(key string, value gauge) {
	for i := 1; i < 6; i += 2 {
		_, err := p.Exec(`MERGE INTO metrics_gauge AS target
		USING (VALUES ($1::text, $2::double precision)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, value)
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

func (p *PostgresStorage) UpdateCounter(key string, value counter) {
	for i := 1; i < 6; i += 2 {
		_, err := p.Exec(`MERGE INTO metrics_counter AS target
		USING (VALUES ($1::text, $2::bigint)) AS source (key, value)
		ON (target.key = source.key)
		WHEN MATCHED THEN
		UPDATE SET value = source.value
		WHEN NOT MATCHED THEN
		INSERT (key, value) VALUES (source.key, source.value)`, key, value)
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

func (p *PostgresStorage) AddCounter(key string, value counter) {	
	for i := 1; i < 6; i += 2 {
		newValue, ok := p.GetCounter(key)
		if !ok {
			GlobalSugar.Infoln("Error Get counter")
			break
		}
		newValue += value
		_, err := p.Exec(`MERGE INTO metrics_counter AS target
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

func (p *PostgresStorage) GetCounter(key string) (counter, bool) {
	var value counter
	var Rows *sql.Rows
	var err error
	for i := 1; i < 6; i += 2 {
		Rows, err = p.Query(`SELECT * FROM metrics_counter WHERE key = $1::text`, key)
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

func (p *PostgresStorage) GetGauge(key string) (gauge, bool) {
	var value gauge
	var Rows *sql.Rows
	var err error
	for i := 1; i < 6; i += 2 {
		Rows, err = p.Query(`SELECT * FROM metrics_gauge WHERE key = $1::text`, key)
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
