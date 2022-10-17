package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rohitchauraisa1997/calorie-tracker/routes"
)

func main() {
	fmt.Println("Running calorie-tracker")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	router.GET("/", routes.Ping)
	router.POST("/entry/add", routes.AddEntry)
	router.GET("/entries", routes.GetEntries)
	router.GET("/entry", routes.GetEntryById)
	router.GET("/ingredient", routes.GetEntriesByIngredient)

	router.PUT("/entry/update", routes.UpdateEntry)
	router.PUT("/ingredient/update", routes.UpdateIngredients)
	router.DELETE("/entry/softdelete", routes.SoftDeleteEntry)
	router.DELETE("/entry/delete", routes.DeleteEntry)

	router.Run(":" + port)
}
