package internal

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultAddr            = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "/tmp/file"
	defaultRestore         = true
	defaultDatabaseDsn     = ""
	//defaultDatabaseDsn = "host=localhost user=postgres password=P@ssw0rd dbname=test"
)

var Addr = flag.String("a", defaultAddr, "Адрес HTTP-сервера")
var StoreInterval = flag.Int("i", defaultStoreInterval, "Интервал времени")
var FileStoragePath = flag.String("f", defaultFileStoragePath, "путь до файла")
var Restore = flag.Bool("r", defaultRestore, "загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
var DatabaseDsn = flag.String("d", defaultDatabaseDsn, "строка подключения к БД")

func GetConfig() {
	flag.Parse()

	addrEnv := os.Getenv("ADDRESS")
	if addrEnv != "" {
		*Addr = addrEnv
	}

	storeIntervalEnv := os.Getenv("STORE_INTERVAL")
	if storeIntervalEnv != "" {
		i, err := strconv.Atoi(storeIntervalEnv)
		if err != nil {
			GlobalSugar.Fatal(err)
		}
		*StoreInterval = i
	}

	fileStoragePathEnv := os.Getenv("FILE_STORAGE_PATH")
	if fileStoragePathEnv != "" {
		*FileStoragePath = fileStoragePathEnv
	}

	restoreEnv := os.Getenv("RESTORE")
	if restoreEnv != "" {
		b, err := strconv.ParseBool(restoreEnv)
		if err != nil {
			GlobalSugar.Fatal(err)
		}
		*Restore = b
	}

	databaseDsnEnv := os.Getenv("DATABASE_DSN")
	if databaseDsnEnv != "" {
		*DatabaseDsn = databaseDsnEnv
	}
	
	if *DatabaseDsn != defaultDatabaseDsn {
		FlagUsePosgresSQL = true
	} else {
		FlagUsePosgresSQL = false
	}
}
