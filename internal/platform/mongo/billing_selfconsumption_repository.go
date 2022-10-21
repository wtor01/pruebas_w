package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewBillingSelfConsumptionRepository(client *mongo.Client, database string, loc *time.Location) *BillingSelfConsumptionRepository {
	return &BillingSelfConsumptionRepository{
		client:     client,
		location:   loc,
		Database:   database,
		Collection: "billing_selfconsumption",
	}
}

type BillingSelfConsumptionRepository struct {
	client     *mongo.Client
	location   *time.Location
	Database   string
	Collection string
}

func (repository BillingSelfConsumptionRepository) GetBySelfConsumptionBetweenDates(ctx context.Context, query billing_measures.QueryGetBillingSelfConsumption) (billing_measures.BillingSelfConsumption, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.Collection)

	filter := bson.D{
		{"selfconsumption_id", query.SelfconsumptionId},
		{"init_date", bson.D{{"$lte", query.EndDate}}},
		{"end_date", bson.D{{"$gte", query.StartDate}}},
		{"status", bson.D{
			{"$in",
				bson.A{
					billing_measures.PendingPoints,
					billing_measures.CalculatingPoints,
					billing_measures.SupervisionSelfConsumptionStatus,
				},
			}},
		},
	}
	var selfConsumption BillingSelfConsumption

	result := coll.FindOne(ctxTimeout, filter)

	if result.Err() != nil {
		return billing_measures.BillingSelfConsumption{}, result.Err()
	}

	err := result.Decode(&selfConsumption)

	if err != nil {
		return billing_measures.BillingSelfConsumption{}, err
	}

	selfConsumptionDomain, err := selfConsumption.toDomain()

	return selfConsumptionDomain, err
}

func (repository BillingSelfConsumptionRepository) Save(ctx context.Context, b billing_measures.BillingSelfConsumption) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.Collection)
	filter := bson.D{
		{"_id", b.Id},
	}

	upsert := true
	_, err := coll.UpdateOne(ctxTimeout, filter, bson.D{
		{"$set", repository.toDb(b)}}, &options.UpdateOptions{
		Upsert: &upsert,
	})

	return err
}

func (repository BillingSelfConsumptionRepository) GetSelfConsumptionByCau(ctx context.Context, query billing_measures.QueryGetBillingSelfConsumptionByCau) ([]billing_measures.BillingSelfConsumption, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	collection := repository.client.Database(repository.Database).Collection(repository.Collection)

	matchStage := bson.D{
		{"$match", bson.D{
			{"distributor_id", query.DistributorId},
			{"CAU", query.CauId},
			{"end_date", bson.M{
				"$gt":  query.StartDate,
				"$lte": query.EndDate,
			},
			},
		}},
	}

	groupStage := bson.D{
		{
			"$group", bson.D{
				{
					"_id", bson.M{
						"end_date":  "$end_date",
						"init_date": "$init_date",
					},
				},
				{
					"first", bson.M{
						"$first": "$$ROOT",
					},
				},
			},
		},
	}

	repleceRootStage := bson.D{
		{
			"$replaceRoot", bson.D{
				{
					"newRoot", "$first",
				},
			},
		},
	}

	aggregations := append(bson.A{}, matchStage, groupStage, repleceRootStage)
	var selfConsumption []BillingSelfConsumption

	cursor, err := collection.Aggregate(ctxTimeout, aggregations)

	if err != nil {
		return []billing_measures.BillingSelfConsumption{}, err
	}

	if err = cursor.All(ctx, &selfConsumption); err != nil {
		return []billing_measures.BillingSelfConsumption{}, err
	}

	return utils.MapSlice(selfConsumption, func(item BillingSelfConsumption) billing_measures.BillingSelfConsumption {
		selfConsumptionDomain, _ := item.toDomain()
		return selfConsumptionDomain
	}), nil
}

/***************

	MODELS

****************/

