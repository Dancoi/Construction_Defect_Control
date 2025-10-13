package models

import "time"

type Defect struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	ProjectID uint `json:"project_id"`
	// Project is the relation to project; add constraint to create FK
	Project     Project    `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"project,omitempty"`
	Title       string     `gorm:"size:255" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Severity    string     `gorm:"size:50" json:"severity"`
	Status      string     `gorm:"size:50" json:"status"`
	AssigneeID  *uint      `json:"assignee_id,omitempty"`
	Assignee    *User      `gorm:"foreignKey:AssigneeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"assignee,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Priority    string     `gorm:"size:50" json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
