package features

import (
	"context"
	"net/http"

	"github.com/awildboop/gosharex/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleImage(images *mongo.Collection, conf *common.Configuration, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		identifier := ctx.Param("identifier")[1:]

		var image common.Image

		if err := images.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&image); err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, conf.Webserver.NotFoundRedirect)
			return
		}

		ctx.HTML(http.StatusOK, "image.html", gin.H{
			"url":            ctx.Request.URL.RawPath,
			"title":          conf.Pages.PageTitles.ImageTitle,
			"image_location": conf.Webserver.BaseURL + "ri/" + identifier,
		})
	}
}

func HandleRawImage(images *mongo.Collection, conf *common.Configuration, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		identifier := ctx.Param("identifier")[1:]

		var image common.Image

		if err := images.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&image); err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, conf.Webserver.NotFoundRedirect)
			return
		}

		ctx.File("./uploads/" + image.FileLocation)
	}
}
