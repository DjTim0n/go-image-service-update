package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nfnt/resize"
)

func createResponse(filename string, folder *string) gin.H {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	BASEURL := os.Getenv("BASEURL") + "images/"

	if folder != nil {
		BASEURL = BASEURL + *folder + "/"
	}

	urls := map[string]string{
		"original": BASEURL + filename,
		"320px":    BASEURL + "320px_" + filename,
		"480px":    BASEURL + "480px_" + filename,
	}

	return gin.H{
		"file":    filename,
		"message": "Image uploaded successfully",
		"urls":    urls,
	}
}

func uploadImageToFolder(c *gin.Context, folder string) {
	file, _ := c.FormFile("image")
	c.SaveUploadedFile(file, "images/"+folder+"/"+file.Filename)
	resizeImage(320, file.Filename, &folder)
	resizeImage(480, file.Filename, &folder)

	response := createResponse(file.Filename, &folder)
	c.IndentedJSON(http.StatusOK, response)

}

func uploadImage(c *gin.Context) {
	file, _ := c.FormFile("image")
	c.SaveUploadedFile(file, "images/"+file.Filename)

	resizeImage(320, file.Filename, nil)
	resizeImage(480, file.Filename, nil)

	response := createResponse(file.Filename, nil)
	c.IndentedJSON(http.StatusOK, response)

}

func resizeImage(width int, filename string, folder *string) string {
	height := 0
	var filepath string

	if folder != nil {
		filepath = "images/" + *folder + "/" + filename
	} else {
		filepath = "images/" + filename
	}

	file, err := os.Open(filepath)
	if err != nil {
		return "error opening file"
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "error decoding image"
	}

	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	var thumbnailPath string
	if folder != nil {
		thumbnailPath = "images/" + *folder + "/" + strconv.Itoa(width) + "px_" + filename
	} else {
		thumbnailPath = "images/" + strconv.Itoa(width) + "px_" + filename
	}

	outFile, err := os.Create(thumbnailPath)
	if err != nil {
		return "error creating thumbnail"
	}
	defer outFile.Close()
	jpeg.Encode(outFile, resizedImage, nil)

	return thumbnailPath
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Выводим в консоль порт, на котором будет работать сервер
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Разрешённые домены
		AllowMethods:     []string{"*"}, // Разрешённые методы
		AllowHeaders:     []string{"*"}, // Разрешённые заголовки
		ExposeHeaders:    []string{"*"}, // Заголовки для клиента
		AllowCredentials: true,          // Разрешить куки
	}))

	router.POST("/uploadImage", uploadImage)
	router.POST("/uploadImage/:folder", func(c *gin.Context) {
		folder := c.Param("folder")
		uploadImageToFolder(c, folder)
	})
	router.Static("/images", "./images")
	router.Run(port)
}
