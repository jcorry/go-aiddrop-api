package controllers

import (
	"aiddrop-api/app/models"
	"net/http"
	"strconv"
	"time"

	"github.com/revel/revel"
)

type CasesCtrl struct {
	GorpController
}

type CasesError struct {
	Message string
}

// CasesRequest is the geo bounding box within which reports will be found
type CasesRequest struct {
	MinLatitude  float64 `json:"minLatitude"`
	MaxLatitude  float64 `json:"maxLatitude"`
	MinLongitude float64 `json:"minLongitude"`
	MaxLongitude float64 `json:"maxLongitude"`
}

// Validate validats the request params
func (r CasesRequest) Validate(v *revel.Validation) {
	v.Check(r.MinLatitude,
		revel.ValidRequired(),
		revel.ValidRange(-85, 85))

	v.Check(r.MaxLatitude,
		revel.ValidRequired(),
		revel.ValidRange(-85, 85))

	v.Check(r.MinLongitude,
		revel.ValidRequired(),
		revel.ValidRange(-180, 180))

	v.Check(r.MaxLongitude,
		revel.ValidRequired(),
		revel.ValidRange(-180, 180))
}

// List returns the cases within a defined region
func (c CasesCtrl) List() revel.Result {
	requestData := c.parseQueryRequest()
	requestData.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(c.Validation.Errors)
	}
	// Now that you have a valid geo bounding box
	// Get all of the reports within that box
	limit := 100

	reportAgeMinutes := revel.Config.IntDefault("case.reportAge", 90)
	reportAge := time.Now().Add(time.Duration(-reportAgeMinutes) * time.Minute)

	var reports []models.Report
	_, err := c.Txn.Select(&reports,
		`SELECT * FROM Report 
		WHERE latitude >= ? 
		AND latitude <= ? 
		AND longitude >= ? 
		AND longitude <= ?
		AND created >= ?
		ORDER BY id DESC
		LIMIT ?`, requestData.MinLatitude, requestData.MaxLatitude, requestData.MinLongitude, requestData.MaxLongitude, reportAge, limit)

	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJSON(err)
	}

	if len(reports) == 0 {
		c.Response.Status = http.StatusNoContent
		c.Response.ContentType = "application/json"
		return c.Render()
	}

	// Consolidate reports that are close together into a single Case
	caseModel := models.Case{}
	cases := caseModel.GetFromReports(reports)

	return c.RenderJSON(cases)
}

func (c CasesCtrl) parseQueryRequest() CasesRequest {
	minLatitude, err := strconv.ParseFloat(c.Params.Query.Get("minLatitude"), 64)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		c.RenderJSON(CasesError{Message: "Could not convert minLatitude to float."})
	}

	maxLatitude, err := strconv.ParseFloat(c.Params.Query.Get("maxLatitude"), 64)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		c.RenderJSON(CasesError{Message: "Could not convert maxLatitude to float."})
	}

	minLongitude, err := strconv.ParseFloat(c.Params.Query.Get("minLongitude"), 64)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		c.RenderJSON(CasesError{Message: "Could not convert minLongitude to float."})
	}

	maxLongitude, err := strconv.ParseFloat(c.Params.Query.Get("maxLongitude"), 64)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		c.RenderJSON(CasesError{Message: "Could not convert maxLongitude to float."})
	}

	return CasesRequest{
		MinLatitude:  minLatitude,
		MaxLatitude:  maxLatitude,
		MinLongitude: minLongitude,
		MaxLongitude: maxLongitude,
	}
}
