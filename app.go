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
	conn, err := grpc.Dial("localhost:3070", grpc.WithInsecure())
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
			log.Fatalf("Could not obtain data from the gRPC service TradeSuggest: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": err})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprintf("Rating data for %v", stockName),
				"data":   res.Suggestions})
		}
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	fmt.Println("\nmicroservice-api has been stopped")
}
