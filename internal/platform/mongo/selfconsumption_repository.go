package mongo

import (
	"bitbucket.org/sercide/data-ingestion/internal/billing_measures"
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

func NewSelfConsumptionRepository(client *mongo.Client, database string, loc *time.Location) *SelfConsumptionRepository {
	return &SelfConsumptionRepository{
		client:     client,
		location:   loc,
		Database:   database,
		Collection: "view_selfconsumption_joined",
	}
}

type PointSelfConsumption struct {
	ID               string               `bson:"_id"`
	Id               string               `bson:"id"`
	CauID            string               `bson:"cau_ID"`
	ServicePointType string               `bson:"service_point_type"`
	CUPS             string               `bson:"CUPS"`
	InitDate         time.Time            `bson:"init_date"`
	EndDate          time.Time            `bson:"end_date"`
	InstalationFlag  int                  `bson:"instalation_flag"`
	WithoutmeterFlag int                  `bson:"withoutmeter_flag"`
	Exent1Flag       int                  `bson:"exent1_flag"`
	Exent2Flag       int                  `bson:"exent2_flag"`
	PartitionCoeff   primitive.Decimal128 `bson:"partition_coeff"`
	DistributorId    string               `bson:"distributor_id"`
}

func (p PointSelfConsumption) toDomain() (billing_measures.PointSelfConsumption, error) {
	partitionCoeff, err := strconv.ParseFloat(p.PartitionCoeff.String(), 64)

	if err != nil {
		return billing_measures.PointSelfConsumption{}, err
	}

	return billing_measures.PointSelfConsumption{
		ID:               p.ID,
		CauID:            p.CauID,
		ServicePointType: billing_measures.PointSelfConsumptionType(p.ServicePointType),
		CUPS:             p.CUPS,
		InitDate:         p.InitDate,
		EndDate:          p.EndDate,
		InstalationFlag:  p.InstalationFlag,
		WithoutmeterFlag: p.WithoutmeterFlag,
		Exent1Flag:       p.Exent1Flag,
		Exent2Flag:       p.Exent2Flag,
		PartitionCoeff:   partitionCoeff,
		DistributorId:    p.DistributorId,
	}, nil
}

type ConfigSelfConsumption struct {
	ID                    string               `bson:"_id"`
	Id                    string               `bson:"id"`
	CauID                 string               `bson:"cau_ID"`
	StatusID              int                  `bson:"status_ID"`
	StatusName            string               `bson:"status_name"`
	InitDate              time.Time            `bson:"init_date"`
	EndDate               time.Time            `bson:"end_date"`
	CreationDate          time.Time            `bson:"creation_date"`
	DistributorId         string               `bson:"distributor_id"`
	CnmcTypeId            int                  `bson:"cnmc_type_id"`
	CnmcTypeName          string               `bson:"cnmc_type_name"`
	CnmcTypeDesc          string               `bson:"cnmc_type_desc"`
	ConfType              string               `bson:"conf_type"`
	ConfTypeDescription   string               `bson:"conf_type_description"`
	ConsumerType          string               `bson:"consumer_type"`
	ParticipantNumber     int                  `bson:"participant_number"`
	ConnType              string               `bson:"conn_type"`
	Excedents             bool                 `bson:"excedents"`
	Compensation          bool                 `bson:"compensation"`
	GenerationPot         primitive.Decimal128 `bson:"generation_pot"`
	GroupSubgroup         int                  `bson:"group_subgroup"`
	AntivertType          string               `bson:"antivert_type"`
	SolarZoneId           int                  `bson:"solar_zone_id"`
	SolarZoneNum          int                  `bson:"solar_zone_num"`
	SolarZoneName         string               `bson:"solar_zone_name"`
	TechnologyId          string               `bson:"technology_id"`
	TechnologyDescription string               `bson:"technology_description"`
}

func (c ConfigSelfConsumption) toDomain() (billing_measures.ConfigSelfConsumption, error) {
	generationPot, err := strconv.ParseFloat(c.GenerationPot.String(), 64)

	if err != nil {
		return billing_measures.ConfigSelfConsumption{}, err
	}

	return billing_measures.ConfigSelfConsumption{
		ID:                    c.ID,
		CauID:                 c.CauID,
		StatusID:              c.StatusID,
		StatusName:            c.StatusName,
		InitDate:              c.InitDate,
		EndDate:               c.EndDate,
		CreationDate:          c.CreationDate,
		DistributorId:         c.DistributorId,
		CnmcTypeId:            c.CnmcTypeId,
		CnmcTypeName:          c.CnmcTypeName,
		CnmcTypeDesc:          c.CnmcTypeDesc,
		ConfType:              billing_measures.ConfigType(c.ConfType),
		ConfTypeDescription:   c.ConfTypeDescription,
		ConsumerType:          c.ConsumerType,
		ParticipantNumber:     c.ParticipantNumber,
		ConnType:              billing_measures.ConnectionType(c.ConnType),
		Excedents:             c.Excedents,
		Compensation:          c.Compensation,
		GenerationPot:         generationPot,
		GroupSubgroup:         c.GroupSubgroup,
		AntivertType:          c.AntivertType,
		SolarZoneId:           c.SolarZoneId,
		SolarZoneNum:          c.SolarZoneNum,
		SolarZoneName:         c.SolarZoneName,
		TechnologyId:          c.TechnologyId,
		TechnologyDescription: c.TechnologyDescription,
	}, nil
}

type SelfConsumption struct {
	ID            string                  `bson:"_id"`
	Id            string                  `bson:"id"`
	CAU           string                  `bson:"CAU"`
	Name          string                  `bson:"name"`
	StatusID      int                     `bson:"status_ID"`
	StatusName    string                  `bson:"status_Name"`
	CcaaId        int                     `bson:"ccaaId"`
	Ccaa          string                  `bson:"ccaa"`
	InitDate      time.Time               `bson:"init_date"`
	EndDate       time.Time               `bson:"end_date"`
	DistributorId string                  `bson:"distributor_id"`
	Configs       []ConfigSelfConsumption `bson:"configs"`
	Points        []PointSelfConsumption  `bson:"points"`
}

func (s SelfConsumption) toDomain() (billing_measures.SelfConsumption, error) {
	configs := make([]billing_measures.ConfigSelfConsumption, 0, cap(s.Configs))
	points := make([]billing_measures.PointSelfConsumption, 0, cap(s.Points))

	for _, config := range s.Configs {
		c, err := config.toDomain()
		if err != nil {
			return billing_measures.SelfConsumption{}, err
		}
		configs = append(configs, c)
	}
	for _, point := range s.Points {
		p, err := point.toDomain()
		if err != nil {
			return billing_measures.SelfConsumption{}, err
		}
		points = append(points, p)
	}

	return billing_measures.SelfConsumption{
		ID:            s.ID,
		CAU:           s.CAU,
		Name:          s.Name,
		StatusID:      s.StatusID,
		StatusName:    s.StatusName,
		CcaaId:        s.CcaaId,
		Ccaa:          s.Ccaa,
		InitDate:      s.InitDate,
		EndDate:       s.EndDate,
		DistributorId: s.DistributorId,
		Configs:       configs,
		Points:        points,
	}, nil
}

type SelfConsumptionRepository struct {
	client     *mongo.Client
	location   *time.Location
	Database   string
	Collection string
}

func (repository SelfConsumptionRepository) GetActiveSelfConsumptionByCUP(ctx context.Context, query billing_measures.QueryGetActiveSelfConsumptionByCUP) (billing_measures.SelfConsumption, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.Collection)

	filter := bson.D{
		{"points", bson.D{{"$elemMatch", bson.D{{"CUPS", query.CUP}}}}},
		{"init_date", bson.D{{"$lte", query.EndDate}}},
		{"end_date", bson.D{{"$gte", query.StartDate}}},
	}
	var selfConsumption SelfConsumption

	result := coll.FindOne(ctxTimeout, filter)

	if result.Err() != nil {
		return billing_measures.SelfConsumption{}, result.Err()
	}

	err := result.Decode(&selfConsumption)

	if err != nil {
		return billing_measures.SelfConsumption{}, err
	}

	selfConsumptionDomain, err := selfConsumption.toDomain()

	return selfConsumptionDomain, err
}

