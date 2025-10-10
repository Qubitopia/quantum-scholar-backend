package handlers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/image/draw"

	"github.com/Qubitopia/quantum-scholar-backend/database"
	"github.com/Qubitopia/quantum-scholar-backend/models"
	"github.com/gin-gonic/gin"
)

type DeleteImageRequest struct {
	TestID   uint32 `json:"test_id" binding:"required"`
	Filename string `json:"filename" binding:"required"`
}

// UploadImage handles authenticated image upload to object storage
func UploadImage(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userRaw.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user type"})
		return
	}

	test_id := c.Param("test_id")
	if test_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing test_id parameter"})
		return
	}

	// Check if test exists and belongs to user
	var test models.Test
	if err := database.DB.Where("test_id = ? AND examiner_id = ?", test_id, user.ID).First(&test).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
		return
	}

	// Deduct QS Coins (1 per image)
	if test.QSCoins < 1 {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient QS Coins to upload image"})
		return
	}
	test.QSCoins -= 1 // Deduct 1 QS Coins per image upload

	// Expect raw file body (no multipart). Enforce 10MB max size.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	// Determine content type from header (allow jpeg or png)
	declaredCT := strings.ToLower(c.GetHeader("Content-Type"))
	if declaredCT == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Content-Type header"})
		return
	}
	if declaredCT != "image/jpeg" && declaredCT != "image/pjpeg" && declaredCT != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only image/jpeg or image/png are supported"})
		return
	}

	// Format new filename: email-datetime-filename.filetype
	email := user.Email
	now := time.Now().UTC().Format("20060102T150405Z")
	// No original filename available in raw body; use a fixed base
	safeBase := "image"
	newName := fmt.Sprintf("%s-%s-%s.%s", email, now, safeBase, "jpg")

	// Read file into memory (<=10MB due to MaxBytesReader limit)
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		// Differentiate size errors if needed
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File too large (max 10MB)"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file: " + err.Error()})
		return
	}

	// Decode image (supports jpeg & png due to registered decoders)
	img, format, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image: " + err.Error()})
		return
	}
	if format != "jpeg" && format != "png" { // extra safety
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported image format"})
		return
	}

	// Resize if width exceeds 640px while preserving aspect ratio
	// max width 640px
	b := img.Bounds()
	origW := b.Dx()
	origH := b.Dy()
	if origW > 640 {
		newW := 640
		newH := int(float64(origH) * (float64(newW) / float64(origW)))
		// Create RGBA canvas and scale
		resized := image.NewRGBA(image.Rect(0, 0, newW, newH))
		draw.ApproxBiLinear.Scale(resized, resized.Bounds(), img, b, draw.Over, nil)
		img = resized
	}

	// Re-encode/compress as JPEG
	var out bytes.Buffer
	if err := jpeg.Encode(&out, img, &jpeg.Options{Quality: 80}); err != nil { // quality can be tuned
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compress image: " + err.Error()})
		return
	}

	contentType := "image/jpeg"
	if err = database.UploadObject(newName, contentType, out.Bytes()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Save filename in test's Images array
	test.Images = append(test.Images, newName)

	// Update test record in database
	if err := database.DB.Save(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test record: " + err.Error()})
		return
	}

	// Generate presigned URL valid for 15 minutes
	url, err := database.GetPresignedURL(newName, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"filename": newName, "url": url})
}

func GetImageURL(c *gin.Context) {
	filename := c.Param("imagename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing filename parameter"})
		return
	}
	url, err := database.GetPresignedURL(filename, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func DeleteImage(c *gin.Context) {
	userRaw, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userRaw, ok := userRaw.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}
	user := userRaw.(models.User)

	// Check if the image exist isn the test and belongs to the user
	var req DeleteImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var test models.Test
	if err := database.DB.Where("test_id = ? AND examiner_id = ?", req.TestID, user.ID).First(&test).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test not found"})
		return
	}

	// Check if the image exists in the test's Images array
	imageIndex := -1
	for i, img := range test.Images {
		if img == req.Filename {
			imageIndex = i
			break
		}
	}
	if imageIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found in test"})
		return
	}
	test.Images = append(test.Images[:imageIndex], test.Images[imageIndex+1:]...)

	// delete the image from object storage
	if err := database.DeleteObject(req.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image from storage: " + err.Error()})
		return
	}

	// update the test record in the database
	if err := database.DB.Save(&test).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test record: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}
