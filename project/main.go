package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rabbit/rabbit_instance"
	"rabbit/routes"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	app := routes.InitRoutes()
	app.LoadHTMLGlob("templates/*")
	numOfRabbitInstance := runtime.NumCPU() * 10
	// Initialize rabbitMQ connection pool
	rabbit_instance.InitConnPool(numOfRabbitInstance, os.Getenv("RABBITMQ_URL"))
	app.NoRoute(func(ctx *gin.Context) {
		notFoundResponse := gin.H{
			"error":   true,
			"status":  "failed",
			"message": "Page not found",
		}
		ctx.JSON(http.StatusNotFound, notFoundResponse)
	})

	server := http.Server{
		Addr:           fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:        app,
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   2 * time.Minute,
		MaxHeaderBytes: 2 << 20,
	}

	log.Println("\n\n Started \n\n ")
	go Receiver()
	log.Println("\n\n Ends \n\n ")

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Could not start app: >> ", err)
		}
	}()
	Sender()

	quitCh := make(chan os.Signal, 1)

	// Signal to quit channel on interruption or termination signal
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh // wait till there is signal, then shutdown

	fmt.Sprintln("Server about to shutdown")
	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err = server.Shutdown(cntx); err != nil {
		log.Fatal("Unexpected error during graceful shutdown", err)
	}

	select {
	case <-cntx.Done():
		fmt.Println("Server gracefully shutdown")
	}
}
