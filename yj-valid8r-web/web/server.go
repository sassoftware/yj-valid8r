package web

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"

	internal "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-common"
	validator "github.com/Huzaib-Sayyed_sasinst/yj-valid8r/yj-valid8r-lib"
)

//go:embed templates/*
var tmpl embed.FS

func StartServer() {
	log.Println("Application started")
	port := "7070"
	router := gin.New()
	router.Use(gin.Recovery())
	// router.SetTrustedProxies([]string{"127.0.0.1", "localhost"})

	templates, err := template.ParseFS(tmpl, "templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	router.SetHTMLTemplate(templates)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "validator.html", nil)
	})
	router.POST("/api/validate", handleValidate)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}

func handleValidate(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	var req internal.ValidationRequest

	switch contentType {
	case "application/json":
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON body: %v", err)})
			return
		}
	case "application/x-yaml":
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to read request body: %v", err)})
			return
		}
		if err := yaml.Unmarshal(body, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid YAML body: %v", err)})
			return
		}
	default:
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported Content-Type. Content-Type Supported are \"application/json\", \"application/x-yaml\""})
		return
	}

	dataStr := strings.TrimSpace(req.Data)
	if dataStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No data provided. Please supply valid YAML or JSON.",
		})
		return
	}

	dataBytes := []byte(req.Data)

	if validator.IsUnknownDataType(dataBytes) {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error": "Provided data is neither valid JSON nor YAML. Please check if your YAML/JSON is correct.",
		})
		return
	}

	checkTrailingWhitespace := true
	if req.CheckTrailingWhitespace != nil {
		checkTrailingWhitespace = *req.CheckTrailingWhitespace
	}

	results := internal.InitValidation(req.Schemas, dataBytes, checkTrailingWhitespace, req.RegexPatternRules, req.SearchPaths, req.Plugins)

	c.JSON(http.StatusOK, results)
}
