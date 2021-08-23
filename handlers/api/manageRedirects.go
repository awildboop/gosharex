package api

import (
	"context"
	"net/http"

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

func CreateRedirect(redirects *mongo.Collection, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		newRedirect := &common.Redirect{
			Identifier:   ctx.Query("identifier"),
			Location:     ctx.Query("location"),
			CreationDate: "",
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

		ctx.Writer.WriteHeader(http.StatusCreated)
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
