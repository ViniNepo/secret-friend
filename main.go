package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ViniNepo/secretfriend/config"
	"github.com/ViniNepo/secretfriend/handler"
	"github.com/ViniNepo/secretfriend/services"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/urfave/negroni"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := initializeDatabase(cfg.AppConfig.DBHost, cfg.AppConfig.DBPort, cfg.AppConfig.DBUser, cfg.AppConfig.DBPassword, cfg.AppConfig.DBName)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	emailService := services.NewEmailService(cfg.AppConfig.FromEmail, cfg.AppConfig.FromEmailSMTP, cfg.AppConfig.FromEmailPassword, cfg.AppConfig.SMTPAddrress)
	friendService := services.NewFriendService(emailService, db)

	// Configurar o roteador Gorilla Mux
	router := mux.NewRouter()

	// Configurar os handlers
	handler.CreatePingHandlers(router)
	friendHandlers := handler.NewFriendHandlers(friendService)
	friendHandlers.CreateFriendHandlers(router)

	// Configurar CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"})

	// Configurar Negroni como middleware
	n := negroni.Classic()
	n.UseHandler(router)

	// Iniciar o servidor com suporte a CORS
	err = http.ListenAndServe(":"+cfg.ServerPort, handlers.CORS(originsOk, headersOk, methodsOk)(n))
	if err != nil {
		log.Println("error starting server:", err)
		return
	}
	log.Println("Server started on port", cfg.ServerPort)
}

func initializeDatabase(host, port, user, password, dbname string) (*sql.DB, error) {
	// Monta a string de conexão
	dsn := "postgresql://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	time.Sleep(5 * time.Second)
	// Abre a conexão com o banco
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Testa a conexão
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
