package main

import (
	"database/sql"
	"log"

	"github.com/SIRIUS-webkit/crud-app/internal/database"
	"github.com/SIRIUS-webkit/crud-app/internal/env"
	_ "modernc.org/sqlite"
)

// @title Go Gin Rest API
// @version 1.0
// @description A rest API in Go using Gin framework.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

// Apply the security definition to your endpoints
// @security BearerAuth

type application struct{
	port int
	jwtSecret string
	models database.Models
}

func main() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil{
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModels(db)
    app := &application{
		port: env.GetEnvInt("port", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "lamote-secret-001"),
		models: models,
	}

	if err := serve(app); err != nil {
		log.Fatal(err)
	}

}