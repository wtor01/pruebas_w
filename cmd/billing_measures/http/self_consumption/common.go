package self_consumption

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
)

func selfConsumptionToResponse(self billing_measures.SelfConsumption) SelfConsumptionUnitConfig {

	scconfigs := make([]SelfConsumptionConfig, 0)

	for _, val := range self.Configs {
		confT := string(val.ConfType)
		connT := string(val.ConnType)
		generationP := float32(val.GenerationPot)
		scconfigs = append(scconfigs, SelfConsumptionConfig{
			AntivertType:          &val.AntivertType,
			CnmcTypeDesc:          &val.CnmcTypeDesc,
			CnmcTypeId:            &val.CnmcTypeId,
			CnmcTypeName:          &val.CnmcTypeName,
			Compensation:          &val.Compensation,
			ConfType:              &confT,
			ConfTypeDescription:   &val.ConfTypeDescription,
			ConnType:              &connT,
			ConsumerType:          &val.ConsumerType,
			EndDate:               &val.EndDate,
			Excedents:             &val.Excedents,
			GenerationPot:         &generationP,
			GroupSubgroup:         &val.GroupSubgroup,
			Id:                    &val.ID,
			InitDate:              &val.InitDate,
			ParticipantNumber:     &val.ParticipantNumber,
			SolarZoneId:           &val.SolarZoneId,
			SolarZoneNum:          &val.SolarZoneNum,
			SolarZoneName:         &val.SolarZoneName,
			StatusId:              &val.StatusID,
			StatusName:            &val.StatusName,
			TechnologyId:          &val.TechnologyId,
			TechnologyDescription: &val.TechnologyDescription,
		})
	}

	scpoints := make([]SelfConsumptionPoint, 0)

	for _, point := range self.Points {

		partitionCoe := float32(point.PartitionCoeff)
		servicePoint := string(point.ServicePointType)

		scpoints = append(scpoints, SelfConsumptionPoint{
			CUPS:             &point.CUPS,
			EndDate:          &point.EndDate,
			Exent1Flag:       &point.Exent1Flag,
			Exent2Flag:       &point.Exent2Flag,
			Id:               &point.ID,
			InitDate:         &point.InitDate,
			InstalationFlag:  &point.InstalationFlag,
			PartitionCoeff:   &partitionCoe,
			ServicePointType: &servicePoint,
			WithoutmeterFlag: &point.WithoutmeterFlag,
		})
	}

	return SelfConsumptionUnitConfig{
		Id:            self.ID,
		CAU:           self.CAU,
		Name:          self.Name,
		StatusId:      self.StatusID,
		StatusName:    self.StatusName,
		CcaaId:        self.CcaaId,
		Ccaa:          self.Ccaa,
		InitDate:      self.InitDate,
		EndDate:       self.EndDate,
		DistributorId: self.DistributorId,
		Configs:       scconfigs,
		Points:        scpoints,
	}
}

func billingSelfConsumptionToResponse(billingSelfConsumption billing_measures.BillingSelfConsumptionUnit) BillingSelfConsumptionUnitInfo {
	return BillingSelfConsumptionUnitInfo{
		Status:    string(billingSelfConsumption.Status),
		StartDate: billingSelfConsumption.InitDate,
		EndDate:   billingSelfConsumption.EndDate,
		CauInfo: BillingCauInfo{
			Id:         billingSelfConsumption.CauInfo.Id,
			ConfigType: string(billingSelfConsumption.CauInfo.ConfigType),
			Name:       billingSelfConsumption.CauInfo.Name,
			Points:     billingSelfConsumption.CauInfo.Points,
			UnitType:   billingSelfConsumption.CauInfo.UnitType,
		},
		CalendarConsumption: utils.MapSlice(billingSelfConsumption.CalendarConsumption, calendarConsumptionToResponse),
		NetGeneration:       utils.MapSlice(billingSelfConsumption.NetGeneration, netGenerationToResponse),
		Totals: BillingTotals{
			AuxConsumption:     billingSelfConsumption.Totals.AuxConsumption,
			GrossGeneration:    billingSelfConsumption.Totals.GrossGeneration,
			NetGeneration:      billingSelfConsumption.Totals.NetGeneration,
			NetworkConsumption: billingSelfConsumption.Totals.NetworkConsumption,
			SelfConsumption:    billingSelfConsumption.Totals.SelfConsumption,
		},
		UnitConsumption: utils.MapSlice(billingSelfConsumption.UnitConsumption, unitConsumptionToResponse),
		CupsList:        utils.MapSlice(billingSelfConsumption.Cups, cupsToResponse),
	}
}

func calendarConsumptionToResponse(calendarConsumption billing_measures.BillingSelfConsumptionCalendarConsumption) BillingCalendarConsumption {
	return BillingCalendarConsumption{
		Date:   calendarConsumption.Date,
		Energy: calendarConsumption.Energy,
		Values: utils.MapSlice(calendarConsumption.Values, calendarConsumptionValuesToResponse),
	}
}

func calendarConsumptionValuesToResponse(values billing_measures.CalendarConsumptionValues) BillingCalendarConsumptionValues {
	return BillingCalendarConsumptionValues{
		Energy: values.Energy,
		Hour:   values.Hour,
	}
}

func netGenerationToResponse(netGeneration billing_measures.BillingSelfConsumptionNetGeneration) BillingNetGeneration {
	return BillingNetGeneration{
		Date:      netGeneration.Date,
		Excedente: netGeneration.Excedente,
		Net:       netGeneration.Net,
	}
}

func unitConsumptionToResponse(unitConsumption billing_measures.BillingSelfConsumptionUnitConsumption) BillingUnitConsumption {
	return BillingUnitConsumption{
		Aux:     unitConsumption.Aux,
		Date:    unitConsumption.Date,
		Network: unitConsumption.Network,
		Self:    unitConsumption.Self,
	}
}

func cupsToResponse(cupsSelfConsumption billing_measures.BillingSelfConsumptionCups) BillingSelfConsumptionCups {
	return BillingSelfConsumptionCups{
		Consumption: cupsSelfConsumption.Consumption,
		Cups:        cupsSelfConsumption.Cups,
		EndDate:     cupsSelfConsumption.EndDate,
		Generation:  cupsSelfConsumption.Generation,
		PsType:      string(cupsSelfConsumption.Type),
		StartDate:   cupsSelfConsumption.StartDate,
	}
}
