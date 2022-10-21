package smarkia_api

import (
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures"
	"bitbucket.org/sercide/data-ingestion/internal/gross_measures/smarkia"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"bitbucket.org/sercide/data-ingestion/pkg/telemetry"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ResponseGetDistributors struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	SearchSize int `json:"searchSize"`
	List       []struct {
		Id               int    `json:"id"`
		Sname            string `json:"sname"`
		CountryId        int    `json:"countryId"`
		ParentGroupingId int    `json:"parentGroupingId"`
		Status           string `json:"status"`
		Synchronizable   bool   `json:"synchronizable"`
	} `json:"list"`
}

type ResponseGetCts struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	SearchSize int `json:"searchSize"`
	List       []struct {
		Id               int    `json:"id"`
		Sname            string `json:"sname"`
		CountryId        int    `json:"countryId"`
		Province         string `json:"province"`
		City             string `json:"city"`
		Telephone        string `json:"telephone,omitempty"`
		ParentGroupingId int    `json:"parentGroupingId"`
		Path             string `json:"path"`
		Status           string `json:"status"`
		Synchronizable   bool   `json:"synchronizable"`
		Grouping         bool   `json:"grouping"`
	} `json:"list"`
}

type ResponseGetEquipments struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	SearchSize int `json:"searchSize"`
	List       []struct {
		Id             int    `json:"id"`
		Sname          string `json:"sname"`
		WorkcentreId   int    `json:"workcentreId"`
		Virtual        bool   `json:"virtual"`
		Active         bool   `json:"active"`
		ModelId        int    `json:"modelId"`
		Path           string `json:"path"`
		OriginalId     string `json:"originalId"`
		Synchronizable bool   `json:"synchronizable"`
		Tags           []int  `json:"tags"`
		SupplyId       int    `json:"supplyId"`
	} `json:"list"`
}

type ResponseMagnitudeByEquipment struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	SearchSize int `json:"searchSize"`
	List       []struct {
		Id                   int    `json:"id"`
		Sname                string `json:"sname"`
		MeasuringEquipmentId int    `json:"measuringEquipmentId"`
		WorkCentreId         int    `json:"workCentreId"`
		GeneralConsumption   bool   `json:"generalConsumption"`
		Active               bool   `json:"active"`
		Status               string `json:"status"`
		OriginalId           string `json:"originalId"`
		Synchronizable       bool   `json:"synchronizable"`
		Tags                 []int  `json:"tags"`
		InheritTags          []int  `json:"inheritTags"`
	} `json:"list"`
}

type ResponseCurva struct {
	OriginalId int `json:"originalId"`
	DataPoints []struct {
		Date      time.Time `json:"date"`
		Data      string    `json:"data"`
		Estimated bool      `json:"estimated"`
		Qualifier string    `json:"qualifier"`
		Period    string    `json:"period"`
	} `json:"dataPoints"`
}

type CurveDataReadingType string

const (
	HourlyCurveDataReadingType  = "HOURLY"
	QuarterCurveDataReadingType = "QUARTER"
)

type CurvaData struct {
	Date          time.Time
	Data          string
	Estimated     bool
	Qualifier     string
	Period        string
	OriginalId    int
	GroupingType  CurveDataReadingType
	MagnitudeName string
}

type ClosingData struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	SearchSize int `json:"searchSize"`
	List       []struct {
		Id                  int       `json:"id"`
		SupplyCode          string    `json:"supplyCode"`
		StartDate           time.Time `json:"startDate"`
		CloseDate           time.Time `json:"closeDate"`
		Ct                  string    `json:"ct"`
		Per                 string    `json:"per"`
		EnergyAbsA          string    `json:"energyAbsA"`
		EnergyIncA          string    `json:"energyIncA"`
		EnergyAbsRi         string    `json:"energyAbsRi"`
		EnergyIncRi         string    `json:"energyIncRi"`
		EnergyAbsRc         string    `json:"energyAbsRc"`
		EnergyIncRc         string    `json:"energyIncRc"`
		ExcessA             string    `json:"excessA"`
		PotMaxA             string    `json:"potMaxA"`
		Fecha               time.Time `json:"fecha"`
		EnergyIncAQualName  string    `json:"energyIncAQualName"`
		EnergyIncRiQualName string    `json:"energyIncRiQualName"`
		EnergyIncRcQualName string    `json:"energyIncRcQualName"`
		ExcessAQualName     string    `json:"excessAQualName"`
		PotMaxAQualName     string    `json:"potMaxAQualName"`
		Reserve7            string    `json:"reserve7"`
		Reserve8            string    `json:"reserve8"`
		Reserve7QualName    string    `json:"reserve7QualName"`
		Reserve8QualName    string    `json:"reserve8QualName"`
	} `json:"list"`
}

