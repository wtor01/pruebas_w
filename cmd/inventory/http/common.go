package http

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"bitbucket.org/sercide/data-ingestion/internal/measures"
	"fmt"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"strconv"
)

func measureEquipmentToResponse(d inventory.MeasureEquipment) MeasureEquipment {

	return MeasureEquipment{
		Id:                d.ID.String(),
		SerialNumber:      &d.SerialNumber,
		Technology:        &d.Technology,
		Type:              &d.Type,
		Brand:             &d.Brand,
		Model:             &d.Model,
		ActiveConstant:    &d.ActiveConstant,
		ReactiveConstant:  &d.ReactiveConstant,
		MaximeterConstant: &d.MaximeterConstant,
	}
}
func meterConfigToResponse(mc measures.MeterConfig) MeterConfig {
	csed := openapi_types.Date{mc.ContractualSituations.EndDate}
	cssd := openapi_types.Date{mc.ContractualSituations.InitDate}
	ed := openapi_types.Date{mc.EndDate}
	sd := openapi_types.Date{mc.StartDate}
	csrc := strconv.Itoa(mc.ContractualSituations.RetailerCode)
	mplc := fmt.Sprint(mc.MeasurePoint.LossesCoef)
	mplp := fmt.Sprint(mc.MeasurePoint.LossesPerc)
	mpType := string(mc.MeasurePoint.Type)
	curveType := string(mc.CurveType)
	cups := mc.Cups()

	mac := mc.Meter.ActiveConstant
	MaximeterConstant := mc.Meter.MaximeterConstant
	ReactiveConstant := mc.Meter.ReactiveConstant
	pc, _ := strconv.ParseFloat(string(mc.PriorityContract), 32)
	PriorityContract := float32(pc)
	enabled := !mc.ServicePoint.Disabled
	pointType := string(mc.ServicePoint.PointType)
	serviceType := string(mc.ServicePoint.ServiceType)
	mType := string(mc.Type)
	tlg := string(mc.TlgCode)

	meterType := string(mc.Meter.Type)
	cs := ContractualSituation{
		EndDate:      &csed,
		Id:           mc.ContractualSituations.ID,
		P1Demand:     &mc.ContractualSituations.P1Demand,
		P2Demand:     &mc.ContractualSituations.P2Demand,
		P3Demand:     &mc.ContractualSituations.P3Demand,
		P4Demand:     &mc.ContractualSituations.P4Demand,
		P5Demand:     &mc.ContractualSituations.P5Demand,
		P6Demand:     &mc.ContractualSituations.P6Demand,
		RetailerCdos: &mc.ContractualSituations.RetailerCdos,
		RetailerCode: &csrc,
		RetailerName: &mc.ContractualSituations.RetailerName,
		StartDate:    &cssd,
		Tariff:       &mc.ContractualSituations.TariffID,
	}
	mp := MeasurePoint{
		Id:         mc.MeasurePoint.ID,
		LossesCoef: &mplc,
		LossesPerc: &mplp,
		Type:       &mpType,
	}
	meter := Meter{
		ActiveConstant:    &mac,
		Brand:             &mc.Meter.Brand,
		Id:                mc.Meter.ID,
		MaximeterConstant: &MaximeterConstant,
		Model:             &mc.Meter.Model,
		ReactiveConstant:  &ReactiveConstant,
		SerialNumber:      &mc.Meter.SerialNumber,
		SmakiaId:          &mc.Meter.SmarkiaID,
		Technology:        &mc.Meter.Technology,
		Type:              &meterType,
	}
	servicePoint := ServicePoint{
		Cups:                &mc.ServicePoint.Cups,
		Enabled:             &enabled,
		Id:                  mc.ServicePoint.ID,
		MeasureTensionLevel: &mc.ServicePoint.MeasureTensionLevel,
		Name:                &mc.ServicePoint.Name,
		PointTensionLevel:   &mc.ServicePoint.PointTensionLevel,
		PointType:           &pointType,
		ServiceType:         &serviceType,
		TensionSection:      &mc.ServicePoint.TensionSection,
	}

	return MeterConfig{
		Ae:                   &mc.AE,
		Ai:                   &mc.AI,
		Calendar:             &mc.CalendarID,
		ContractualSituation: &cs,
		Cups:                 &cups,
		CurveType:            &curveType,
		E:                    &mc.E,
		EndDate:              &ed,
		Id:                   mc.ID,
		MeasurePoint:         &mp,
		Meter:                &meter,
		Mx:                   &mc.M,
		PriorityContract:     &PriorityContract,
		R1:                   &mc.R1,
		R2:                   &mc.R2,
		R3:                   &mc.R3,
		R4:                   &mc.R4,
		ReadingType:          &mc.ReadingType,
		RentingPrince:        &mc.RentingPrice,
		ServicePoint:         &servicePoint,
		StartDate:            &sd,
		TlgCode:              &tlg,
		Type:                 &mType,
	}
}

func distributorToResponse(d inventory.Distributor) Distributor {
	return Distributor{
		Id:        d.ID.String(),
		Name:      &d.Name,
		R1:        &d.R1,
		SmarkiaId: &d.SmarkiaId,
	}
}
