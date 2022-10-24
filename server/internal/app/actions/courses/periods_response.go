package courses

import "github.com/sava-cska/SPbSU-EMKN/internal/app/models"

type PeriodsResponseWrapper struct {
	Periods []*models.Period `json:"periods"`
}

type PeriodsResponse struct {
	Response *PeriodsResponseWrapper `json:"response"`
}

func (p PeriodsResponse) Bind() {}
