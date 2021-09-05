package main

import (
	"context"
	"log"
	"time"

	"github.com/awildboop/gosharex/common"
	"github.com/awildboop/gosharex/handlers/api"
	"github.com/awildboop/gosharex/handlers/features"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	todo := context.TODO()
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
	texts := client.Database(conf.MongoDB.DB).Collection("texts")
	images := client.Database(conf.MongoDB.DB).Collection("images")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// r = redirect (shortener), i = image, t = text, f = file,
	// potentially merge image/text/file into a single one since they really are all just files in the end
	if cfFeatures.EnableRedirector {
		r.GET("/r/*identifier", features.HandleRedirect(redirects, conf, todo))
		r.POST("/r", api.CreateRedirect(redirects, conf, todo))
		r.PUT("/r", api.CreateRedirect(redirects, conf, todo))
	}

	if cfFeatures.EnableText {
		r.GET("/t/*identifier", features.HandleText(texts, conf, todo))
		r.POST("/t", api.CreateText(texts, conf, todo))
		r.PUT("/t", api.CreateText(texts, conf, todo))
	}

	if cfFeatures.EnableImages {
		r.GET("/i/*identifier", features.HandleImage(images, conf, todo))
		r.GET("/ri/*identifier", features.HandleRawImage(images, conf, todo))
		r.POST("/i", api.CreateImage(images, conf, todo))
		r.PUT("/i", api.CreateImage(images, conf, todo))
	}

	if cfFeatures.API.EnableAPI {
		apiFeatures := cfFeatures.API
		v1 := r.Group("v1")

		if apiFeatures.ManageRedirects {
			v1.GET("/r", api.GetRedirect(redirects, todo)) // Returns information & stats about the redirect
			v1.DELETE("/r", api.DeleteRedirect(redirects, todo))
		}

		// if apiFeatures.ManageText {
		// }

		// if apiFeatures.ManageImages {
		// }
	}

	r.Run(conf.GetWebserverAddress())
}