type BillingSelfConsumption struct {
	Id                string                                        `json:"_id" bson:"_id"`
	SelfconsumptionId string                                        `json:"selfconsumption_id" bson:"selfconsumption_id"`
	CAU               string                                        `json:"CAU" bson:"CAU"`
	Name              string                                        `json:"name" bson:"name"`
	DistributorId     string                                        `json:"distributor_id" bson:"distributor_id"`
	InitDate          time.Time                                     `json:"init_date" bson:"init_date"`
	EndDate           time.Time                                     `json:"end_date" bson:"end_date"`
	GenerationDate    time.Time                                     `json:"generation_date" bson:"generation_date"`
	Status            billing_measures.BillingSelfConsumptionStatus `json:"status" bson:"status"`
	Config            BillingSelfConsumptionConfig                  `json:"config" bson:"config"`
	Points            []BillingSelfConsumptionPoint                 `json:"points" bson:"points"`
	Curve             []BillingSelfConsumptionCurve                 `json:"curve" bson:"curve"`
	GraphHistory      *billing_measures.Graph                       `json:"graph_history" bson:"graph_history"`
	SurplusType       billing_measures.SurplusType                  `json:"surplus_type" bson:"surplus_type"`
}

func (b BillingSelfConsumption) toDomain() (billing_measures.BillingSelfConsumption, error) {
	return billing_measures.BillingSelfConsumption{
		Id:                b.Id,
		SelfconsumptionId: b.SelfconsumptionId,
		CAU:               b.CAU,
		Name:              b.Name,
		DistributorId:     b.DistributorId,
		InitDate:          b.InitDate,
		EndDate:           b.EndDate,
		GenerationDate:    b.GenerationDate,
		Status:            b.Status,
		Config:            b.Config.toDomain(),
		Points: utils.MapSlice(b.Points, func(item BillingSelfConsumptionPoint) billing_measures.BillingSelfConsumptionPoint {
			return item.toDomain()
		}),
		Curve: utils.MapSlice(b.Curve, func(item BillingSelfConsumptionCurve) billing_measures.BillingSelfConsumptionCurve {
			return item.toDomain()
		}),
		SurplusType:  b.SurplusType,
		GraphHistory: b.GraphHistory,
	}, nil
}

type BillingSelfConsumptionPointMvhReceived struct {
	BillingMeasureId string    `json:"billing_measure_id" bson:"billing_measure_id"`
	InitDate         time.Time `json:"init_date" bson:"init_date"`
	EndDate          time.Time `json:"end_date" bson:"end_date"`
}

type BillingSelfConsumptionPoint struct {
	Id               string                                   `json:"_id" bson:"_id"`
	ServicePointType string                                   `json:"service_point_type" bson:"service_point_type"`
	CUPS             string                                   `json:"CUPS" bson:"CUPS"`
	CUPSgd           *string                                  `json:"CUPS_gd" bson:"CUPS_gd"`
	InstalationFlag  int                                      `json:"instalation_flag" bson:"instalation_flag"`
	WithoutmeterFlag int                                      `json:"withoutmeter_flag" bson:"withoutmeter_flag"`
	Exent1Flag       int                                      `json:"exent1_flag" bson:"exent1_flag"`
	Exent2Flag       int                                      `json:"exent2_flag" bson:"exent2_flag"`
	PartitionCoeff   float64                                  `json:"partition_coeff" bson:"partition_coeff"`
	MvhReceived      []BillingSelfConsumptionPointMvhReceived `json:"mvh_received" bson:"mvh_received"`
	InitDate         time.Time                                `json:"init_date,omitempty" bson:"init_date,omitempty"`
	EndDate          time.Time                                `json:"end_date,omitempty" bson:"end_date,omitempty"`
}

