package blogService

import (
	"errors"

	
	"github.com/Martinpasaribu/Golang-V1/internal/models"
	blogRepository "github.com/Martinpasaribu/Golang-V1/internal/repositories/blog"
	"go.mongodb.org/mongo-driver/bson/primitive" // <-- TAMBAHKAN INI
)

type BlogService interface {
	CreateBlog(blog *models.Blog) (*models.Blog, error)
	// GetBlogBySlug(slug string) (*models.Blog, error)
	// UpdateBlog(id string, blog *models.Blog) error
	// DeleteBlog(id string) error
	
	
	FindBlogByID(id string) (*models.Blog, error)
	FindBlogBySlug(slug string) (*models.Blog, error)
	GetAllBlogs(page, limit int) ([]models.Blog, int64, error)
	IncrementView(id string) error 

	// ALl Element Articles
	GetArticles01() (models.Blog, []models.Blog, error) // ✅ BENAR
	GetArticlesList() ([]models.Blog, []models.Blog, error) // ✅ BENAR

	// Get Category
	GetCategory(category string, page int64, limit int64) (models.Blog, []models.Blog, []models.Blog, error);
	GetCategoryNavbar(category string) (models.Blog, []models.Blog, error)


	FindRelatedBlogsByCategory(category string, excludeID string, limit int) ([]models.Blog, error)


}

type blogService struct {
	repo blogRepository.BlogRepository
}

func NewBlogService(repo blogRepository.BlogRepository) BlogService {
	return &blogService{repo: repo}
}

func (s *blogService) CreateBlog(blog *models.Blog) (*models.Blog, error) {
	// Validasi tambahan bisa ditambahkan di sini
	return s.repo.CreateBlog(blog)
}

// func (s *blogService) GetBlogByID(id string) (*models.Blog, error) {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, errors.New("invalid ID format")
// 	}
// 	return s.repo.FindBlogByID(objID)
// }


// func (s *blogService) UpdateBlog(id string, blog *models.Blog) error {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return errors.New("invalid ID format")
// 	}
// 	return s.repo.UpdateBlog(objID, blog)
// }


// func (s *blogService) DeleteBlog(id string) error {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return errors.New("invalid ID format")
// 	}
// 	return s.repo.DeleteBlog(objID)
// }


func (s *blogService) IncrementView(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}
	return s.repo.IncrementView(objID)
}


func (s *blogService) GetAllBlogs(page int, limit int) ([]models.Blog, int64, error) {
	return s.repo.GetAllBlogs(page, limit)
}

func (s *blogService) FindBlogBySlug(slug string) (*models.Blog, error) {
	return s.repo.FindBlogBySlug(slug)
}

func (s *blogService) FindBlogByID(id string) (*models.Blog, error) {
	return s.repo.FindBlogByID(id) // kirim string saja
}

// Function Category

func (s *blogService) GetCategory(category string, page int64, limit int64) (models.Blog, []models.Blog, []models.Blog, error) {
	return s.repo.GetCategory(category, page, limit)
}

func (s *blogService) GetCategoryNavbar(category string) (models.Blog, []models.Blog, error) {
	return s.repo.GetCategoryNavbar(category);
}


// ALl Element Articles

func (s *blogService) GetArticles01() (models.Blog, []models.Blog, error) {
	return s.repo.GetArticles01()
}

func (s *blogService) GetArticlesList() ([]models.Blog, []models.Blog, error) {
	return s.repo.GetArticlesList()
}


func (s *blogService) FindRelatedBlogsByCategory(category string, excludeID string, limit int) ([]models.Blog, error) {
	return s.repo.FindRelatedBlogsByCategory(category, excludeID, int64(limit))
}