func NewApi(token string, host string, loc *time.Location) *Api {
	return &Api{token: token, host: host, location: loc}
}

type Api struct {
	location *time.Location
	token    string
	host     string
}

func (a Api) GetEquipments(ctx context.Context, dto smarkia.GetEquipmentsQuery) ([]smarkia.Equipment, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetEquipments")
	defer span.End()

	numberRequest := 1
	pageNumber := 1
	pageSize := 100
	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("%v/restapi-monitor/v2/measuringequipments", a.host)
	var err error
	ctsList := make([]smarkia.Equipment, 0)

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			break
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.token))

		pageNumber = pageNumber * numberRequest
		query := req.URL.Query()

		if dto.Id != "" {
			query.Add("workCentreId", dto.Id)
		} else {
			query.Add("sname", dto.Cups)
		}

		query.Add("pageNumber", strconv.Itoa(pageNumber))
		query.Add("pageSize", strconv.Itoa(pageSize))

		req.URL.RawQuery = query.Encode()

		res, err := client.Do(req)
		numberRequest += 1
		if err != nil {
			break
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			break
		}

		_ = res.Body.Close()

		var r ResponseGetEquipments
		err = json.Unmarshal(body, &r)
		if err != nil {
			break
		}

		for _, ct := range r.List {
			ctsList = append(ctsList, smarkia.Equipment{
				ID:   strconv.Itoa(ct.Id),
				CUPS: ct.Sname,
				CtId: strconv.Itoa(ct.WorkcentreId),
			})
		}

		if len(r.List) != r.PageSize {
			break
		}

	}

	return ctsList, err
}

func (a Api) GetCTs(ctx context.Context, distributorID string) ([]smarkia.Ct, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetCTs")
	defer span.End()
	numberRequest := 1
	pageNumber := 1
	pageSize := 100
	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("%v/restapi-monitor/v2/workcentres", a.host)
	var err error
	ctsList := make([]smarkia.Ct, 0)

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			break
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.token))

		pageNumber = pageNumber * numberRequest
		query := req.URL.Query()

		query.Add("parentGroupingId", distributorID)
		query.Add("pageNumber", strconv.Itoa(pageNumber))
		query.Add("pageSize", strconv.Itoa(pageSize))

		req.URL.RawQuery = query.Encode()

		res, err := client.Do(req)
		numberRequest += 1
		if err != nil {
			break
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			break
		}

		_ = res.Body.Close()

		var r ResponseGetCts
		err = json.Unmarshal(body, &r)
		if err != nil {
			break
		}

		for _, ct := range r.List {
			ctsList = append(ctsList, smarkia.Ct{
				ID: strconv.Itoa(ct.Id),
			})
		}

		if len(r.List) != r.PageSize {
			break
		}

	}

	return ctsList, err
}

func (a Api) GetDistributors(ctx context.Context) ([]smarkia.Distributor, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetDistributors")
	defer span.End()

	numberRequest := 1
	pageNumber := 1
	pageSize := 100
	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("%v/restapi-monitor/v2/groupings", a.host)
	var err error
	distributorList := make([]smarkia.Distributor, 0)

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			break
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.token))

		pageNumber = pageNumber * numberRequest
		query := req.URL.Query()

		query.Add("pageNumber", strconv.Itoa(pageNumber))
		query.Add("pageSize", strconv.Itoa(pageSize))
		req.URL.RawQuery = query.Encode()

		res, err := client.Do(req)
		numberRequest += 1

		if err != nil {
			break
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			break
		}

		_ = res.Body.Close()

		var r ResponseGetDistributors
		err = json.Unmarshal(body, &r)
		if err != nil {
			break
		}
		for _, distributor := range r.List {
			distributorList = append(distributorList, smarkia.Distributor{
				SmarkiaCode: strconv.Itoa(distributor.Id),
			})
		}

		if len(r.List) != r.PageSize {
			break
		}

	}

	return distributorList, err
}

func (a Api) getMagnitude(ctx context.Context, channel chan CurvaData, magnitudeId int, magnitudeName string, groupingType CurveDataReadingType, date time.Time) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "getMagnitude")
	defer span.End()

	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	url := fmt.Sprintf("%v/restapi-monitor/v2/measurements/%v/timeseries", a.host, magnitudeId)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.token))

	query := req.URL.Query()
	beginDate := date.AddDate(0, 0, -1)

	endDate := date
	query.Set("beginDate", beginDate.Format("2006-01-02T15:04:05.000Z07:00"))
	query.Set("endDate", endDate.Format("2006-01-02T15:04:05.000Z07:00"))
	query.Add("groupingType", string(groupingType))

	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	_ = res.Body.Close()

	var httpResponse ResponseCurva
	err = json.Unmarshal(body, &httpResponse)
	if err != nil {
		return
	}

	for _, point := range httpResponse.DataPoints {
		channel <- CurvaData{
			Date:          point.Date,
			Data:          point.Data,
			Estimated:     point.Estimated,
			Qualifier:     point.Qualifier,
			Period:        point.Period,
			OriginalId:    httpResponse.OriginalId,
			GroupingType:  groupingType,
			MagnitudeName: magnitudeName,
		}
	}
}

