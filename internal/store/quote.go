package store

import (
	"gorm.io/gorm"
)

type Quote struct {
	Model
	Text   string  `gorm:"size:255" json:"text"`
	Person string  `gorm:"size:10,index" json:"person"`
	Date   *string `gorm:"size:15" json:"date"`
}

// CreateQuote on a TODO list
func CreateQuote(db *gorm.DB, text, person string, date *string) *Quote {
	quote := Quote{
		Text:   text,
		Person: person,
		Date:   date,
	}

	if tx := db.Create(&quote); tx.Error != nil {
		return nil
	}

	return &quote
}

// ListQuotes on a todo list
func ListQuotes(db *gorm.DB) []Quote {
	var entries []Quote

	_ = db.Order("id DESC").Find(&entries)

	return entries
}

// FindQuote by its unique id
func FindOuote(db *gorm.DB, id string) *Quote {
	var quote Quote

	if tx := db.First(&quote, id); tx.Error != nil {
		return nil
	}

	return &quote
}
