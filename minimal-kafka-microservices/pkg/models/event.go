package models

import "time"

type Event struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Data    string    `json:"data"`
	Created time.Time `json:"created"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
