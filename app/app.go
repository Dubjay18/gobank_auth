package app

import (
	"fmt"
	"github.com/Dubjay18/gobank_auth/domain"
	"github.com/Dubjay18/gobank_auth/service"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func GetEnvVar() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

}

func SanityCheck() {
	if os.Getenv("SERVER_ADDRESS") == "" || os.Getenv("SERVER_PORT") == "" {
		log.Fatal("Environment variables not defined...")
	}
}

func getDbClient() *sqlx.DB {
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	constr := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPass + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"

	db, err := sqlx.Open("postgres", constr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)

	}
	return db
}
func Start() {
	SanityCheck()
	//mux := http.NewServeMux()
	dbClient := getDbClient()
	authRepository := domain.NewAuthRepository(dbClient)
	ah := AuthHandler{service.NewLoginService(authRepository, domain.GetRolePermissions())}
	r := mux.NewRouter()
	r.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)
	port := os.Getenv("SERVER_PORT")
	address := os.Getenv("SERVER_ADDRESS")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), r))
}
