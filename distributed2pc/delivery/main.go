package main

import (
	"log"

	delivery "github.com/amrishshah/distributed2pc/delivery/svc"
	"github.com/amrishshah/distributed2pc/io"
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

	router.POST("/delivery/agent/reserve", func(ctx *gin.Context) {
		agent, err := delivery.ReverseAgent()

		if err != nil {
			log.Print(err.Error())
			ctx.JSON(429, gin.H{"error": err.Error()})
			return
		} else {
			ctx.IndentedJSON(200, agent)
		}
	})

	router.POST("/delivery/agent/book", func(ctx *gin.Context) {

		var req delivery.Agent
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatus(400)
		}
		// var s []byte
		// ctx.Request.Body.Read(s)
		//log.Println(string(*req.OrderId))

		agent, err := delivery.BookAgent(string(*req.OrderId))

		if err != nil {
			ctx.JSON(429, err)

		} else {
			ctx.IndentedJSON(200, agent)
		}
	})

	router.Run("localhost:8080")
}

// // getAlbums responds with the list of all albums as JSON.
// func getAlbums(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, albums)
// }

// // postAlbums adds an album from JSON received in the request body.
// func postAlbums(c *gin.Context) {
// 	var newAlbum album

// 	// Call BindJSON to bind the received JSON to
// 	// newAlbum.
// 	if err := c.BindJSON(&newAlbum); err != nil {
// 		return
// 	}

// 	// Add the new album to the slice.
// 	albums = append(albums, newAlbum)
// 	c.IndentedJSON(http.StatusCreated, newAlbum)
// }

// // getAlbumByID locates the album whose ID value matches the id
// // parameter sent by the client, then returns that album as a response.
// func getAlbumByID(c *gin.Context) {
// 	id := c.Param("id")

// 	// Loop through the list of albums, looking for
// 	// an album whose ID value matches the parameter.
// 	for _, a := range albums {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }
