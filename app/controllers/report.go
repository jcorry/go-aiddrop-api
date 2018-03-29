package controllers

import (
	"aiddrop-api/app/models"
	"encoding/json"

	"github.com/revel/revel"
)

// ReportsCtrl controller is the handler for requests made to the /reports endpoints
type ReportsCtrl struct {
	GorpController
}

// ReportsError object returned when there's an error
type ReportsError struct {
	Message string
}

// Create is the createResource handler
func (c ReportsCtrl) Create() revel.Result {
	if report, err := c.parseReport(); err != nil {
		return c.RenderJSON(ReportsError{Message: "Unable to parse report from JSON."})
	} else {
		// Validate the model
		report.Validate(c.Validation)
		if c.Validation.HasErrors() {
			return c.RenderJSON(c.Validation.Errors)
		} else {
			if err := c.Txn.Insert(&report); err != nil {
				return c.RenderJSON(ReportsError{Message: err.Error()})
			} else {
				return c.RenderJSON(report)
			}
		}
	}
}

// Get readResource handler
func (c ReportsCtrl) Get(id int64) revel.Result {
	report := new(models.Report)
	err := c.Txn.SelectOne(report,
		`SELECT * FROM Report WHERE id = ?`, id)
	if err != nil {
		return c.RenderJSON(ReportsError{Message: "Error. Item doesn't exist."})
	}
	return c.RenderJSON(report)
}

// List is the listResource handler
func (r ReportsCtrl) List() revel.Result {

}

func (r ReportsCtrl) Delete() revel.Result {

}

// Report request parser
func (c ReportsCtrl) parseReport() (models.Report, error) {
	report := models.Report{}
	err := json.NewDecoder(c.Request.GetBody()).Decode(&report)
	return report, err
}
