package handler

import "time"

// RegisterRequest is used for swagger to describe register payload
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
	Name     string `json:"name" example:"John Doe"`
}

// LoginRequest is used for swagger to describe login payload
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret123"`
}

// UserResponse represents a user returned by the API
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"John Doe"`
	Role      string    `json:"role" example:"engineer"`
	CreatedAt time.Time `json:"created_at" example:"2025-10-12T12:00:00Z"`
}

// AuthResponse represents login/register response
type AuthResponse struct {
	Token string        `json:"token,omitempty" example:"eyJhbGci..."`
	User  *UserResponse `json:"user,omitempty"`
}

// ProjectResponse represents a project
type ProjectResponse struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"New Building"`
	Address   string    `json:"address" example:"123 Main St, City"`
	CreatedAt time.Time `json:"created_at" example:"2025-10-12T12:00:00Z"`
}

// DefectResponse represents a defect
type DefectResponse struct {
	ID          uint      `json:"id" example:"1"`
	ProjectID   uint      `json:"project_id" example:"1"`
	Title       string    `json:"title" example:"Cracked wall"`
	Description string    `json:"description" example:"Long vertical crack on east wall"`
	Severity    string    `json:"severity" example:"major"`
	Status      string    `json:"status" example:"open"`
	CreatedAt   time.Time `json:"created_at" example:"2025-10-12T12:00:00Z"`
}

// AttachmentResponse represents an attachment
type AttachmentResponse struct {
	ID          uint      `json:"id" example:"1"`
	DefectID    uint      `json:"defect_id" example:"1"`
	UploaderID  uint      `json:"uploader_id" example:"2"`
	Filename    string    `json:"filename" example:"photo.jpg"`
	ContentType string    `json:"content_type" example:"image/jpeg"`
	Size        int64     `json:"size" example:"23456"`
	URL         string    `json:"url" example:"/uploads/2025/10/12/uuid-photo.jpg"`
	CreatedAt   time.Time `json:"created_at" example:"2025-10-12T12:00:00Z"`
}

// CreateProjectRequest used in swagger for creating projects
type CreateProjectRequest struct {
	Name    string `json:"name" example:"New Building"`
	Address string `json:"address" example:"123 Main St"`
}

// CreateDefectRequest used in swagger for creating defects
type CreateDefectRequest struct {
	ProjectID   uint   `json:"project_id" example:"1"`
	Title       string `json:"title" example:"Leaking pipe"`
	Description string `json:"description" example:"Pipe leaking near ceiling"`
	Severity    string `json:"severity" example:"minor"`
}
