package api

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/awildboop/gosharex/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateImage(images *mongo.Collection, conf *common.Configuration, todo context.Context) func(*gin.Context) {
	return func(ctx *gin.Context) {
		upload_date := time.Now().Format(time.RFC3339)
		identifier := ctx.Query("identifier")
		upload, err := ctx.FormFile("upload")
		if common.HandleErr(err, ctx) {
			return
		}

		// we are assuming the random string is unused, TODO: add checking of random identifier,
		// potentially split no provided/provided identifier into it's own if/else block since they'll be quite similar.
		if identifier == "" {
			identifier = common.RandomString(8)

		}

		newImage := &common.Image{
			Identifier:   identifier,
			DeletionKey:  common.RandomString(32),
			CreationDate: upload_date,
			FileSize:     upload.Size,
			Views:        0,
		}

		if err := images.FindOne(context.TODO(), bson.M{"identifier": newImage.Identifier}).Err(); err == nil {
			ctx.JSON(http.StatusConflict, &common.ErrorResponse{Message: "Identifier already exists, must be unique.", Code: http.StatusConflict})
			return
		}

		fh, err := upload.Open()
		if common.HandleErr(err, ctx) {
			return
		}
		defer fh.Close()

		buffer, err := io.ReadAll(bufio.NewReader(fh))
		if common.HandleErr(err, ctx) {
			return
		}

		// use DetectContentType to ensure file  type isn't being faked in http header
		mime_type := http.DetectContentType(buffer)

		if !common.HasAnySuffix(mime_type, "png", "jpeg") {
			ctx.JSON(http.StatusBadRequest, &common.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid file type.",
			})
			return
		}

		newImage.FileLocation = fmt.Sprintf("%s-%s.%s", identifier, upload_date, strings.Split(mime_type, "/")[1])

		file, err := os.Create("./uploads/" + newImage.FileLocation)
		if common.HandleErr(err, ctx) {
			return
		}

		_, err = file.Write(buffer)
		if common.HandleErr(err, ctx) {
			return
		}

		if _, err = images.InsertOne(todo, newImage); err != nil {
			ctx.JSON(http.StatusInternalServerError, &common.ErrorResponse{Message: err.Error(), Code: http.StatusInternalServerError})
			return
		}

		var return_url = conf.Webserver.BaseURL + "%s/" + identifier
		if conf.Features.Extra.UseRawImageURL {
			return_url = fmt.Sprintf(return_url, "ri")
		} else {
			return_url = fmt.Sprintf(return_url, "i")
		}

		ctx.JSON(http.StatusCreated, common.ImageCreated{
			LocationURL: return_url,
			FileSize:    newImage.FileSize,
			DeletionURL: "NOT IMPLEMENTED",
		})
	}
}
