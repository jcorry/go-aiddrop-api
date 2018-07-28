package models

import (
	"time"

	"github.com/revel/revel"
)

// Report represents a user report of someone needing aid
type Report struct {
	ID              int64     `db:"id" json:"id"`
	Latitude        float64   `db:"latitude" json:"latitude"`
	Longitude       float64   `db:"longitude" json:"longitude"`
	Description     string    `db:"description" json:"description"`
	RecipientsCount int64     `db:"recipients_count" json:"recipientsCount"`
	Created         time.Time `db:"created" json:"created"`
}

// Validate will validate the Report struct
func (r *Report) Validate(v *revel.Validation) {
	v.Check(r.Latitude,
		revel.ValidRequired(),
		revel.ValidRange(-80, 85))

	v.Check(r.Longitude,
		revel.ValidRequired(),
		revel.ValidRange(-180, 180))

	v.Check(r.Description,
		revel.ValidMaxSize(255))
}
