package models

import "time"

type Attachment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DefectID    uint      `json:"defect_id"`
	Defect      Defect    `gorm:"foreignKey:DefectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"defect,omitempty"`
	UploaderID  uint      `json:"uploader_id"`
	Uploader    User      `gorm:"foreignKey:UploaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"uploader,omitempty"`
	Path        string    `gorm:"size:1024" json:"path"`
	Filename    string    `gorm:"size:512" json:"filename"`
	ContentType string    `gorm:"size:255" json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}
