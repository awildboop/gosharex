package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/awildboop/gosharex/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteRedirect(redirects *mongo.Collection, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		identifier := ctx.Query("identifier")

		if err := redirects.FindOne(context.TODO(), bson.M{"identifier": identifier}).Err(); err != nil {
			ctx.JSON(http.StatusNotFound, &common.ErrorResponse{Message: "Identifier does not exist.", Code: http.StatusNotFound})
			return
		}

		redirects.DeleteOne(todo, bson.M{"identifier": identifier})
	}
}

func CreateRedirect(redirects *mongo.Collection, conf *common.Configuration, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		identifier := ctx.Query("identifier")

		// we are assuming the random string is unused, TODO: add checking of random identifier,
		// potentially split no provided/provided identifier into it's own if/else block since they'll be quite similar.
		if identifier == "" {
			identifier = common.RandomString(8)
		}

		newRedirect := &common.Redirect{
			Identifier:   identifier,
			Location:     strings.Trim(ctx.Query("location"), " $/^\\"),
			CreationDate: time.Now().Format(time.RFC3339),
			DeletionKey:  common.RandomString(32), // unused, todo: redo delete API to check deletion key (also api key)
			Clicks:       0,
		}

		if err := redirects.FindOne(context.TODO(), bson.M{"identifier": newRedirect.Identifier}).Err(); err == nil {
			ctx.JSON(http.StatusConflict, &common.ErrorResponse{Message: "Identifier already exists, must be unique.", Code: http.StatusConflict})
			return
		}

		_, err := redirects.InsertOne(todo, newRedirect)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &common.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
			return
		}

		ctx.JSON(http.StatusCreated, &common.RedirectCreated{
			ShortenedURL: conf.Webserver.BaseURL + "r/" + newRedirect.Identifier,
			TargetURL:    newRedirect.Location,
			DeletionURL:  "NOT IMPLEMENTED",
		})
	}
}

func GetRedirect(redirects *mongo.Collection, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		identifier := ctx.Query("identifier")
		var redirect common.Redirect

		if err := redirects.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&redirect); err != nil {
			ctx.JSON(http.StatusNotFound, &common.ErrorResponse{Message: "Identifier does not exist.", Code: http.StatusNotFound})
			return
		}

		ctx.JSON(http.StatusOK, redirect)
	}
}
