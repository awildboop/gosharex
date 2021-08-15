package main

import (
	"context"
	"log"
	"time"

	"github.com/awildboop/gosharex/common"
	"github.com/awildboop/gosharex/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	conf, err := common.LoadConfiguration("./config.yaml")
	if err != nil {
		log.Fatalf("Encountered error while loading configuration file\n%v\n", err)
	}
	cfFeatures := conf.Features
	// cfAPI := conf.Features.API

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.GetURI()))
	if err != nil {
		log.Fatalf("Encountered error while connecting to MongoDB database\n%v\n", err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Encountered error while testing MongoDB database\n%v\n", err)
	}

	redirects := client.Database(conf.MongoDB.DB).Collection("redirects")

	r := gin.Default()
	v1 := r.Group("v1")

	// r = redirect (shortener), i = image, t = text, f = file,
	// potentiall merge image/text/file into a single one since they really are all just files in the end
	if cfFeatures.EnableRedirector {
		v1.GET("/r/", handlers.HandleRedirect(redirects))
	}

	// if cfAPI.EnableAPI {
	// }
	// TODO: Load config, initiate web server
}
