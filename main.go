package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/HotCodeGroup/warscript-utils/models"
	"github.com/jackc/pgx"
	"google.golang.org/grpc"

	"github.com/jcftang/logentriesrus"

	"github.com/sirupsen/logrus"
)

var authManager models.AuthClient
var logger *logrus.Logger

func main() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// собираем логи в хранилище
	le, err := logentriesrus.NewLogentriesrusHook(os.Getenv("LOGENTRIESRUS_TOKEN"))
	if err != nil {
		log.Printf("can not create logrus logger %s", err)
		return
	}
	logger.AddHook(le)

	dbPort, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 16)
	if err != nil {
		logger.Errorf("incorrect database port: %s", err.Error())
		return
	}

	pgxConn, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Database: os.Getenv("DB_NAME"),
			Port:     uint16(dbPort),
		},
	})
	if err != nil {
		logger.Errorf("cant connect to postgresql database: %s", err.Error())
		return
	}
	defer pgxConn.Close()

	games := &GamesManager{}
	grpcPort := os.Getenv("GRPC_PORT")
	listenGRPCPort, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Errorf("grpc port listener error: %s", err)
		return
	}

	serverGRPCGames := grpc.NewServer()
	models.RegisterGamesServer(serverGRPCGames, games)

	logger.Infof("Games gRPC service successfully started at port %s", grpcPort)
	go func() {
		if err := serverGRPCGames.Serve(listenGRPCPort); err != nil {
			logger.Fatalf("Games gRPC service failed at port %s", grpcPort)
			os.Exit(1)
		}
	}()

	grcpConn, err := grpc.Dial(
		os.Getenv("GRPC_AUTH_ADDRESS"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to auth grpc client")
	}
	defer grcpConn.Close()

	authManager = models.NewAuthClient(grcpConn)

	r := mux.NewRouter().PathPrefix("/v1").Subrouter()
	r.HandleFunc("/games", GetGameList).Methods("GET")
	r.HandleFunc("/games/{game_slug}", GetGame).Methods("GET")
	r.HandleFunc("/games/{game_slug}/leaderboard", GetGameLeaderboard).Methods("GET")
	r.HandleFunc("/games/{game_slug}/leaderboard/count", GetGameTotalPlayers).Methods("GET")

	httpPort := os.Getenv("HTTP_PORT")
	logger.Infof("Games HTTP service successfully started at port %s", httpPort)
	err = http.ListenAndServe(":"+httpPort, r)
	if err != nil {
		logger.Errorf("cant start main server. err: %s", err.Error())
		return
	}
}
