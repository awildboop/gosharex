package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleRedirect(redirects *mongo.Collection) func(*gin.Context) {
	return func(ctx *gin.Context) {

	}
}
