package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/validations"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type ValidationRule struct {
	ModelEntity
	Name                 string                 `gorm:"column:name"`
	Params               JSONMap                `gorm:"column:params"`
	Action               string                 `gorm:"column:action;type:validation_rules_action;not null;"`
	Enabled              bool                   `gorm:"column:enabled"`
	ReadingType          string                 `gorm:"column:reading_type;type:validation_rules_reading_type;not null;"`
	Type                 string                 `gorm:"column:type;type:validation_rules_type;not null;"`
	Code                 string                 `gorm:"column:code;unique"`
	Message              string                 `gorm:"column:message"`
	Description          string                 `gorm:"column:description"`
	ValidationRuleConfig []ValidationRuleConfig `gorm:"foreignKey:ValidationRuleID"`
	ModelRegisterUserActions
}

func (v ValidationRule) toDomain() validations.ValidationMeasure {
	updatedByID := ""

	if v.UpdatedByID != nil {
		updatedByID = *v.UpdatedByID
	}
	b, _ := json.Marshal(v.Params)

	var params validations.ValidationMeasureParams

	_ = json.Unmarshal(b, &params)

	return validations.ValidationMeasure{
		Id:          v.ID.String(),
		Name:        v.Name,
		Action:      v.Action,
		Enabled:     v.Enabled,
		MeasureType: v.ReadingType,
		Type:        v.Type,
		Code:        v.Code,
		Message:     v.Message,
		Description: v.Description,
		Params:      params,
		CreatedByID: v.CreatedByID,
		UpdatedByID: updatedByID,
	}
}

func validationRuleToDB(v validations.ValidationMeasure) ValidationRule {

	b, _ := json.Marshal(v.Params)

	params := make(map[string]interface{})

	_ = json.Unmarshal(b, &params)

	var updatedBy string

	if value, err := uuid.FromString(v.UpdatedByID); err == nil {
		updatedBy = value.String()
	}

	return ValidationRule{
		ModelEntity: ModelEntity{
			ID:        uuid.FromStringOrNil(v.Id),
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Name:        v.Name,
		Params:      params,
		Action:      v.Action,
		Enabled:     v.Enabled,
		ReadingType: v.MeasureType,
		Type:        v.Type,
		Code:        v.Code,
		Message:     v.Message,
		Description: v.Description,
		ModelRegisterUserActions: ModelRegisterUserActions{
			CreatedByID: v.CreatedByID,
			UpdatedByID: &updatedBy,
		},
	}
}

type ValidationRuleConfig struct {
	ModelEntity
	DistributorID    string         `gorm:"column:distributor_id;type:uuid;"`
	Distributor      Distributor    `gorm:"distributor_id" gorm:"foreignKey:DistributorID"`
	ValidationRuleID string         `gorm:"column:validation_rule_id;type:uuid;"`
	ValidationRule   ValidationRule `gorm:"validation_rule_id" gorm:"foreignKey:ValidationRuleID"`

	Params  JSONMap `gorm:"column:params"`
	Action  string  `gorm:"column:action;type:validation_rules_action;not null;"`
	Enabled bool    `gorm:"column:enabled"`

	ModelRegisterUserActions
}

func (v ValidationRuleConfig) toDomain() validations.ValidationMeasureConfig {
	updatedByID := ""

	if v.UpdatedByID != nil {
		updatedByID = *v.UpdatedByID
	}
	b, _ := json.Marshal(v.Params)

	var params validations.ValidationMeasureParams

	_ = json.Unmarshal(b, &params)

	return validations.ValidationMeasureConfig{
		Id:                v.ID.String(),
		DistributorID:     v.DistributorID,
		ValidationMeasure: v.ValidationRule.toDomain(),
		Action:            v.Action,
		Enabled:           v.Enabled,
		Params:            params,
		CreatedByID:       v.CreatedByID,
		UpdatedByID:       updatedByID,
	}
}

func validationRuleConfigToDB(v validations.ValidationMeasureConfig) ValidationRuleConfig {

	b, _ := json.Marshal(v.Params)

	params := make(map[string]interface{})

	_ = json.Unmarshal(b, &params)

	var updatedBy string

	if value, err := uuid.FromString(v.UpdatedByID); err == nil {
		updatedBy = value.String()
	}

	return ValidationRuleConfig{
		ModelEntity: ModelEntity{
			ID:        uuid.FromStringOrNil(v.Id),
			CreatedAt: time.Now(),
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		DistributorID:    v.DistributorID,
		ValidationRuleID: v.ValidationMeasure.Id,
		ValidationRule:   validationRuleToDB(v.ValidationMeasure),
		Params:           params,
		Action:           v.Action,
		Enabled:          v.Enabled,
		ModelRegisterUserActions: ModelRegisterUserActions{
			CreatedByID: v.CreatedByID,
			UpdatedByID: &updatedBy,
		},
	}
}

type HistoryValidationRuleConfig struct {
	ModelEntity
	ValidationRuleConfigID string               `gorm:"column:validation_rule_config_id;type:uuid"`
	ValidationRuleConfig   ValidationRuleConfig `gorm:"column:validation_rule_config_id" gorm:"foreignKey:ValidationRuleConfigID"`
	UserID                 uuid.UUID            `gorm:"column:user_id;type:uuid"`
	User                   User                 `gorm:"column:user_id" gorm:"foreignKey:UserID"`
	Params                 JSONMap              `gorm:"column:params"`
	Action                 string               `gorm:"column:action;type:validation_rules_action;not null;"`
	Enabled                bool                 `gorm:"column:enabled"`
}

func (HistoryValidationRuleConfig) TableName() string {
	return "hist_validation_rules_configs"
}
