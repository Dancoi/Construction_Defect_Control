package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DefectID  uint      `json:"defect_id"`
	Defect    Defect    `gorm:"foreignKey:DefectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"defect,omitempty"`
	AuthorID  *uint     `json:"author_id"`
	Author    *User     `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"author,omitempty"`
	Body      string    `gorm:"type:text" json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
