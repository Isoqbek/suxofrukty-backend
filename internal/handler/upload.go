package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler { return &UploadHandler{} }

func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not provided"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	if privateKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage not configured (IMAGEKIT_PRIVATE_KEY missing)"})
		return
	}

	fileName := fmt.Sprintf("%d-%s", time.Now().UnixMilli(), header.Filename)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fw, _ := w.CreateFormField("fileName")
	fw.Write([]byte(fileName))
	fw2, _ := w.CreateFormField("folder")
	fw2.Write([]byte("/products"))
	fw3, _ := w.CreateFormFile("file", fileName)
	fw3.Write(data)
	w.Close()

	req, _ := http.NewRequest("POST", "https://upload.imagekit.io/api/v1/files/upload", &body)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(privateKey+":")))
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload request failed"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body2, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "imagekit error: " + string(body2)})
		return
	}

	var result struct {
		URL string `json:"url"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	c.JSON(http.StatusOK, gin.H{"url": result.URL})
}
