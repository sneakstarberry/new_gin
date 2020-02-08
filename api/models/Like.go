package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//Like when we liked post maked that post
type Like struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	PostID    uint64    `gorm:"not null" json:"post_id"`
	CreatedAt time.Time `gorm:"default_CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default_CURRENT_TIMESTAMP" jsont:"updated_at"`
}

//SaveLike Save like post
func (l *Like) SaveLike(db *gorm.DB) (*Like, error) {
	// Check if the auth user has liked this post before:
	err := db.Debug().Model(&Like{}).Where("post_id =? AND user_id = ?", l.PostID, l.UserID).Take(&l).Error
	if err != nil {
		if err.Error() == "record not found" {
			// the user has not liked this post before, so lets save incomming like:
			err = db.Debug().Model(&Like{}).Create(&l).Error
			if err != nil {
				return &Like{}, err
			}
		}
	} else {
		err = errors.New("double like")
		return &Like{}, err
	}
	return l, nil
}

//GetLikesInfo get what we like post information
func (l *Like) GetLikesInfo(db *gorm.DB, pid uint64) (*[]Like, error) {
	likes := []Like{}
	err := db.Debug().Model(&Like{}).Where("post_id=?", pid).Find(&likes).Error
	if err != nil {
		return &[]Like{}, err
	}
	return &likes, err
}

//DeleteUserLikes When a post is deleted, we also delete the likes that the post had
func (l *Like) DeleteUserLikes(db *gorm.DB, uid uint32) (int64, error) {
	likes := []Like{}
	db = db.Debug().Model(&Like{}).Where("user_id =?", uid).Find(&likes).Delete(&likes)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffeted, nil
}

//DeletePostLikes When a post is deleted, we also delete the likes that the post had
func (l *Like) DeletePostLikes(db *gorm.DB, pid uint64) (int64, error) {
	likes := []Like{}
	db = db.Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Delete(&likes)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffeted, nil
}
