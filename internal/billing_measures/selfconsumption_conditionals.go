package billing_measures

import (
	"context"
)

type IsSelfConsumption struct {
	ID string
}

func NewIsSelfConsumption() *IsSelfConsumption {
	return &IsSelfConsumption{
		ID: "IS_SELF_CONSUMPTION",
	}
}

func (i IsSelfConsumption) Eval(ctx context.Context) bool {
	return true
}

type IsConfigType struct {
	ID    string
	Types []ConfigType
	c     *BillingSelfConsumption `bson:"-"`
}

func NewIsConfigType(c *BillingSelfConsumption, types []ConfigType) *IsConfigType {
	return &IsConfigType{
		ID:    "IS_CONFIG_TYPE",
		Types: types,
		c:     c,
	}
}

func (i IsConfigType) Eval(ctx context.Context) bool {
	for _, t := range i.Types {
		if t == i.c.Config.ConfType {
			return true
		}
	}
	return false
}

type IsTypeConnection struct {
	ID   string
	Type ConnectionType
	c    *BillingSelfConsumption `bson:"-"`
}

func NewIsTypeConnection(c *BillingSelfConsumption, t ConnectionType) *IsTypeConnection {
	return &IsTypeConnection{
		ID:   "IS_TYPE_CONNECTION",
		Type: t,
		c:    c,
	}
}

func (i IsTypeConnection) Eval(ctx context.Context) bool {
	return i.Type == i.c.Config.ConnType
}

type IsIndividualGeneration struct {
	ID string
	c  *BillingSelfConsumption `bson:"-"`
}

func NewIsIndividualGeneration(c *BillingSelfConsumption) *IsIndividualGeneration {
	return &IsIndividualGeneration{
		ID: "IS_INDIVIDUAL_GENERATION",
		c:  c,
	}
}

func (i IsIndividualGeneration) Eval(ctx context.Context) bool {

	count := 0

	for _, sp := range i.c.Points {
		if sp.ServicePointType == GdServicePointType {
			count += 1
		}
	}

	if count > 1 {
		return false
	}

	return true
}

type IsIndividualConsumer struct {
	ID string
	c  *BillingSelfConsumption `bson:"-"`
}

func NewIsIndividualConsumer(c *BillingSelfConsumption) *IsIndividualConsumer {
	return &IsIndividualConsumer{
		ID: "IS_INDIVIDUAL_CONSUMER",
		c:  c,
	}
}

func (i IsIndividualConsumer) Eval(ctx context.Context) bool {
	//TODO : REFACTOR TYPES
	return i.c.Config.ConsumerType == "Individual"
}

type IsCompensation struct {
	ID string
	c  *BillingSelfConsumption `bson:"-"`
}

func NewIsCompensation(c *BillingSelfConsumption) *IsCompensation {
	return &IsCompensation{
		ID: "IS_COMPENSATION",
		c:  c,
	}
}

func (i IsCompensation) Eval(ctx context.Context) bool {
	return i.c.Config.Compensation
}

type AreExcedents struct {
	ID string
	c  *BillingSelfConsumption `bson:"-"`
}

func NewAreExcedents(c *BillingSelfConsumption) *AreExcedents {
	return &AreExcedents{
		ID: "ARE_EXCEDENTS",
		c:  c,
	}
}

func (i AreExcedents) Eval(ctx context.Context) bool {
	//TODO : REFACTOR TYPES
	return i.c.Config.Excedents
}

type IsSSAA struct {
	ID string
	c  *BillingSelfConsumption `bson:"-"`
}

func NewIsSSAA(c *BillingSelfConsumption) *IsSSAA {
	return &IsSSAA{
		ID: "IS_SSAA",
		c:  c,
	}
}

func (i IsSSAA) Eval(ctx context.Context) bool {
	for _, sp := range i.c.Points {
		if sp.ServicePointType == SsaaServicePointType {
			return true
		}
	}
	return false
}
