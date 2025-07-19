package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment struct contoh
type Comment struct {
	User    string `json:"user" bson:"user"`
	Message string `json:"message" bson:"message"`
	Date    time.Time `json:"date" bson:"date"`
}

// Image struct contoh
// type Image struct {
// 	URL     string `json:"url" bson:"url"`
// 	Caption string `json:"caption" bson:"caption"`
// }

// Blog adalah model utama
type Blog struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title          string             `json:"title" bson:"title"`
	Desc           string             `json:"desc" bson:"desc"`
	Sub_Desc        string             `json:"sub_desc" bson:"sub_desc"`
	Slug           string             `json:"slug" bson:"slug"`
	Content        string             `json:"content" bson:"content"`
	Comment        []Comment          `json:"comment" bson:"comment"`
	View           int             		`json:"view" bson:"view"`
	Status         string             `json:"status" bson:"status"`
	Image_bg string             		`json:"image_bg" bson:"image_bg"`
	Images          []string            `json:"images" bson:"images"`
	Category       string             `json:"category" bson:"category"`
	Tags           []string           `json:"tags" bson:"tags"`
	Author         string           `json:"author" bson:"author"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

// BeforeCreate dijalankan sebelum membuat dokumen baru
func (b *Blog) BeforeCreate() {
	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now
}

// BeforeUpdate dijalankan sebelum memperbarui dokumen
func (b *Blog) BeforeUpdate() {
	b.UpdatedAt = time.Now().UTC()
}


// Category
// IT, Since, 

// Nasional / Internasional

// Politik

// Ekonomi / Bisnis

// Technology

// Science

// Entertainment

// Kesehatan

// Pendidikan


// Gaya Hidup / Lifestyle

// Hiburan

// Olahraga

// Travel / Wisata

// Otomotif

// Kuliner / Makanan