func (b BillingSelfConsumptionPoint) toDomain() billing_measures.BillingSelfConsumptionPoint {
	return billing_measures.BillingSelfConsumptionPoint{
		ID:               b.Id,
		ServicePointType: billing_measures.PointSelfConsumptionType(b.ServicePointType),
		CUPS:             b.CUPS,
		CUPSgd:           b.CUPSgd,
		InstalationFlag:  b.InstalationFlag,
		WithoutmeterFlag: b.WithoutmeterFlag,
		Exent1Flag:       b.Exent1Flag,
		Exent2Flag:       b.Exent2Flag,
		PartitionCoeff:   b.PartitionCoeff,
		MvhReceived: utils.MapSlice(b.MvhReceived, func(item BillingSelfConsumptionPointMvhReceived) billing_measures.BillingSelfConsumptionPointMvhReceived {
			return billing_measures.BillingSelfConsumptionPointMvhReceived{
				BillingMeasureId: item.BillingMeasureId,
				InitDate:         item.InitDate,
				EndDate:          item.EndDate,
			}
		}),
		InitDate: b.InitDate,
		EndDate:  b.EndDate,
	}
}

type BillingSelfConsumptionCurvePoint struct {
	CUPS string   `json:"CUPS" bson:"CUPS"`
	AI   float64  `json:"AI" bson:"AI"`
	AE   float64  `json:"AE" bson:"AE"`
	EHGN *float64 `json:"EHGN" bson:"EHGN"`
	EHCR *float64 `json:"EHCR" bson:"EHCR"`
	EHEX *float64 `json:"EHEX" bson:"EHEX"`
	EHAU *float64 `json:"EHAU" bson:"EHAU"`
	EHSA *float64 `json:"EHSA" bson:"EHSA"`
	EHDC *float64 `json:"EHDC" bson:"EHDC"`
}

func (b BillingSelfConsumptionCurvePoint) toDomain() billing_measures.BillingSelfConsumptionCurvePoint {
	return billing_measures.BillingSelfConsumptionCurvePoint{
		CUPS: b.CUPS,
		AI:   b.AI,
		AE:   b.AE,
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHGN: b.EHGN,
			EHCR: b.EHCR,
			EHEX: b.EHEX,
			EHAU: b.EHAU,
			EHSA: b.EHSA,
			EHDC: b.EHDC,
		},
	}
}

type BillingSelfConsumptionCurve struct {
	EndDate time.Time                          `json:"end_date" bson:"end_date"`
	Points  []BillingSelfConsumptionCurvePoint `json:"points" bson:"points"`
	EHGN    *float64                           `json:"EHGN" bson:"EHGN"`
	EHCR    *float64                           `json:"EHCR" bson:"EHCR"`
	EHEX    *float64                           `json:"EHEX" bson:"EHEX"`
	EHAU    *float64                           `json:"EHAU" bson:"EHAU"`
	EHSA    *float64                           `json:"EHSA" bson:"EHSA"`
	EHDC    *float64                           `json:"EHDC" bson:"EHDC"`
}

func (b BillingSelfConsumptionCurve) toDomain() billing_measures.BillingSelfConsumptionCurve {
	return billing_measures.BillingSelfConsumptionCurve{
		EndDate: b.EndDate,
		Points: utils.MapSlice(b.Points, func(item BillingSelfConsumptionCurvePoint) billing_measures.BillingSelfConsumptionCurvePoint {
			return item.toDomain()
		}),
		BillingSelfConsumptionValues: billing_measures.BillingSelfConsumptionValues{
			EHGN: b.EHGN,
			EHCR: b.EHCR,
			EHEX: b.EHEX,
			EHAU: b.EHAU,
			EHSA: b.EHSA,
			EHDC: b.EHDC,
		},
	}
}

