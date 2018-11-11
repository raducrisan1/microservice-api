package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raducrisan1/microservice-api/tradesuggest"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:3015", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	tradeSuggestClient := tradesuggest.NewTradeSuggestServiceClient(conn)
	apiSrv := gin.Default()
	apiSrv.GET("/api/:stockname", func(c *gin.Context) {
		stockName := c.Param("stockname")
		startTime, _ := time.Parse(time.RFC3339, "2018-11-10 09:30Z")
		endTime, _ := time.Parse(time.RFC3339, "2018-11-10 10:00Z")
		req := &tradesuggest.TradeSuggestRequest{
			Resolution: 300
		}		
		if res, err := tradeSuggestClient.GetSuggestions(c, req); err != nil {
			log.Fatalf("Could not obtain data from the gRPC service TradeSuggest: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": err})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprintf("Will bring some data for %v", stockName),
				"data":   res})
		}
	})
	if err := apiSrv.Run(":3030"); err != nil {
		log.Fatalf("Could not start the gin server: %v", err)
	}
}
