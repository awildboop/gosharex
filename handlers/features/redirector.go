package handlers

import (
	"context"
	"net/http"

	"github.com/awildboop/gosharex/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleRedirect(redirects *mongo.Collection, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		identifier := ctx.Param("identifier")[1:]

		var redirect common.Redirect

		if err := redirects.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&redirect); err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "https://awildboop.com")
			return
		}

		redirects.UpdateOne(todo, bson.M{"identifier": identifier}, bson.M{"$set": bson.M{"clicks": redirect.Clicks + 1}})
		ctx.Redirect(http.StatusTemporaryRedirect, redirect.Location)
	}
}
