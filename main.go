package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

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
		c.String(http.StatusOK, "Isoku mek Hello World")
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
			ContentType: mime.TypeByExtension(filepath.Ext(filename)),
		})
		if err != nil {
			panic(err)
		}

		// file, _ := c.FormFile("file")
		// filename := filepath.Base(file.Filename)

		// if err := c.SaveUploadedFile(file, filename); err != nil {
		// 	// Handling FPutObject :(
		// 	if _, err := minioClient.FPutObject(context.Background(), bucket, filename, file, minio.PutObjectOptions{
		// 		ContentType: "multipart/form-data",
		// 	}); err != nil {
		// 		log.Fatalln(err)
		// 	}
		// }
		c.JSON(http.StatusOK, gin.H{
			"message":   "ok",
			"x-dotwell": "384kbps",
		})
	})

	// List bucket test
	// buckets, err := minioClient.ListBuckets(context.Background())
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// for _, bucket := range buckets {
	// 	log.Println(bucket)
	// }

	// debs
	// GetEp := GetKey("ENDPOINT")
	// fmt.Println("my ep is", GetEp)

	r.Run(":8080")
}

// some notes 1 https://github.com/gin-gonic/examples/blob/master/upload-file/single/main.go
// some notes 2 https://github.com/minio/minio-go/blob/master/examples/s3/fputobject.go
// some notes 3 https://github.com/siva-chegondi/filehandler
