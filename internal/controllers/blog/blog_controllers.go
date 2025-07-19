package blogController

import (
	"net/http"
	"strconv"

	// "log"
	"bytes"
    // "context"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "os"
    // "path/filepath"
    "time"


	// "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/Martinpasaribu/Golang-V1/internal/config"
    "github.com/Martinpasaribu/Golang-V1/internal/utils"
	"github.com/Martinpasaribu/Golang-V1/internal/models"
	blogService "github.com/Martinpasaribu/Golang-V1/internal/services/blog"
	"github.com/gin-gonic/gin"
)

type BlogController struct {
	service blogService.BlogService
}

func NewBlogController(service blogService.BlogService) *BlogController {
	return &BlogController{service: service}
}


func uploadToImageKit(file *multipart.FileHeader, folder string) (string, error) {
    // Dapatkan keys dari config
    _, privateKey, _ := config.GetImageKitKeys()

    // Buka file
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    // Siapkan form data
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    
    // Tambahkan file ke form
    part, err := writer.CreateFormFile("file", file.Filename)
    if err != nil {
        return "", err
    }
    _, err = io.Copy(part, src)
    if err != nil {
        return "", err
    }
    
    // Tambahkan parameter lainnya
    _ = writer.WriteField("fileName", fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename))
    _ = writer.WriteField("useUniqueFileName", "true")
    _ = writer.WriteField("folder", folder)
    
    // Tutup writer
    err = writer.Close()
    if err != nil {
        return "", err
    }

    // Buat HTTP request
    req, err := http.NewRequest("POST", "https://upload.imagekit.io/api/v1/files/upload", body)
    if err != nil {
        return "", err
    }
    
    // Set headers
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    // Basic auth menggunakan private key
    auth := base64.StdEncoding.EncodeToString([]byte(privateKey + ":"))
    req.Header.Set("Authorization", "Basic "+auth)

    // Kirim request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Baca response
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // Cek status code
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("imagekit upload failed: %d - %s", resp.StatusCode, string(respBody))
    }

    // Parse JSON response
    var result struct {
        URL string `json:"url"`
    }
    if err := json.Unmarshal(respBody, &result); err != nil {
        return "", err
    }

    return result.URL, nil
}

func (c *BlogController) UploadBlogImages(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Form tidak valid", "error": err.Error()})
		return
	}

	imageBg := form.File["image_bg"]
	images := form.File["images"]

	var imageBgURL string
	if len(imageBg) > 0 {
		url, err := uploadToImageKit(imageBg[0], "blogs/bg")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal upload image_bg", "error": err.Error()})
			return
		}
		imageBgURL = url
	}

	var imageURLs []string
	for _, file := range images {
		url, err := uploadToImageKit(file, "blogs/images")
		if err != nil {
			fmt.Printf("❌ Gagal upload: %v\n", err)
			continue
		}
		imageURLs = append(imageURLs, url)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"image_bg_url": imageBgURL,
		"images_url":   imageURLs,
	})
}


func (c *BlogController) CreateBlog(ctx *gin.Context) {
	var blog models.Blog


	// body, _ := io.ReadAll(ctx.Request.Body)
	// fmt.Println("🔥 RAW JSON BODY:", string(body))
	// ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // agar bisa dibaca ulang oleh ShouldBindJSON
	
	if os.Getenv("APP_DEBUG") == "true" {
		body, _ := io.ReadAll(ctx.Request.Body)
		fmt.Println("🔥 RAW JSON BODY:", string(body))
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	// Bind JSON request ke struct
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		fmt.Println("❌ Gagal bind JSON:", err.Error()) // ✅ Log ke terminal

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"id":      nil,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}


	// Set nilai default/otomatis
	// blog.View = 1 // ⬅️ Default jumlah view

	// Panggil method service untuk menyimpan blog
	result, err := c.service.CreateBlog(&blog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"id":      nil,
			"data":    nil,
			"message": "Failed to create blog",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"id":      result.ID.Hex(),
		"data":    result,
		"message": "Blog successfully created",
	})
}


