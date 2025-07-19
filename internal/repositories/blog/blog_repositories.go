package blogRepository

import (
	"context"
	"errors"

	"github.com/Martinpasaribu/Golang-V1/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BlogRepository interface mendefinisikan operasi database untuk blog
type BlogRepository interface {
	CreateBlog(blog *models.Blog) (*models.Blog, error)
	UpdateBlog(id primitive.ObjectID, blog *models.Blog) error
	DeleteBlog(id primitive.ObjectID) error
	
	IncrementView(id primitive.ObjectID) error 
	FindBlogByID(id string) (*models.Blog, error) // pastikan juga ini ada
	FindBlogBySlug(slug string) (*models.Blog, error)
	GetAllBlogs(page int, limit int) ([]models.Blog, int64, error) // 🟢 Tambahkan ini
	GetCategory(category string, page int64, limit int64) (models.Blog, []models.Blog, []models.Blog, error)
	GetCategoryNavbar(category string) (models.Blog,[]models.Blog, error)

	GetArticles01() (models.Blog, []models.Blog, error) // ✅ BENAR
	GetArticlesList() ([]models.Blog, []models.Blog, error)
	FindRelatedBlogsByCategory(category string, excludeID string, limit int64) ([]models.Blog, error)

}

type blogRepository struct {
	collection *mongo.Collection
}

// Perbaiki: Tambahkan parameter *mongo.Database
func NewBlogRepository(db *mongo.Database) BlogRepository {
	return &blogRepository{
		collection: db.Collection("blogs"), // Gunakan database dari parameter
	}
}

// CreateBlog menyimpan blog baru ke database
func (r *blogRepository) CreateBlog(blog *models.Blog) (*models.Blog, error) {
	// Set timestamps
	blog.BeforeCreate()
	
	res, err := r.collection.InsertOne(context.Background(), blog)
	if err != nil {
		return nil, err
	}
	
	blog.ID = res.InsertedID.(primitive.ObjectID)
	return blog, nil
}

// FindBlogByID mencari blog berdasarkan ID
// func (r *blogRepository) FindBlogByID(id primitive.ObjectID) (*models.Blog, error) {
// 	var blog models.Blog
// 	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&blog)
	
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, errors.New("blog not found")
// 		}
// 		return nil, err
// 	}
	
// 	return &blog, nil
// }

// UpdateBlog memperbarui blog yang ada
func (r *blogRepository) UpdateBlog(id primitive.ObjectID, blog *models.Blog) error {
	// Update updated_at
	blog.BeforeUpdate()
	
	update := bson.M{
		"$set": bson.M{
			"title":      blog.Title,
			"slug":       blog.Slug,
			"content":    blog.Content,
			"updated_at": blog.UpdatedAt,
		},
	}
	
	_, err := r.collection.UpdateByID(context.Background(), id, update)
	return err
}

// DeleteBlog menghapus blog berdasarkan ID
func (r *blogRepository) DeleteBlog(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

// GetAllBlogs mengambil semua blog dengan paginasi



func (r *blogRepository) GetAllBlogs(page int, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	skip := (page - 1) * limit

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &blogs); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}


func (r *blogRepository) FindBlogBySlug(slug string) (*models.Blog, error) {
	var blog models.Blog
	err := r.collection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&blog)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("blog not found")
		}
		return nil, err
	}
	
	return &blog, nil
}


func (r *blogRepository) FindBlogByID(id string) (*models.Blog, error) {
	// Konversi string id dari query ke ObjectID MongoDB
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	var blog models.Blog

	// Query ke MongoDB berdasarkan _id (ObjectID)
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("blog not found")
		}
		return nil, err
	}

	return &blog, nil
}

func (r *blogRepository) FindRelatedBlogsByCategory(category string, excludeID string, limit int64) ([]models.Blog, error) {
	
	objID, err := primitive.ObjectIDFromHex(excludeID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"category": category,
		"_id": bson.M{
			"$ne": objID, // Exclude blog utama
		},
	}

	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.M{"createdAt": -1}) // Ambil yang terbaru

	cursor, err := r.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var blogs []models.Blog
	if err := cursor.All(context.TODO(), &blogs); err != nil {
		return nil, err
	}

	return blogs, nil
}



func (r *blogRepository) IncrementView(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"view": 1}}

	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}


