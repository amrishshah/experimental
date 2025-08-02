package main

import (
	"log"

	"github.com/amrishshah/distributed2pc/io"
	store "github.com/amrishshah/distributed2pc/store/svc"
	"github.com/gin-gonic/gin"
)

func main() {

	// Replace with your DB details
	driver := "mysql"
	dsn := "root:password@tcp(127.0.0.1:6306)/demo_12"

	// Initialize DB
	if err := io.InitDB(driver, dsn); err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.POST("/store/packet/reserve", func(ctx *gin.Context) {

		var req store.Packet
		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Println(err.Error())
			ctx.AbortWithStatus(400)
			return
		}
		log.Print(req.FoodId)
		agent, err := store.ReversePacket(req.FoodId)

		if err != nil {
			log.Print(err.Error())
			ctx.JSON(429, gin.H{"error": err.Error()})
			return
		} else {
			ctx.IndentedJSON(200, agent)
		}
	})

	router.POST("/store/packet/book", func(ctx *gin.Context) {

		var req store.Packet
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatus(400)
		}
		// var s []byte
		// ctx.Request.Body.Read(s)
		//log.Println(string(*req.OrderId))

		agent, err := store.BookPacket(string(*req.OrderId), req.FoodId)

		if err != nil {
			ctx.JSON(429, err)

		} else {
			ctx.IndentedJSON(200, agent)
		}
	})

	router.Run("localhost:8081")
}