// func (c *BlogController) GetBlogByID(ctx *gin.Context) {
// 	id := ctx.Param("id")
	
// 	blog, err := c.service.GetBlogByID(id)
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
// 		return
// 	}
	
// 	ctx.JSON(http.StatusOK, blog)
// }

// func (c *BlogController) UpdateBlog(ctx *gin.Context) {
// 	id := ctx.Param("id")
	
// 	var blog models.Blog
// 	if err := ctx.ShouldBindJSON(&blog); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
	
// 	if err := c.service.UpdateBlog(id, &blog); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
// 		return
// 	}
	
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully"})
// }

// func (c *BlogController) DeleteBlog(ctx *gin.Context) {
// 	id := ctx.Param("id")
	
// 	if err := c.service.DeleteBlog(id); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog"})
// 		return
// 	}
	
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
// }


func (c *BlogController) FindBlogBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	id := ctx.Query("id")

	var blog *models.Blog
	var err error

	blog, err = c.service.FindBlogBySlug(slug);

	if err != nil && id != "" {
		blog, err = c.service.FindBlogByID(id);
	}

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Blog not found"});
		return
	}

	ip := ctx.ClientIP() // ambil IP client

	// ✅ Cek dan tambahkan view hanya jika IP belum pernah akses blog ini
	if !utils.HasViewedBlog(blog.ID.Hex(), ip) {
		err = c.service.IncrementView(blog.ID.Hex())
		if err != nil {
			fmt.Println("❌ Gagal menambah view:", err.Error())
		}
		utils.SetViewedBlog(blog.ID.Hex(), ip, 24*time.Hour) // simpan IP di Redis selama 24 jam
	}

		// ✅ Ambil 3 blog lain dengan kategori yang sama
	var relatedBlogs []models.Blog
	relatedBlogs, err = c.service.FindRelatedBlogsByCategory(blog.Category, blog.ID.Hex(), 3)
	if err != nil {
		fmt.Println("❌ Gagal mengambil blog terkait:", err.Error())
		relatedBlogs = []models.Blog{} // fallback kosong
	}


	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success get blog",
		"data":    blog,
		"relatedBlogs": relatedBlogs,
	})
}


func (c *BlogController) GetAllBlogs(ctx *gin.Context) {
	// Ambil query parameter (page dan limit)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Fallback nilai default jika tidak valid
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Panggil service
	blogs, total, err := c.service.GetAllBlogs(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to fetch blogs",
			"data":    nil,
		})
		return
	}

	// Respon sukses
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success get blogs",
		"data":    blogs,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// Get Category
func (c *BlogController) GetCategory(ctx *gin.Context) {
	category := ctx.Query("category")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "6") // default 6 item

	page, _ := strconv.ParseInt(pageStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	latest, nextThree, fullData, err := c.service.GetCategory(category, page, limit)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil artikel"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Success Get Category",
		"latest":    latest,
		"nextThree": nextThree,
		"fullData":  fullData,
	})
}


func (c *BlogController) GetCategoryNavbar(ctx *gin.Context) {
	category := ctx.Query("category")

	main, main_sub, err := c.service.GetCategoryNavbar(category)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil GetCategory Navbar"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Success Get GetCategoryNavbar",
		"main":    main,
		"main_sub": main_sub,
	})
}

// ALl Element Articles

func (c *BlogController) GetArticles01(ctx *gin.Context) {
	
	latest, nextThree, err := c.service.GetArticles01()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil artikel"})
		return
	}

	// Respon sukses
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success GetArticles01",
		"data": []any{},
		"latest":    latest,
		"nextThree":   nextThree,
	})
}


func (c *BlogController) GetArticlesList(ctx *gin.Context) {

	nextThree, popularThree, err := c.service.GetArticlesList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil artikel"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        	http.StatusOK,
		"message":       	"Berhasil mengambil data artikel",
		"nextThree": 		nextThree,
		"popularThree":     popularThree,
	})
}
