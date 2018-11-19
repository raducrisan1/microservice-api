package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raducrisan1/microservice-api/tradesuggest"
	"google.golang.org/grpc"
)

func main() {
	addr := os.Getenv("REPORTS_GRPC_ADDR")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	tradeSuggestClient := tradesuggest.NewTradeSuggestServiceClient(conn)
	router := gin.Default()
	router.GET("/api/:stockname", func(c *gin.Context) {
		stockName := c.Param("stockname")
		req := &tradesuggest.TradeSuggestRequest{
			Resolution: 300}
		if res, err := tradeSuggestClient.GetSuggestions(c, req); err != nil {
			log.Printf("Could not obtain data from the gRPC service TradeSuggest: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": err})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprintf("Rating data for %v", stockName),
				"data":   res.Suggestions})
		}
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"result": "success"})
	})

	srv := &http.Server{
		Addr:    ":3030",
		Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown microservice-api...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown: ", err)
	}
	log.Println("Microservice-api has been stopped")
}