// GetSelfConsumptionByCUP Function to obtain the self-consumption units to which a cup belongs for a given date
func (repository SelfConsumptionRepository) GetSelfConsumptionByCUP(ctx context.Context, query billing_measures.GetSelfConsumptionByCUP) (billing_measures.SelfConsumption, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	coll := repository.client.Database(repository.Database).Collection(repository.Collection)

	filter := bson.D{
		{"distributor_id", query.DistributorId},
		{"points", bson.D{{"$elemMatch", bson.D{{"CUPS", query.CUP}}}}},
		{"init_date", bson.D{{"$lte", query.Date}}},
		{"end_date", bson.D{{"$gte", query.Date}}},
	}
	var selfConsumption SelfConsumption

	result := coll.FindOne(ctxTimeout, filter)

	if result.Err() != nil {
		return billing_measures.SelfConsumption{}, result.Err()
	}

	err := result.Decode(&selfConsumption)

	if err != nil {
		return billing_measures.SelfConsumption{}, err
	}

	selfConsumptionDomain, err := selfConsumption.toDomain()

	return selfConsumptionDomain, err
}

// GetSelfConsumptionActiveByDistributor Function to list the self-consumption units of a distributor active on a certain date
func (repository SelfConsumptionRepository) GetSelfConsumptionActiveByDistributor(ctx context.Context, query billing_measures.GetSelfConsumptionByDistributortDto) ([]billing_measures.SelfConsumption, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	var selfConsumption []SelfConsumption
	opts := options.Find()
	opts.SetLimit(int64(query.Limit)).SetSkip(int64(query.Offset))

	match := bson.D{
		{"distributor_id", query.DistributorId},
		{"init_date", bson.D{{"$lte", query.Date}}},
		{"end_date", bson.D{{"$gte", query.Date}}},
	}

	collection := repository.client.Database(repository.Database).Collection(repository.Collection)

	coll, err := collection.Find(ctxTimeout, match, opts)

	if err != nil {
		return []billing_measures.SelfConsumption{}, 0, err
	}

	if err = coll.All(ctx, &selfConsumption); err != nil {
		return []billing_measures.SelfConsumption{}, 0, err
	}

	results := utils.MapSlice(selfConsumption, func(item SelfConsumption) billing_measures.SelfConsumption {
		i, _ := item.toDomain()
		return i
	})

	count, err := collection.CountDocuments(ctxTimeout, match)

	if err != nil {
		return results, 0, err
	}

	return results, int(count), nil
}
