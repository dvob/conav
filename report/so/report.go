package so

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dvob/conav"
)

var _ conav.CaseReporter = &CaseReporter{}

const ReportURL = "https://corona.so.ch/bevoelkerung/daten/fallzahlen-nach-gemeinden/"

type CaseReporter struct {
	report    *conav.CantonReport
	lastFetch time.Time
	reportURL string
}

func NewCaseReporter() *CaseReporter {
	return &CaseReporter{
		reportURL: ReportURL,
	}
}

func (cr *CaseReporter) GetCaseReport(ctx context.Context) (*conav.CantonReport, error) {
	if cr.report != nil && cr.lastFetch.Add(time.Hour*2).Before(time.Now()) {
		return cr.report, nil
	}
	return cr.fetch()
}

func (cr *CaseReporter) GetCaseReportByMunicipality(ctx context.Context, name string) (*conav.MunicipalityReport, error) {
	report, err := cr.GetCaseReport(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range report.MunicipalityReports {
		if m.Name == name {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("municipality with name '%s' not found")
}

func (cr *CaseReporter) fetch() (*conav.CantonReport, error) {
	resp, err := http.Get(cr.reportURL)
	if err != nil {
		return nil, err
	}
	return parseReport(resp.Body)
}

func parseReport(r io.Reader) (*conav.CantonReport, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	report := &conav.CantonReport{
		MunicipalityReports: []conav.MunicipalityReport{},
	}

	var parseErr error
	doc.Find("table").Each(func(j int, s *goquery.Selection) {
		// For each item found, get the band and title
		if strings.Contains(s.Find("td").First().Text(), "Laborbestätigte Infektionen") {
			totalNewCases, err := parseInt(s.Find("td").Eq(2).Text())
			if err != nil {
				parseErr = err
				return
			}
			report.NewCases = totalNewCases
		} else if strings.Contains(s.Find("td").First().Text(), "Gemeinde") {
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				tds := s.Find("td")
				gemeinde := tds.Eq(0).Text()
				gemeinde = strings.TrimSpace(gemeinde)
				gemeinde = strings.TrimRight(gemeinde, "*")
				if gemeinde == "Gemeinde" || gemeinde == "Total" || gemeinde == "Übrige Gemeinden" || gemeinde == "" {
					return
				}
				population, err := parseInt(tds.Eq(1).Text())
				if err != nil {
					//fmt.Println("gemeinde", i, "bezirk", j)
					parseErr = err
					return
				}
				diff, err := parseInt(tds.Eq(3).Text())
				if err != nil {
					//fmt.Println("gemeinde", i, "bezirk", j)
					parseErr = err
					return
				}
				report.MunicipalityReports = append(report.MunicipalityReports, conav.MunicipalityReport{
					Name:       gemeinde,
					NewCases:   diff,
					Population: population,
				})
			})
		}
	})

	if parseErr != nil {
		return nil, parseErr
	}
	return report, nil
}

func parseInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "'", "")
	return strconv.Atoi(s)
}
