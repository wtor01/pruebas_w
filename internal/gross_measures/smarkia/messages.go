package smarkia

import (
	"bitbucket.org/sercide/data-ingestion/pkg/event"
	"time"
)

type ProcessName = string

const (
	ProcessCurve ProcessName = "curve"
	ProcessClose ProcessName = "close"
)

type Types = string

const (
	TypeRequestSmarkia     = "SMARKIA/REQUEST"
	TypeDistributorProcess = "SMARKIA/DISTRIBUTOR"
	TypeCtProcess          = "SMARKIA/CT"
	TypeEquipmentProcess   = "SMARKIA/EQUIPMENT"
)

type RequestSmarkiaPayload struct {
	ProcessName string    `json:"process_name"`
	Date        time.Time `json:"date"`
}

type RequestSmarkiaEvent = event.Message[RequestSmarkiaPayload]

type MessageDistributorProcessPayload struct {
	DistributorId   string    `json:"distributor_id"`
	ProcessName     string    `json:"process_name"`
	SmarkiaId       string    `json:"smarkia_id"`
	DistributorCDOS string    `json:"distributor_cdos"`
	Date            time.Time `json:"date"`
}

type MessageDistributorProcess = event.Message[MessageDistributorProcessPayload]

func NewMessageDistributorProcess(distributorID, smarkiaId, processName, distributorCDOS string, date time.Time) MessageDistributorProcess {
	return MessageDistributorProcess{
		Type: TypeDistributorProcess,
		Payload: MessageDistributorProcessPayload{
			DistributorId:   distributorID,
			ProcessName:     processName,
			SmarkiaId:       smarkiaId,
			DistributorCDOS: distributorCDOS,
			Date:            date,
		},
	}
}

type CtProcessEvent = event.Message[CtProcessEventPayload]

type CtProcessEventPayload struct {
	MessageDistributorProcessPayload
	CtId string `json:"ct_id"`
}

func NewCtProcessEvent(distributorID, smarkiaId, processName, distributorCDOS, CtId string, date time.Time) CtProcessEvent {
	return CtProcessEvent{
		Type: TypeCtProcess,
		Payload: CtProcessEventPayload{
			MessageDistributorProcessPayload: MessageDistributorProcessPayload{
				DistributorId:   distributorID,
				ProcessName:     processName,
				DistributorCDOS: distributorCDOS,
				SmarkiaId:       smarkiaId,
				Date:            date,
			},
			CtId: CtId,
		},
	}
}

type EquipmentProcessEventPayload struct {
	CtProcessEventPayload
	EquipmentId string `json:"equipment_id"`
	CUPS        string `json:"cups"`
}

type MessageEquipmentDto struct {
	ProcessName     string
	DistributorId   string
	SmarkiaId       string
	CtId            string
	DistributorCDOS string
	Date            time.Time
}

type EquipmentProcessEvent = event.Message[EquipmentProcessEventPayload]

func NewEquipmentProcessEvent(dto MessageEquipmentDto, EquipmentId, cups string) EquipmentProcessEvent {
	return EquipmentProcessEvent{
		Type: TypeEquipmentProcess,
		Payload: EquipmentProcessEventPayload{
			CtProcessEventPayload: CtProcessEventPayload{
				MessageDistributorProcessPayload: MessageDistributorProcessPayload{
					DistributorId:   dto.DistributorId,
					ProcessName:     dto.ProcessName,
					DistributorCDOS: dto.DistributorCDOS,
					SmarkiaId:       dto.SmarkiaId,
					Date:            dto.Date,
				},
				CtId: dto.CtId,
			},
			EquipmentId: EquipmentId,
			CUPS:        cups,
		},
	}
}
