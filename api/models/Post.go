package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Post Post table
type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"text;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//Prepare prep Post
func (p *Post) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

//Validate Contents are they blank or not
func (p *Post) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)
	if p.Title == "" {
		err = errors.New("Required Title")
		errorMessages["Required_title"] = err.Error()
	}

	if p.Content == "" {
		err = errors.New("Required Content")
		errorMessages["Required_content"] = err.Error()
	}
	if p.AuthorID < 1 {
		err = errors.New("Required Author")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}

//SavePost save written Post
func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

//FindAllPosts Find all current posts
func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Limit(100).Order("Created_at desc").Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

//FindPostByID find post by id
func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

//UpdatedAPost update a post
func (p *Post) UpdatedAPost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, err
}

//DeleteAPost delete a post
func (p *Post) DeleteAPost(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Post{}).Where("id =?", p.ID).Take(&Post{}).Delete(&Post{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//FindUserPosts find user's posts
func (p *Post) FindUserPosts(db *gorm.DB, uid uint32) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Where("author_id=?", uid).Limit(100).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

//DeleteUserPosts delete all user's posts
func (p *Post) DeleteUserPosts(db *gorm.DB, uid uint32) (int64, error) {
	posts := []Post{}
	db = db.Debug().Model(&Post{}).Where("author_id = ?", uid).Find(&posts).Delete(&posts)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
