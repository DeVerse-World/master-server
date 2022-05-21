package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/hyperjiang/gin-skeleton/config"
	"github.com/hyperjiang/gin-skeleton/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := flag.String("addr", config.Server.Addr, "Address to listen and serve")
	flag.Parse()

	if config.Server.Mode == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	uiHost := os.Getenv("UI_HOST")

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{uiHost},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == uiHost
		},
		MaxAge: 12 * time.Hour,
	}))
	// app.Use(cors.Default())
	app.Static("/images", filepath.Join(config.Server.StaticDir, "img"))
	app.StaticFile("/favicon.ico", filepath.Join(config.Server.StaticDir, "img/favicon.ico"))
	app.LoadHTMLGlob(config.Server.ViewDir + "/*")
	app.MaxMultipartMemory = config.Server.MaxMultipartMemory << 20

	router.Route(app)

	// Listen and Serve
	// app.Use(cors.Default())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{uiHost},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == uiHost
		},
		MaxAge: 12 * time.Hour,
	}))
	app.Run(*addr)
}
