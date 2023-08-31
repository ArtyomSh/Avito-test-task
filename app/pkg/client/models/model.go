package models

import "gorm.io/gorm"

type User struct {
	ID       uint       `json:"ID"`
	Segments []*Segment `json:"segments,omitempty" gorm:"many2many:user_segments"`
}

type Segment struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	Name    string  `json:"name"`
	Users   []*User `json:"users,omitempty" gorm:"many2many:user_segments"`
	Deleted gorm.DeletedAt
}
