package main

import (
	"context"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetKey(key string) string {
	return os.Getenv(key)
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error read .env")
		os.Exit(1)
	}
}

var expired int = 300

func main() {
	LoadEnv()

	endpoint := GetKey("ENDPOINT")
	accessKeyID := GetKey("ACCESS_KEY_ID")
	accessSecretKey := GetKey("ACCESS_SECRET_KEY")
	bucket := GetKey("BUCKET")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, accessSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello!")
	})

	r.POST("/ups", func(c *gin.Context) {
		f, err := c.FormFile("file")
		if err != nil {
			panic(err)
		}

		filename := f.Filename
		size := f.Size
		r, err := f.Open()
		if err != nil {
			panic(err)
		}

		defer r.Close()

		_, err = minioClient.PutObject(c.Request.Context(), bucket, filename, r, size, minio.PutObjectOptions{
			ContentType:  mime.TypeByExtension(filepath.Ext(filename)),
			CacheControl: "max-age=31556952",
		})
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})

		reqParams := make(url.Values)
		reqParams.Set("content-disposition", "inline")
		presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucket, filename, time.Duration(expired)*time.Second, reqParams)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(presignedURL)
	})

	r.Run(":8080")
}

// some notes 1 https://github.com/gin-gonic/examples/blob/master/upload-file/single/main.go
// some notes 2 https://github.com/minio/minio-go/blob/master/examples/s3/fputobject.go
// some notes 3 https://github.com/siva-chegondi/filehandler
