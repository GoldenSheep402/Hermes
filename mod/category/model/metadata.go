package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

// Metadata is the model for metadata
/*
It's used to store the metadata of a category.
Example:
For book category, the metadata can be the author, publisher, etc.
For movie category, the metadata can be the director, producer, etc.
*/
type Metadata struct {
	stdao.Model
	CategoryID   string `json:"category_id"`
	Order        int    `json:"order"`
	Key          string `json:"key"`
	Type         string `json:"type"`
	Description  string `json:"description"`
	Value        string `json:"value"`
	DefaultValue string `json:"default_value"`
}
