package internal

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultAddr            = "localhost:8080"
	defaultStoreInterval   = 0
	defaultFileStoragePath = "/tmp/file"
	defaultRestore         = true
)

var Addr = flag.String("a", defaultAddr, "Адрес HTTP-сервера")
var StoreInterval = flag.Int("i", defaultStoreInterval, "Интервал времени")
var FileStoragePath = flag.String("f", defaultFileStoragePath, "путь до файла")
var Restore = flag.Bool("r", defaultRestore, "загружать или нет ранее сохранённые значения из указанного файла при старте сервера")

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
}