func (a Api) GetMagnitudes(ctx context.Context, equipmentID string, date time.Time) ([]gross_measures.MeasureCurveWrite, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetMagnitudes")
	defer span.End()

	numberRequest := 1
	pageNumber := 1
	pageSize := 100
	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("%v/restapi-monitor/v2/measurements", a.host)
	var err error

	channel := make(chan CurvaData)
	var wg sync.WaitGroup
	for {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			break
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.token))

		pageNumber = pageNumber * numberRequest
		query := req.URL.Query()

		query.Add("measuringEquipmentId", equipmentID)
		query.Add("pageNumber", strconv.Itoa(pageNumber))
		query.Add("pageSize", strconv.Itoa(pageSize))

		req.URL.RawQuery = query.Encode()

		res, err := client.Do(req)
		numberRequest += 1
		if err != nil {
			break
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			break
		}

		_ = res.Body.Close()

		var r ResponseMagnitudeByEquipment
		err = json.Unmarshal(body, &r)
		if err != nil {
			break
		}

		for _, rr := range r.List {
			if rr.OriginalId == "E1" || rr.OriginalId == "E2" || rr.OriginalId == "E3" || rr.OriginalId == "E4" || rr.OriginalId == "E5" || rr.OriginalId == "E6" {
				wg.Add(2)

				go func(ctx context.Context, channel chan CurvaData, magnitudeId int, originalId string) {
					defer wg.Done()
					a.getMagnitude(ctx, channel, magnitudeId, originalId, QuarterCurveDataReadingType, date)
				}(ctx, channel, rr.Id, rr.OriginalId)

				go func(ctx context.Context, channel chan CurvaData, magnitudeId int, originalId string) {
					defer wg.Done()
					a.getMagnitude(ctx, channel, magnitudeId, originalId, HourlyCurveDataReadingType, date)
				}(ctx, channel, rr.Id, rr.OriginalId)
			}
		}

		if len(r.List) != r.PageSize {
			break
		}
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	measuresProcessed := make(map[string]gross_measures.MeasureCurveWrite)

	now := time.Now()
	for c := range channel {
		key := fmt.Sprintf("%v|%v", c.Date, c.GroupingType)

		m, ok := measuresProcessed[key]
		if !ok {
			m = gross_measures.MeasureCurveWrite{
				EndDate:     c.Date.UTC(),
				ReadingDate: now,
				Type:        measures.Incremental,
				ReadingType: measures.Curve,
				CurveType:   measures.HourlyMeasureCurveReadingType,
				Origin:      measures.STM,
				File:        a.generateFile(now),
				Qualifier:   c.Qualifier,
			}
			if c.GroupingType == QuarterCurveDataReadingType {
				m.CurveType = measures.QuarterMeasureCurveReadingType
			}
		}

		switch c.MagnitudeName {
		case "E1":
			{
				val, _ := strconv.ParseFloat(c.Data, 64)
				m.AI = val * 1000
				break
			}
		case "E2":
			{
				val, _ := strconv.ParseFloat(c.Data, 64)
				m.AE = val * 1000
				break
			}
		case "E3":
			{
				val, _ := strconv.ParseFloat(c.Data, 64)
				m.R1 = val * 1000
				break
			}
		case "E4":
			{
				val, _ := strconv.ParseFloat(c.Data, 64)
				m.R2 = val * 1000
				break
			}
		case "E5":
			{
				val, _ := strconv.ParseFloat(c.Data, 64)
				m.R3 = val * 1000
				break
			}
		case "E6":
			{
				val, _ := strconv.ParseFloat(c.Data, 64)
				m.R4 = val * 1000
				break
			}
		}
		measuresProcessed[key] = m
	}

	measureList := make([]gross_measures.MeasureCurveWrite, 0, len(measuresProcessed))

	for _, v := range measuresProcessed {
		measureList = append(measureList, v)
	}

	sort.Slice(measureList, func(i, j int) bool {
		return measureList[i].EndDate.Before(measureList[j].EndDate)
	})

	return measureList, err
}

func (a Api) generateFile(date time.Time) string {
	return fmt.Sprintf("SMARKIA/%s", date.In(a.location).Format("2006-01-02"))
}

