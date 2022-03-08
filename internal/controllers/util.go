package controllers

import (
	"time"
)

const (
	HtmlLocalDateTimeFormat = "2006-01-02"
	OutputDateFormat        = "2nd Jan 2006"
)

// parsHtmlLocalDateTime
func parsHtmlLocalDateTime(dateString string) *string {
	date, err := time.Parse(HtmlLocalDateTimeFormat, dateString)
	if err != nil {
		return nil
	}

	formatted := date.Format("2nd Jan 2006")
	return &formatted
}
