package main

import (
	"divar.ir/api/routers"
	"divar.ir/internal/mongodb"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic("Can't load <.env> file.")
	}

	if os.Getenv("MODE") == "RELEASE" {
		fmt.Println("Server started by <RELEASE MODE>")
		gin.SetMode(gin.ReleaseMode)
	}
	serv := gin.Default()

	err = mongodb.Init_Mongo(os.Getenv("MONGO_ADDRESS"), os.Getenv("MONGO_DBNAME"))
	if err != nil {
		panic("Failed to connect mongodb")
	}
	routers.IncludeRouters(serv)

	servAddr := os.Getenv("SERVER_ADDRESS")
	servPort := os.Getenv("SERVER_PORT")
	if servAddr == "" || servPort == "" {
		panic("Not valid <SERVER_ADDRESS> or <SERVER_PORT> in <.env> file.")
	}

	err = serv.Run(servAddr + ":" + servPort)
	if err != nil {
		panic("Failed to start server.")
	}

}
