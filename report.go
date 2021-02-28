package conav

import (
	"context"
)

type CaseReporter interface {
	GetCaseReport(ctx context.Context) (*CantonReport, error)
	GetCaseReportByMunicipality(ctx context.Context, name string) (*MunicipalityReport, error)
}

type CantonReport struct {
	NewCases            int                  `json:"new_cases"`
	MunicipalityReports []MunicipalityReport `json:"municipality_reports"`
}

type MunicipalityReport struct {
	Name       string `json:"name"`
	NewCases   int    `json:"new_cases"`
	Population int    `json:"population"`
}