// Get Category 
func (r *blogRepository) GetCategory(category string, page int64, limit int64) (models.Blog, []models.Blog, []models.Blog, error) {
	var latest models.Blog
	var nextThree []models.Blog
	var fullData []models.Blog

	projection := bson.M{
		"content": 0,
		"desc":    0,
	}

	filter := bson.M{
		"category": category,
	}

	// Ambil 1 data terbaru
	err := r.collection.FindOne(
		context.TODO(),
		filter,
		options.FindOne().
			SetSort(bson.M{"created_at": -1}).
			SetProjection(projection),
	).Decode(&latest)

	if err != nil {
		return models.Blog{}, nil, nil, err
	}

	// Ambil 3 data berikutnya
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(1).
		SetLimit(3).
		SetProjection(projection)


	cursor, err := r.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return models.Blog{}, nil, nil, err
	}

	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &nextThree); err != nil {
		return models.Blog{}, nil, nil, err
	}

	// Full data (pagination)
	fullOpts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetProjection(projection)

	fullCursor, err := r.collection.Find(context.TODO(), filter, fullOpts)
	if err != nil {
		return models.Blog{}, nil, nil, err
	}
	defer fullCursor.Close(context.TODO())
	if err := fullCursor.All(context.TODO(), &fullData); err != nil {
		return models.Blog{}, nil, nil, err
	}

	return latest, nextThree, fullData, nil
}

// Get Category Navbar
func (r *blogRepository) GetCategoryNavbar(category string) (models.Blog, []models.Blog, error){

	var main models.Blog
	var main_sub []models.Blog

	projection := bson.M{
		"content": 0,
		"desc":    0,
		"images":  0,
		"view":    0,
		"comment": 0,
	}

	filter := bson.M{ "category": category,}

	// Ambil 1 data terbaru
	err := r.collection.FindOne(
		context.TODO(),
		filter,
		options.FindOne().
			SetSort(bson.M{"created_at": -1}).
			SetProjection(projection),
	).Decode(&main)

	if err != nil {
		return models.Blog{}, nil, err
	}

	// Ambil 6 data berikutnya
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(1).
		SetLimit(6).
		SetProjection(projection);

	cursor, err := r.collection.Find(context.TODO(), filter, opts);

	if err != nil {
		return models.Blog{}, nil, err
	}

	defer cursor.Close(context.TODO());
	if err := cursor.All(context.TODO(), &main_sub); err != nil {
		return models.Blog{}, nil, err
	}

	return main, main_sub, nil;

}


// ALL Element Articles 
func (r *blogRepository) GetArticles01() (models.Blog, []models.Blog, error) {
	var latest models.Blog
	var nextThree []models.Blog

	// ❌ Exclude content dan author, sisanya MongoDB akan tampilkan otomatis
	projection := bson.M{
		"content": 0,
		"desc":  0,
	}

	// Ambil 1 data terbaru
	err := r.collection.FindOne(
		context.TODO(),
		bson.M{},
		options.FindOne().
			SetSort(bson.M{"created_at": -1}).
			SetProjection(projection),
	).Decode(&latest)

	if err != nil {
		return models.Blog{}, nil, err
	}

	// Ambil 3 data setelahnya
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(1).
		SetLimit(3).
		SetProjection(projection)

	cursor, err := r.collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return models.Blog{}, nil, err
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &nextThree); err != nil {
		return models.Blog{}, nil, err
	}

	return latest, nextThree, nil
}


func (r *blogRepository) GetArticlesList() ([]models.Blog, []models.Blog, error) {
	var nextThree []models.Blog
	var popularThree []models.Blog

	// Projection: sembunyikan content dan desc
	projection := bson.M{
		"content": 0,
		"desc":    0,
	}

	// 👉 Ambil 3 artikel setelah artikel pertama (artikel ke-2, 3, 4)
	optsNext := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(1).
		SetLimit(3).
		SetProjection(projection)

	cursorNext, err := r.collection.Find(context.TODO(), bson.M{}, optsNext)
	if err != nil {
		return nil, nil, err
	}
	defer cursorNext.Close(context.TODO())

	if err := cursorNext.All(context.TODO(), &nextThree); err != nil {
		return nil, nil, err
	}

	// 👉 Ambil 3 artikel terpopuler berdasarkan views
	optsPopular := options.Find().
		SetSort(bson.M{"views": -1}).
		SetLimit(3).
		SetProjection(projection)

	cursorPopular, err := r.collection.Find(context.TODO(), bson.M{}, optsPopular)
	if err != nil {
		return nil, nil, err
	}
	defer cursorPopular.Close(context.TODO())

	if err := cursorPopular.All(context.TODO(), &popularThree); err != nil {
		return nil, nil, err
	}

	return nextThree, popularThree, nil
}
