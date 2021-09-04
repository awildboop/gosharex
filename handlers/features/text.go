package features

import (
	"context"
	"net/http"

	"github.com/awildboop/gosharex/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleText(texts *mongo.Collection, conf *common.Configuration, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		identifier := ctx.Param("identifier")[1:]

		var text common.Text

		if err := texts.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&text); err != nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "https://awildboop.com")
			return
		}

		texts.UpdateOne(todo, bson.M{"identifier": identifier}, bson.M{"$set": bson.M{"views": text.Views + 1}})

		ctx.HTML(http.StatusOK, "text.html", gin.H{
			"url":     ctx.Request.URL.RawPath,
			"title":   conf.Pages.PageTitles.TextTitle,
			"content": text.Content,
			"preview": text.Preview,
		})
	}
}
