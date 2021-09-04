package api

import (
	"context"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/awildboop/gosharex/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const previewLength = 27

func makePreview(content string) (p string) {
	p = strconv.Quote(content)
	p = strings.Trim(p, "\"")
	p = strings.TrimSuffix(p, "\\n")

	if utf8.RuneCountInString(p) >= previewLength {
		p = strings.TrimSuffix(p[:previewLength]+"...", "\\n")
	}

	return html.EscapeString(p)
}

func CreateText(texts *mongo.Collection, conf *common.Configuration, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		identifier := ctx.Query("identifier")
		content := common.TrimPrefixesRecursive(strings.Trim(ctx.Query("content"), " $/^\\"), "\\n", "\\r")

		// we are assuming the random string is unused, TODO: add checking of random identifier,
		// potentially split no provided/provided identifier into it's own if/else block since they'll be quite similar.
		if identifier == "" {
			identifier = common.RandomString(8)
		}

		newText := &common.Text{
			Identifier:   identifier,
			Content:      content,
			Preview:      makePreview(content),
			CreationDate: time.Now().Format(time.RFC3339),
			DeletionKey:  common.RandomString(32), // unused, todo: redo delete API to check deletion key (also api key)
			Views:        0,
		}

		if err := texts.FindOne(context.TODO(), bson.M{"identifier": newText.Identifier}).Err(); err == nil {
			ctx.JSON(http.StatusConflict, &common.ErrorResponse{Message: "Identifier already exists, must be unique.", Code: http.StatusConflict})
			return
		}

		_, err := texts.InsertOne(todo, newText)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &common.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
			return
		}

		ctx.JSON(http.StatusCreated, &common.TextCreated{
			LocationURL: conf.Webserver.BaseURL + "t/" + newText.Identifier,
			Content:     newText.Content,
			DeletionURL: "NOT IMPLEMENTED",
		})
	}
}
