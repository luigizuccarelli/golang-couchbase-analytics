package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-couchbase-analytics/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-couchbase-analytics/pkg/handlers"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-couchbase-analytics/pkg/validator"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
)

func startHttpServer(con connectors.Clients) *http.Server {
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/sankeydata", func(w http.ResponseWriter, req *http.Request) {
		handlers.SankeyChartHandler(w, req, con)
	}).Methods("POST")

	r.HandleFunc("/api/v2/sys/info/isalive", handlers.IsAlive).Methods("GET")

	sh := http.StripPrefix("/api/v2/web/", http.FileServer(http.Dir("./charts/")))
	r.PathPrefix("/api/v2/web/").Handler(sh)
	http.Handle("/", r)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			con.Error("Httpserver: ListenAndServe() error: " + err.Error())
		}
	}()

	return srv
}

func main() {

	var logger *simple.Logger

	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}

	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}

	conn := connectors.NewClientConnections(logger)

	defer conn.Close()

	srv := startHttpServer(conn)
	logger.Info("Starting server on port " + srv.Addr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan

	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
	logger.Info("Server shutdown successfully")
	os.Exit(code)
}
