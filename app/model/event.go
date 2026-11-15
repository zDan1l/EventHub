package model

import "time"

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Quota       int       `json:"quota"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Quota       int       `json:"quota"`
	IsActive    bool      `json:"is_active"`
}

type PatchEventRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Quota       *int       `json:"quota"`
	IsActive    *bool      `json:"is_active"`
}