func (a Api) GetClosinger(ctx context.Context, equipmentID string, date time.Time) ([]gross_measures.MeasureCloseWrite, error) {
	tracer := telemetry.GetTracer()
	ctx, span := tracer.Start(ctx, "GetClosinger")
	defer span.End()

	numberRequest := 1
	pageNumber := 1
	pageSize := 100
	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("%v/restapi-monitor/v2/closures", a.host)

	var err error

	measuresClosed := make(map[string]gross_measures.MeasureCloseWrite, 0)
	measuresClosedPeriods := make(map[string]gross_measures.MeasureClosePeriod, 0)
	separatorPeriod := "|||"
	now := time.Now()

	beginDate := date.AddDate(0, -1, -date.Day())

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			break
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", a.token))

		pageNumber = pageNumber * numberRequest
		query := req.URL.Query()

		query.Add("deviceId", equipmentID)
		query.Add("protocol", "IEC8705102")
		query.Set("beginDate", beginDate.Format("2006-01-02T15:04:05.000Z07:00"))
		query.Set("endDate", date.AddDate(0, 1, 5).Format("2006-01-02T15:04:05.000Z07:00"))
		query.Add("pageNumber", strconv.Itoa(pageNumber))
		query.Add("pageSize", strconv.Itoa(pageSize))
		req.URL.RawQuery = query.Encode()

		res, err := client.Do(req)
		numberRequest += 1
		if err != nil {
			break
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			break
		}

		_ = res.Body.Close()

		var cd ClosingData
		err = json.Unmarshal(body, &cd)
		if err != nil {
			break
		}

		for _, measure := range cd.List {
			// El contrato 2 no lo estamos utilizando para la medida, por lo que lo descartamos
			if measure.Ct == "II" {
				continue
			}
			ai, _ := strconv.ParseFloat(measure.EnergyAbsA, 64)
			ae, _ := strconv.ParseFloat(measure.EnergyAbsA, 64)
			r1, _ := strconv.ParseFloat(measure.EnergyIncRi, 64)
			r2, _ := strconv.ParseFloat(measure.EnergyIncRc, 64)
			r3, _ := strconv.ParseFloat(measure.EnergyIncRi, 64)
			r4, _ := strconv.ParseFloat(measure.EnergyIncRc, 64)
			m := gross_measures.MeasureCloseWrite{
				StartDate:      time.Date(measure.StartDate.Year(), measure.StartDate.Month(), measure.StartDate.Day(), 0, 0, 0, 0, a.location).UTC(),
				EndDate:        time.Date(measure.CloseDate.Year(), measure.CloseDate.Month(), measure.CloseDate.Day(), 0, 0, 0, 0, a.location).UTC(),
				ReadingDate:    now,
				Type:           measures.Absolute,
				ReadingType:    measures.BillingClosure,
				ConcentratorID: measure.Ct,
				File:           a.generateFile(now),
				Origin:         measures.STM,
				Periods:        make([]gross_measures.MeasureClosePeriod, 0, 6),
			}
			if measure.Per == "T" {
				measure.Per = "0"
			}
			keyPeriod := fmt.Sprintf("%v%v%s%s%s%s", m.StartDate, m.EndDate, m.MeterSerialNumber, m.ConcentratorID, separatorPeriod, measure.Per)
			keyMeasure := fmt.Sprintf("%v%v%s%s", m.StartDate, m.EndDate, m.MeterSerialNumber, m.ConcentratorID)
			period, _ := measuresClosedPeriods[keyPeriod]
			period.Period = measures.PeriodKey(fmt.Sprintf("P%s", measure.Per))

			if measure.Ct == "I" {
				mx, _ := strconv.ParseFloat(measure.PotMaxA, 64)
				e, _ := strconv.ParseFloat(measure.ExcessA, 64)

				period.AI = ai * 1000
				period.R1 = r1 * 1000
				period.R2 = r2 * 1000
				period.MX = mx * 1000
				period.FX = measure.Fecha
				period.E = e * 10000
			} else if measure.Ct == "III" {
				period.AE = ae * 1000
				period.R4 = r4 * 1000
				period.R2 = r2 * 1000
				period.R3 = r3 * 1000
			}
			measuresClosedPeriods[keyPeriod] = period
			measuresClosed[keyMeasure] = m
		}

		if len(cd.List) != cd.PageSize {
			break
		}
	}

	for keyPeriod, p := range measuresClosedPeriods {
		keys := strings.Split(keyPeriod, separatorPeriod)
		if len(keys) != 2 {
			continue
		}
		keyMeasure := keys[0]
		m, ok := measuresClosed[keyMeasure]
		if !ok {
			continue
		}
		m.Periods = append(m.Periods, p)
		measuresClosed[keyMeasure] = m
	}

	return utils.MapToSlice(measuresClosed), err
}