type BillingSelfConsumptionConfig struct {
	Id                    string    `json:"_id" bson:"_id"`
	CauID                 string    `json:"cau_ID" bson:"cau_ID"`
	StatusID              int       `json:"status_ID" bson:"status_ID"`
	StatusName            string    `json:"status_name" bson:"status_name"`
	InitDate              time.Time `json:"init_date" bson:"init_date"`
	EndDate               time.Time `json:"end_date" bson:"end_date"`
	CnmcTypeId            int       `json:"cnmc_type_id" bson:"cnmc_type_id"`
	CnmcTypeName          string    `json:"cnmc_type_name" bson:"cnmc_type_name"`
	CnmcTypeDesc          string    `json:"cnmc_type_desc" bson:"cnmc_type_desc"`
	ConfType              string    `json:"conf_type" bson:"conf_type"`
	ConfTypeDescription   string    `json:"conf_type_description" bson:"conf_type_description"`
	ConsumerType          string    `json:"consumer_type" bson:"consumer_type"`
	ParticipantNumber     int       `json:"participant_number" bson:"participant_number"`
	ConnType              string    `json:"conn_type" bson:"conn_type"`
	Excedents             bool      `json:"excedents" bson:"excedents"`
	Compensation          bool      `json:"compensation" bson:"compensation"`
	GenerationPot         float64   `json:"generation_pot" bson:"generation_pot"`
	GroupSubgroup         int       `json:"group_subgroup" bson:"group_subgroup"`
	AntivertType          string    `json:"antivert_type" bson:"antivert_type"`
	SolarZoneId           int       `json:"solar_zone_id" bson:"solar_zone_id"`
	SolarZoneNum          int       `json:"solar_zone_num" bson:"solar_zone_num"`
	SolarZoneName         string    `json:"solar_zone_name" bson:"solar_zone_name"`
	TechnologyId          string    `json:"technology_id" bson:"technology_id"`
	TechnologyDescription string    `json:"technology_description" bson:"technology_description"`
}

func (s BillingSelfConsumptionConfig) toDomain() billing_measures.BillingSelfConsumptionConfig {
	return billing_measures.BillingSelfConsumptionConfig{
		Id:                    s.Id,
		CauID:                 s.CauID,
		StatusID:              s.StatusID,
		StatusName:            s.StatusName,
		InitDate:              s.InitDate,
		EndDate:               s.EndDate,
		CnmcTypeId:            s.CnmcTypeId,
		CnmcTypeName:          s.CnmcTypeName,
		CnmcTypeDesc:          s.CnmcTypeDesc,
		ConfType:              billing_measures.ConfigType(s.ConfType),
		ConfTypeDescription:   s.ConfTypeDescription,
		ConsumerType:          s.ConsumerType,
		ParticipantNumber:     s.ParticipantNumber,
		ConnType:              billing_measures.ConnectionType(s.ConnType),
		Excedents:             s.Excedents,
		Compensation:          s.Compensation,
		GenerationPot:         s.GenerationPot,
		GroupSubgroup:         s.GroupSubgroup,
		AntivertType:          s.AntivertType,
		SolarZoneId:           s.SolarZoneId,
		SolarZoneNum:          s.SolarZoneNum,
		SolarZoneName:         s.SolarZoneName,
		TechnologyId:          s.TechnologyId,
		TechnologyDescription: s.TechnologyDescription,
	}
}

/***************

	TO DB

****************/

func (repository BillingSelfConsumptionRepository) toDb(b billing_measures.BillingSelfConsumption) BillingSelfConsumption {
	return BillingSelfConsumption{
		Id:                b.Id,
		SelfconsumptionId: b.SelfconsumptionId,
		CAU:               b.CAU,
		Name:              b.Name,
		DistributorId:     b.DistributorId,
		InitDate:          b.InitDate,
		EndDate:           b.EndDate,
		GenerationDate:    b.GenerationDate,
		Status:            b.Status,
		Config:            repository.billingSelfConsumptionConfigToDb(b.Config),
		Points: utils.MapSlice(b.Points, func(item billing_measures.BillingSelfConsumptionPoint) BillingSelfConsumptionPoint {
			return repository.billingSelfConsumptionPointToDb(item)
		}),
		Curve: utils.MapSlice(b.Curve, func(item billing_measures.BillingSelfConsumptionCurve) BillingSelfConsumptionCurve {
			return repository.billingSelfConsumptionCurveToDb(item)
		}),
		GraphHistory: b.GraphHistory,
		SurplusType:  b.SurplusType,
	}
}

