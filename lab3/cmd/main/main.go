package main

import (
	"fmt"

	_ "lab3/docs"
	"lab3/internal/app/config"
	"lab3/internal/app/dsn"
	"lab3/internal/app/handler"
	"lab3/internal/app/repository"
	"lab3/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Big O Request API
// @version 1.0
// @description API для управления расчётами времени и сложности Классов сложности
// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email support@bigorequest.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.NewRepository(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
