package main

import (
	"github.com/UsachevDmitry/mcaas/cmd/server/internal"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	// "sync"
)

func main() {
	// var wg sync.WaitGroup

	internal.GetConfig()
	internal.ImportDataFromFile(*internal.FileStoragePath, *internal.Restore)

	// wg.Add(1)
	// go func() {
	go internal.SaveDataInFile(time.Duration(*internal.StoreInterval), *internal.FileStoragePath)
	// 	defer wg.Done()
	// }()

	router := mux.NewRouter()
	router.HandleFunc("/", internal.WithLoggingGet(internal.GzipHandle(internal.HandleIndex()))).Methods(http.MethodGet)
	router.HandleFunc("/update/", internal.WithLoggingGet(internal.GzipHandle(internal.HandlePostMetricsJSON()))).Methods(http.MethodPost)
	router.HandleFunc("/update/{type}/{name}/{value}", internal.WithLoggingGet(internal.GzipHandle(internal.HandlePostMetrics()))).Methods(http.MethodPost)
	router.HandleFunc("/value/", internal.WithLoggingGet(internal.GzipHandle(internal.HandleGetMetricsJSON()))).Methods(http.MethodPost)
	router.HandleFunc("/value/{type}/{name}", internal.WithLoggingGet(internal.GzipHandle(internal.HandleGetValue()))).Methods(http.MethodGet)
	
	internal.Logger()
	internal.GlobalSugar.Infow(
		"Starting server",
		"ADDRESS", *internal.Addr,
		"STORE_INTERVAL", *internal.StoreInterval,
		"FILE_STORAGE_PATH", *internal.FileStoragePath, 
		"RESTORE", *internal.Restore, 
	)
	
	if err := http.ListenAndServe(*internal.Addr, router); err != nil {
		internal.GlobalSugar.Fatalw(err.Error(), "event", "start server")
	}
}
