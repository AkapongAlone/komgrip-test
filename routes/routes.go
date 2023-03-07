package routes

import (
	"net/http"
	 _"github.com/AkapongAlone/komgrip-test/docs"
	db "github.com/AkapongAlone/komgrip-test/database"
	_ "github.com/AkapongAlone/komgrip-test/src/beers/handler"
	"github.com/AkapongAlone/komgrip-test/src/beers/repositories"
	"github.com/gin-gonic/gin"
	"github.com/AkapongAlone/komgrip-test/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		// add header Access-Control-Allow-Origin
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})
	// add swagger
	
	handle := repositories.NewBeerHandler(db.Db)
	beerGroup := r.Group("/api/v1")
	beerGroup.GET("/beer", handle.GetBeer)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	beerGroup.Use(middlewares.LoggingInfoMiddleware())
		{beerGroup.POST("/beer", handle.PostBeer)
		beerGroup.PUT("/beer/:id", handle.UpdateBeer)
		beerGroup.DELETE("/beer/:id", handle.DeleteBeer)
		}
	r.StaticFS("/file/", http.Dir("images"))
	return r
}
