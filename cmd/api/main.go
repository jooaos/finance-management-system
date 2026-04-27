package main

import (
	"finance-management/config"
	"finance-management/internal/delivery/server"
	"finance-management/internal/repository"
	"finance-management/internal/service"
	"finance-management/tools/helpers"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var (
	logger *slog.Logger
)

func init() {
	logger = slog.Default()

	var envFile string
	flag.StringVar(&envFile, "env", "", "path to .env file")
	flag.Parse()

	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			logger.Error("could not get env file", helpers.ErrLoggingKey, err)
		}
	}
}

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		logger.Error("could not parse config", helpers.ErrLoggingKey, err)
	}

	logger.Info("Bootstrapping API",
		"DBUser", cfg.DBUser,
		"DBHost", cfg.DBHost,
		"DBName", cfg.DBName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("could not open database connection", helpers.ErrLoggingKey, err)
		panic("database open connection error")
	}

	dbConn, err := db.DB()
	if err != nil {
		logger.Error("could not get db instance", helpers.ErrLoggingKey, err)
		panic("database instance error")
	}
	defer dbConn.Close()

	userRepository := repository.NewUserRepository(db)
	categoriaRepository := repository.NewCategoriaRepository(db)
	transacoesRepository := repository.NewTransacoesRepository(db)
	orcamentosRepository := repository.NewOrcamentosRepository(db)

	usuarioService := service.NewUsuarioService(userRepository, categoriaRepository)
	transacaoService := service.NewTransacaoService(transacoesRepository, categoriaRepository, orcamentosRepository)
	categoriaService := service.NewCategoriaService(categoriaRepository)
	orcamentoService := service.NewOrcamentoService(orcamentosRepository)

	httpServer := server.NewHTTPServer(
		logger,
		usuarioService,
		transacaoService,
		categoriaService,
		orcamentoService,
	)
	router := httpServer.InitServer()

	http.Handle("/", router)

	logger.Info("starting server", "port", cfg.Port)
	err = http.ListenAndServe(cfg.Port, nil)
	if err != nil {
		logger.Error("error starting server", helpers.ErrLoggingKey, err)
		return
	}
}
