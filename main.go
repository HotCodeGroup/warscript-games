package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/HotCodeGroup/warscript-utils/balancer"
	"github.com/HotCodeGroup/warscript-utils/logging"
	"github.com/HotCodeGroup/warscript-utils/middlewares"
	"github.com/HotCodeGroup/warscript-utils/models"
	"github.com/HotCodeGroup/warscript-utils/postgresql"
	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"

	consulapi "github.com/hashicorp/consul/api"
	vaultapi "github.com/hashicorp/vault/api"
)

var authGPRC models.AuthClient
var logger *logrus.Logger

func main() {
	// коннекстим логер
	var err error
	logger, err = logging.NewLogger(os.Stdout, os.Getenv("LOGENTRIESRUS_TOKEN"))
	if err != nil {
		log.Printf("can not create logger: %s", err)
		return
	}

	// коннектим консул
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = os.Getenv("CONSUL_ADDR")
	consul, err := consulapi.NewClient(consulConfig)
	if err != nil {
		logger.Errorf("can not connect consul service: %s", err)
		return
	}

	// коннектим волт
	vaultConfig := vaultapi.DefaultConfig()
	vaultConfig.Address = os.Getenv("VAULT_ADDR")
	vault, err := vaultapi.NewClient(vaultConfig)
	if err != nil {
		logger.Errorf("can not connect vault service: %s", err)
		return
	}
	vault.SetToken(os.Getenv("VAULT_TOKEN"))

	// получаем порты, на которых будем стартовать
	httpPort, grpcPort, err := balancer.GetPorts("warscript-games/bounds", "warscript-games", consul)
	if err != nil {
		logger.Errorf("can not find empry port: %s", err)
		return
	}

	// получаем конфиг на постгрес и стартуем
	postgreConf, err := vault.Logical().Read("warscript-games/postgres")
	if err != nil || postgreConf == nil || len(postgreConf.Warnings) != 0 {
		logger.Errorf("can read warscript-games/postges key: %+v; %+v", err, postgreConf)
		return
	}

	pgxConn, err = postgresql.Connect(postgreConf.Data["user"].(string), postgreConf.Data["pass"].(string),
		postgreConf.Data["host"].(string), postgreConf.Data["port"].(string), postgreConf.Data["database"].(string))
	if err != nil {
		logger.Errorf("can not connect to postgresql database: %s", err.Error())
		return
	}
	defer pgxConn.Close()

	// коннектимся к серверу warscript-users по grpc
	authGPRCConn, err := balancer.ConnectClient(consul, "warscript-users-grpc")
	if err != nil {
		logger.Errorf("can not connect to auth grpc: %s", err.Error())
		return
	}
	defer authGPRCConn.Close()
	authGPRC = models.NewAuthClient(authGPRCConn)

	// регаем http сервис
	httpServiceID := fmt.Sprintf("warscript-games-http:%d", httpPort)
	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      httpServiceID,
		Name:    "warscript-games-http",
		Port:    httpPort,
		Address: "127.0.0.1",
	})
	if err != nil {
		logger.Errorf("can not register warscript-games-http: %s", err.Error())
		return
	}
	defer func() {
		err = consul.Agent().ServiceDeregister(httpServiceID)
		if err != nil {
			logger.Errorf("can not derigister http service: %s", err)
		}
		logger.Info("successfully derigister http service")
	}()

	// регаем grpc сервис
	grpcServiceID := fmt.Sprintf("warscript-games-grpc:%d", grpcPort)
	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      grpcServiceID,
		Name:    "warscript-games-grpc",
		Port:    grpcPort,
		Address: "127.0.0.1",
	})
	if err != nil {
		logger.Errorf("can not register warscript-games-grpc: %s", err.Error())
		return
	}
	defer func() {
		err = consul.Agent().ServiceDeregister(grpcServiceID)
		if err != nil {
			logger.Errorf("can not derigister grpc service: %s", err)
		}
		logger.Info("successfully derigister grpc service")
	}()

	// стартуем свой grpc
	games := &GamesManager{}
	listenGRPCPort, err := net.Listen("tcp", ":"+strconv.Itoa(grpcPort))
	if err != nil {
		logger.Errorf("grpc port listener error: %s", err)
		return
	}

	serverGRPCGames := grpc.NewServer()
	models.RegisterGamesServer(serverGRPCGames, games)
	logger.Infof("Games gRPC service successfully started at port %d", grpcPort)
	go func() {
		if err := serverGRPCGames.Serve(listenGRPCPort); err != nil {
			logger.Fatalf("Games gRPC service failed at port %d", grpcPort)
			os.Exit(1)
		}
	}()

	// стартуем http
	r := mux.NewRouter().PathPrefix("/v1").Subrouter()
	r.HandleFunc("/games", GetGameList).Methods("GET")
	r.HandleFunc("/games/{game_slug}", GetGame).Methods("GET")
	r.HandleFunc("/games/{game_slug}/leaderboard", GetGameLeaderboard).Methods("GET")
	r.HandleFunc("/games/{game_slug}/leaderboard/count", GetGameTotalPlayers).Methods("GET")

	logger.Infof("Games HTTP service successfully started at port %d", httpPort)
	err = http.ListenAndServe(":"+strconv.Itoa(httpPort),
		middlewares.RecoverMiddleware(middlewares.AccessLogMiddleware(r, logger), logger))
	if err != nil {
		logger.Errorf("cant start main server. err: %s", err.Error())
		return
	}
}
