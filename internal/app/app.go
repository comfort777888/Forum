package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"forum/config"
	"forum/internal/controller"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/usecase"
	"forum/pkg/database"
)

type App struct {
	config *config.Config
}

func New(conf *config.Config) *App {
	return &App{
		config: conf,
	}
}

func (a *App) Start() {
	// initialise repository
	db, err := database.InitDB(a.config)
	if err != nil {
		log.Fatalf("app - start - repository init error: %v\n", err)
	}

	// initialise tables - what the best way to place creation of tables?
	if err := database.CreateTables(db); err != nil {
		log.Fatalf("app - start - create tables: %v\n", err)
	}

	// repository layer
	userRepository := repository.NewRepository(db, a.config)
	// usecase layer
	useCase := usecase.NewUseCase(userRepository)
	// handler
	handler := controller.NewHandler(useCase)

	router := controller.SetupRouter(handler)

	server := server.NewServer(a.config, router)

	// waiting signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	// checking if we receiving signal to shut down server
	select {
	case s := <-interrupt:
		log.Printf("app: start: signal: " + s.String())
	case err = <-server.Notify():
		log.Printf("app: start: server.Notify: %v", err)
	}
	// shutdown server
	err = server.Shutdown()
	if err != nil {
		log.Printf("app: start: server.Shutdown: %v\n", err)
	}
}
