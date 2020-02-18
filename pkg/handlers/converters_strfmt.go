package handlers

import (
	"time"

	"github.com/go-openapi/strfmt"
)

// These functions facilitate converting from the go types the db uses
// into the strfmt types that go-swagger uses for payloads.

// FmtDatePtrToPopPtr converts go-swagger type to pop type
func FmtDatePtrToPopPtr(date *strfmt.Date) *time.Time {
	if date == nil {
		return nil
	}

	fmtDate := time.Time(*date)
	return &fmtDate
}

// FmtDateTimePtrToPopPtr converts go-swagger type to pop type
func FmtDateTimePtrToPopPtr(date *strfmt.DateTime) *time.Time {
	if date == nil {
		return nil
	}

	fmtDate := time.Time(*date)
	return &fmtDate
}