func (repository BillingSelfConsumptionRepository) billingSelfConsumptionConfigToDb(s billing_measures.BillingSelfConsumptionConfig) BillingSelfConsumptionConfig {
	return BillingSelfConsumptionConfig{
		Id:                    s.Id,
		CauID:                 s.CauID,
		StatusID:              s.StatusID,
		StatusName:            s.StatusName,
		InitDate:              s.InitDate,
		EndDate:               s.EndDate,
		CnmcTypeId:            s.CnmcTypeId,
		CnmcTypeName:          s.CnmcTypeName,
		CnmcTypeDesc:          s.CnmcTypeDesc,
		ConfType:              string(s.ConfType),
		ConfTypeDescription:   s.ConfTypeDescription,
		ConsumerType:          s.ConsumerType,
		ParticipantNumber:     s.ParticipantNumber,
		ConnType:              string(s.ConnType),
		Excedents:             s.Excedents,
		Compensation:          s.Compensation,
		GenerationPot:         s.GenerationPot,
		GroupSubgroup:         s.GroupSubgroup,
		AntivertType:          s.AntivertType,
		SolarZoneId:           s.SolarZoneId,
		SolarZoneNum:          s.SolarZoneNum,
		SolarZoneName:         s.SolarZoneName,
		TechnologyId:          s.TechnologyId,
		TechnologyDescription: s.TechnologyDescription,
	}
}

func (repository BillingSelfConsumptionRepository) billingSelfConsumptionCurveToDb(b billing_measures.BillingSelfConsumptionCurve) BillingSelfConsumptionCurve {
	return BillingSelfConsumptionCurve{
		EndDate: b.EndDate,
		Points: utils.MapSlice(b.Points, func(item billing_measures.BillingSelfConsumptionCurvePoint) BillingSelfConsumptionCurvePoint {
			return repository.billingSelfConsumptionCurvePointToDb(item)
		}),
		EHGN: b.EHGN,
		EHCR: b.EHCR,
		EHEX: b.EHEX,
		EHAU: b.EHAU,
		EHSA: b.EHSA,
		EHDC: b.EHDC,
	}
}

func (repository BillingSelfConsumptionRepository) billingSelfConsumptionCurvePointToDb(b billing_measures.BillingSelfConsumptionCurvePoint) BillingSelfConsumptionCurvePoint {
	return BillingSelfConsumptionCurvePoint{
		CUPS: b.CUPS,
		AI:   b.AI,
		AE:   b.AE,
		EHGN: b.EHGN,
		EHCR: b.EHCR,
		EHEX: b.EHEX,
		EHAU: b.EHAU,
		EHSA: b.EHSA,
		EHDC: b.EHDC,
	}
}

func (repository BillingSelfConsumptionRepository) billingSelfConsumptionPointToDb(b billing_measures.BillingSelfConsumptionPoint) BillingSelfConsumptionPoint {
	return BillingSelfConsumptionPoint{
		Id:               b.ID,
		ServicePointType: string(b.ServicePointType),
		CUPS:             b.CUPS,
		CUPSgd:           b.CUPSgd,
		InstalationFlag:  b.InstalationFlag,
		WithoutmeterFlag: b.WithoutmeterFlag,
		Exent1Flag:       b.Exent1Flag,
		Exent2Flag:       b.Exent2Flag,
		PartitionCoeff:   b.PartitionCoeff,
		MvhReceived: utils.MapSlice(b.MvhReceived, func(item billing_measures.BillingSelfConsumptionPointMvhReceived) BillingSelfConsumptionPointMvhReceived {
			return BillingSelfConsumptionPointMvhReceived{
				BillingMeasureId: item.BillingMeasureId,
				InitDate:         item.InitDate,
				EndDate:          item.EndDate,
			}
		}),
		InitDate: b.InitDate,
		EndDate:  b.EndDate,
	}
}
