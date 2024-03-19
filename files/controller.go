package files

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"job_board/helpers"
)

type Error struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

var uploadFlyKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	uploadFlyKey = os.Getenv("UPLOAD_FLY")
	if uploadFlyKey == "" {
		panic("Error loading upload fly key")
	}
}

func uploadHandler(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	log.Println(file.Filename)

	newClient := NewAPIClient(uploadFlyKey)
	resp, err := newClient.Post(UploadRequest{
		File:             *file,
		Filename:         file.Filename,
		MaxFileSize:      "8MB",
		AllowedFileTypes: ".jpg,.jpeg,.png,.pdf",
	})
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.CreateResponse(ctx, helpers.Response{
			Message:    err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	//check for response
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var failed Error
		err = json.Unmarshal(body, &failed)
		if err != nil {
			helpers.CreateResponse(ctx, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}

		helpers.CreateResponse(ctx, helpers.Response{
			Message:    errors.New("received non-200 response" + failed.Message).Error(),
			StatusCode: http.StatusInternalServerError,
			Data:       nil,
		})
		return
	}

	// Upload the file to specific dst.
	// ctx.SaveUploadedFile(file, dst)
}
