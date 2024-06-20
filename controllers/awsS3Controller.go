package controllers

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
)

func AwsS3Controller(f multipart.FileHeader) (string, error) {
	//AKIAZQ3DSUWKAERD5VRK
	//Hhko9f3cFpHAq3ry+jT1IdGdtvCNhoDh2qbg70SW
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*aws.String("ap-south-1")))

	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}

	client := s3.NewFromConfig(cfg)

	file, err := f.Open()

	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}
	defer file.Close()

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{

		Bucket:      aws.String("gecko-thumbnail"),
		Key:         aws.String(f.Filename),
		Body:        file,
		ACL:         "public-read",
		ContentType: aws.String("image/png"),
		Expires:     aws.Time(time.Now().Add(time.Duration(time.Now().Year()) * 100)),
	})

	if err != nil {
		log.Printf("error uploader: %v", err)
		return "", err
	}
	return result.Location, nil
}

func UploadThumbnailFiles(c *gin.Context) {
	var body struct {
		image multipart.FileHeader `form:"image" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	file, err := c.FormFile("image")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	result, err := AwsS3Controller(*file)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to upload file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"url":     result,
	})

}
