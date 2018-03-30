package controllers

import (
	"aiddrop-api/app/models"
	"net/http"
	"time"

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
	report := c.parseReport()
	report.Created = time.Now().Unix()

	log := c.Log.New("method", "ReportsCtrl.Create")

	// Validate the model
	report.Validate(c.Validation)
	if c.Validation.HasErrors() {
		if revel.DevMode {
			log.Debugf("%v", report)
		}
		return c.RenderJSON(c.Validation.Errors)
	} else {
		if err := c.Txn.Insert(&report); err != nil {
			return c.RenderJSON(ReportsError{Message: err.Error()})
		} else {
			return c.RenderJSON(report)
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
func (c ReportsCtrl) List() revel.Result {
	lastId := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	reports, err := c.Txn.Select(models.Report{},
		`SELECT * FROM Report WHERE id > ? LIMIT ?`, lastId, limit)

	if err != nil {
		return c.RenderJSON(ReportsError{Message: err.Error()})
	}

	return c.RenderJSON(reports)
}

// Delete is the deleteResource handler
func (c ReportsCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Report{ID: id})
	if err != nil || success == 0 {
		return c.RenderJSON(ReportsError{Message: "Failed to delete report."})
	}
	c.Response.Status = http.StatusNoContent
	return c.Render()
}

func (c ReportsCtrl) parseReport() models.Report {
	report := models.Report{}
	c.Params.BindJSON(&report)
	return report
}
