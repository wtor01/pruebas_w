package billing_measures

import (
	"bitbucket.org/sercide/data-ingestion/pkg/utils"
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

type BillingSelfConsumptionStatus string

const (
	PendingPoints                    BillingSelfConsumptionStatus = "PENDING_POINTS"
	CalculatingPoints                BillingSelfConsumptionStatus = "CALCULATING"
	SupervisionSelfConsumptionStatus BillingSelfConsumptionStatus = "Supervision"
	CalculatedSelfConsumptionStatus  BillingSelfConsumptionStatus = "CALCULATED"
)

type ConfigType string

const (
	ConfigTypeA  ConfigType = "A"
	ConfigTypeB  ConfigType = "B"
	ConfigTypeC  ConfigType = "C"
	ConfigTypeD  ConfigType = "D"
	ConfigTypeE1 ConfigType = "E1"
)

type BillingSelfConsumptionConfig struct {
	Id                    string         `json:"_id"`
	CauID                 string         `json:"cau_ID"`
	StatusID              int            `json:"status_ID"`
	StatusName            string         `json:"status_name"`
	InitDate              time.Time      `json:"init_date"`
	EndDate               time.Time      `json:"end_date"`
	CnmcTypeId            int            `json:"cnmc_type_id"`
	CnmcTypeName          string         `json:"cnmc_type_name"`
	CnmcTypeDesc          string         `json:"cnmc_type_desc"`
	ConfType              ConfigType     `json:"conf_type"`
	ConfTypeDescription   string         `json:"conf_type_description"`
	ConsumerType          string         `json:"consumer_type"`
	ParticipantNumber     int            `json:"participant_number"`
	ConnType              ConnectionType `json:"conn_type"`
	Excedents             bool           `json:"excedents"`
	Compensation          bool           `json:"compensation"`
	GenerationPot         float64        `json:"generation_pot"`
	GroupSubgroup         int            `json:"group_subgroup"`
	AntivertType          string         `json:"antivert_type"`
	SolarZoneId           int            `json:"solar_zone_id"`
	SolarZoneNum          int            `json:"solar_zone_num"`
	SolarZoneName         string         `json:"solar_zone_name"`
	TechnologyId          string         `json:"technology_id"`
	TechnologyDescription string         `json:"technology_description"`
}

type BillingSelfConsumptionPointMvhReceived struct {
	BillingMeasureId string    `json:"billing_measure_id"`
	InitDate         time.Time `json:"init_date"`
	EndDate          time.Time `json:"end_date"`
}

type BillingSelfConsumptionPoint struct {
	ID               string                                   `json:"_id"`
	ServicePointType PointSelfConsumptionType                 `json:"service_point_type"`
	CUPS             string                                   `json:"CUPS"`
	CUPSgd           *string                                  `json:"CUPS_gd" bson:"CUPS_gd"`
	InstalationFlag  int                                      `json:"instalation_flag"`
	WithoutmeterFlag int                                      `json:"withoutmeter_flag"`
	Exent1Flag       int                                      `json:"exent1_flag"`
	Exent2Flag       int                                      `json:"exent2_flag"`
	PartitionCoeff   float64                                  `json:"partition_coeff"`
	MvhReceived      []BillingSelfConsumptionPointMvhReceived `json:"mvh_received"`
	InitDate         time.Time                                `json:"init_date,omitempty"`
	EndDate          time.Time                                `json:"end_date,omitempty"`
}

func NewBillingSelfConsumptionPoint(point PointSelfConsumption, billingMeasure BillingMeasure) BillingSelfConsumptionPoint {
	mvhReceived := make([]BillingSelfConsumptionPointMvhReceived, 0)

	if billingMeasure.CUPS == point.CUPS {
		mvhReceived = append(mvhReceived, BillingSelfConsumptionPointMvhReceived{
			BillingMeasureId: billingMeasure.Id,
			InitDate:         billingMeasure.InitDate,
			EndDate:          billingMeasure.EndDate,
		})
	}

	return BillingSelfConsumptionPoint{
		ID:               point.ID,
		ServicePointType: point.ServicePointType,
		CUPS:             point.CUPS,
		InstalationFlag:  point.InstalationFlag,
		WithoutmeterFlag: point.WithoutmeterFlag,
		Exent1Flag:       point.Exent1Flag,
		Exent2Flag:       point.Exent2Flag,
		PartitionCoeff:   point.PartitionCoeff,
		MvhReceived:      mvhReceived,
		InitDate:         point.InitDate,
		EndDate:          point.EndDate,
	}
}

type BillingSelfConsumptionValues struct {
	EHGN *float64 `json:"EHGN"`
	EHCR *float64 `json:"EHCR"`
	EHEX *float64 `json:"EHEX"`
	EHAU *float64 `json:"EHAU"`
	EHSA *float64 `json:"EHSA"`
	EHDC *float64 `json:"EHDC"`
}

type BillingSelfConsumptionCurvePoint struct {
	CUPS string  `json:"CUPS"`
	AI   float64 `json:"AI"`
	AE   float64 `json:"AE"`
	BillingSelfConsumptionValues
}

type BillingSelfConsumptionCurve struct {
	EndDate time.Time                          `json:"end_date"`
	Points  []BillingSelfConsumptionCurvePoint `json:"points"`
	BillingSelfConsumptionValues
}

type SurplusType string

const (
	NaSurplusType = "NA"
	DcSurplusType = "DC"
	GdSurplusType = "GD"
)

type BillingSelfConsumption struct {
	Id                string                        `json:"_id"`
	SurplusType       SurplusType                   `json:"surplus_type"`
	SelfconsumptionId string                        `json:"selfconsumption_id"`
	CAU               string                        `json:"CAU"`
	Name              string                        `json:"name"`
	DistributorId     string                        `json:"distributor_id"`
	InitDate          time.Time                     `json:"init_date"`
	EndDate           time.Time                     `json:"end_date"`
	GenerationDate    time.Time                     `json:"generation_date"`
	Status            BillingSelfConsumptionStatus  `json:"status"`
	Config            BillingSelfConsumptionConfig  `json:"config"`
	Points            []BillingSelfConsumptionPoint `json:"points"`
	Curve             []BillingSelfConsumptionCurve `json:"curve"`
	GraphHistory      *Graph                        `json:"graph_history"`
}

func NewBillingSelfConsumption(selfConsumption SelfConsumption, billingMeasure BillingMeasure) (BillingSelfConsumption, error) {
	id, err := utils.GenerateId()

	if err != nil {
		return BillingSelfConsumption{}, err
	}

	points := make([]BillingSelfConsumptionPoint, 0, cap(selfConsumption.Points))

	for _, point := range selfConsumption.Points {
		points = append(points, NewBillingSelfConsumptionPoint(point, billingMeasure))
	}

	curve := make([]BillingSelfConsumptionCurve, 0)

	for _, loadCurve := range billingMeasure.BillingLoadCurve {
		curve = append(curve, BillingSelfConsumptionCurve{
			EndDate: loadCurve.EndDate,
			Points: []BillingSelfConsumptionCurvePoint{
				{
					CUPS: billingMeasure.CUPS,
					AI:   loadCurve.AI,
					AE:   loadCurve.AE,
				},
			},
		})
	}

	config := selfConsumption.GetBillingSelfConsumptionConfig(billingMeasure)

	if config.Id == "" {
		return BillingSelfConsumption{}, errors.New("invalid billing self consumption config")
	}

	b := BillingSelfConsumption{
		Id:                id,
		SelfconsumptionId: selfConsumption.ID,
		CAU:               selfConsumption.CAU,
		Name:              selfConsumption.Name,
		DistributorId:     selfConsumption.DistributorId,
		InitDate:          billingMeasure.InitDate,
		EndDate:           billingMeasure.EndDate,
		GenerationDate:    time.Now(),
		Status:            PendingPoints,
		Config:            config,
		Points:            points,
		Curve:             curve,
	}

	return b, nil
}

func (billingSelfConsumption *BillingSelfConsumption) SetDates(billingMeasure BillingMeasure) {
	if billingMeasure.InitDate.Before(billingSelfConsumption.InitDate) {
		billingSelfConsumption.InitDate = billingMeasure.InitDate
	}

	if billingSelfConsumption.EndDate.Before(billingMeasure.EndDate) {
		billingSelfConsumption.EndDate = billingMeasure.EndDate
	}
}

func (billingSelfConsumption *BillingSelfConsumption) AddBillingMeasure(billingMeasure BillingMeasure) {

	billingSelfConsumption.SetDates(billingMeasure)

	for i, point := range billingSelfConsumption.Points {
		if point.CUPS != billingMeasure.CUPS {
			continue
		}

		added := false

		for _, mvhReceived := range point.MvhReceived {
			if mvhReceived.BillingMeasureId == billingMeasure.Id {
				added = true
			}
		}

		if added {
			break
		}

		billingSelfConsumption.Points[i].MvhReceived = append(point.MvhReceived, BillingSelfConsumptionPointMvhReceived{
			BillingMeasureId: billingMeasure.Id,
			InitDate:         billingMeasure.InitDate,
			EndDate:          billingMeasure.EndDate,
		})

	}

	curvesToAdd := make(map[time.Time]BillingLoadCurve)

	for _, loadCurve := range billingMeasure.BillingLoadCurve {
		curvesToAdd[loadCurve.EndDate] = loadCurve
	}

	for i, loadCurve := range billingSelfConsumption.Curve {
		curveToAdd, ok := curvesToAdd[loadCurve.EndDate]

		if !ok {
			continue
		}
		updated := false

		for j, point := range loadCurve.Points {
			if point.CUPS == billingMeasure.CUPS {
				billingSelfConsumption.Curve[i].Points[j].AI = curveToAdd.AI
				billingSelfConsumption.Curve[i].Points[j].AE = curveToAdd.AE
				updated = true
			}
		}
		if !updated {
			billingSelfConsumption.Curve[i].Points = append(billingSelfConsumption.Curve[i].Points, BillingSelfConsumptionCurvePoint{
				CUPS: billingMeasure.CUPS,
				AI:   curveToAdd.AI,
				AE:   curveToAdd.AE,
			})
		}
		delete(curvesToAdd, loadCurve.EndDate)
	}

	for key, curve := range curvesToAdd {
		billingSelfConsumption.Curve = append(billingSelfConsumption.Curve, BillingSelfConsumptionCurve{
			EndDate: key,
			Points: []BillingSelfConsumptionCurvePoint{
				{
					CUPS: billingMeasure.CUPS,
					AI:   curve.AI,
					AE:   curve.AE,
				},
			},
		})
	}
}

func (billingSelfConsumption BillingSelfConsumption) IsRedyToProcess() bool {

	totalHoursShouldHaveBillingMeasuresByCup := billingSelfConsumption.EndDate.Sub(billingSelfConsumption.InitDate).Hours()

	for _, point := range billingSelfConsumption.Points {
		totalHoursByCup := .0
		for _, received := range point.MvhReceived {
			totalHoursByCup += received.EndDate.Sub(received.InitDate).Hours()
		}
		if totalHoursShouldHaveBillingMeasuresByCup != totalHoursByCup {
			return false
		}
	}

	return true
}

func (billingSelfConsumption *BillingSelfConsumption) SetStatusPendingPoints() {
	billingSelfConsumption.Status = PendingPoints
}

func (billingSelfConsumption *BillingSelfConsumption) FinishProcess() {
	billingSelfConsumption.Status = CalculatedSelfConsumptionStatus
}

func (billingSelfConsumption *BillingSelfConsumption) SetStatusSupervision() {
	billingSelfConsumption.Status = SupervisionSelfConsumptionStatus
}

func (billingSelfConsumption *BillingSelfConsumption) StartProcess() {
	billingSelfConsumption.Status = CalculatingPoints
}

func (billingSelfConsumption BillingSelfConsumption) GetMapCurves() map[string]BillingSelfConsumptionCurvePoint {
	curvesToAdd := make(map[string]BillingSelfConsumptionCurvePoint)

	for _, curve := range billingSelfConsumption.Curve {
		for _, point := range curve.Points {
			key := fmt.Sprintf("%s|%s", point.CUPS, curve.EndDate.Format("2006-01-02 15-04-05"))
			curvesToAdd[key] = point
		}
	}

	return curvesToAdd
}

type BillingSelfConsumptionPointMvhReceivedWithType struct {
	BillingSelfConsumptionPointMvhReceived
	ServicePointType PointSelfConsumptionType
}

func (billingSelfConsumption BillingSelfConsumption) GetMvhReceived() []BillingSelfConsumptionPointMvhReceivedWithType {
	mvhReceived := make([]BillingSelfConsumptionPointMvhReceivedWithType, 0)

	for _, point := range billingSelfConsumption.Points {
		for _, received := range point.MvhReceived {
			mvhReceived = append(mvhReceived, BillingSelfConsumptionPointMvhReceivedWithType{
				BillingSelfConsumptionPointMvhReceived: BillingSelfConsumptionPointMvhReceived{
					BillingMeasureId: received.BillingMeasureId,
					InitDate:         received.InitDate,
					EndDate:          received.EndDate,
				},
				ServicePointType: point.ServicePointType,
			})
		}
	}

	return mvhReceived
}

func (billingSelfConsumption BillingSelfConsumption) UpdateBillingMeasureCurve(
	ctx context.Context,
	billingMeasureRepository BillingMeasureRepository,
) ([]BillingMeasure, error) {
	mvhReceived := billingSelfConsumption.GetMvhReceived()
	curvesToAdd := billingSelfConsumption.GetMapCurves()

	g, ctx := errgroup.WithContext(ctx)

	billingMeasuresToUpdate := make([]BillingMeasure, len(mvhReceived))

	for indexMvhReceived, mvh := range mvhReceived {

		g.Go(func(indexMvhReceived int, mvh BillingSelfConsumptionPointMvhReceivedWithType) func() error {

			return func() error {

				billingMeasure, errGo := billingMeasureRepository.Find(ctx, mvh.BillingMeasureId)

				if errGo != nil {
					return errGo
				}

				for i, curve := range billingMeasure.BillingLoadCurve {
					key := fmt.Sprintf("%s|%s", billingMeasure.CUPS, curve.EndDate.Format("2006-01-02 15-04-05"))
					c, ok := curvesToAdd[key]
					if !ok {
						continue
					}

					switch mvh.ServicePointType {
					case SsaaServicePointType:
						billingMeasure.SetLoadCurveSelfConsumptionMagnitude(i, AIAutoSelfConsumptionMagnitude, c.EHSA)
					case ConsumoServicePointType, FronteraServicePointType:
						{
							billingMeasure.SumLoadCurveSelfConsumptionMagnitude(i, AIAutoSelfConsumptionMagnitude, c.EHCR)
							if billingSelfConsumption.SurplusType == DcSurplusType {
								billingMeasure.SumLoadCurveSelfConsumptionMagnitude(i, AeAutoSelfConsumptionMagnitude, c.EHEX)
								billingMeasure.SetLoadCurveSelfConsumptionMagnitude(i, EHEXSelfConsumptionMagnitude, c.EHEX)
							} else {
								val := .0
								billingMeasure.SetLoadCurveSelfConsumptionMagnitude(i, AeAutoSelfConsumptionMagnitude, &val)
							}
							billingMeasure.SetLoadCurveSelfConsumptionMagnitude(i, EHAUSelfConsumptionMagnitude, c.EHAU)
						}
					case GdServicePointType:
						{
							if billingSelfConsumption.SurplusType == GdSurplusType {
								billingMeasure.SumLoadCurveSelfConsumptionMagnitude(i, AeAutoSelfConsumptionMagnitude, c.EHEX)
								billingMeasure.SetLoadCurveSelfConsumptionMagnitude(i, EHEXSelfConsumptionMagnitude, c.EHEX)
							} else {
								val := .0
								billingMeasure.SetLoadCurveSelfConsumptionMagnitude(i, AeAutoSelfConsumptionMagnitude, &val)
								billingMeasure.SetLoadCurveSelfConsumptionMagnitude(i, EHEXSelfConsumptionMagnitude, &val)
							}
							billingMeasure.SetLoadCurveSelfConsumptionMagnitude(i, EHGNSelfConsumptionMagnitude, c.EHGN)

						}
					}
				}

				billingMeasure.Calculated()
				billingMeasuresToUpdate[indexMvhReceived] = billingMeasure

				return nil
			}
		}(indexMvhReceived, mvh))
	}

	err := g.Wait()

	return billingMeasuresToUpdate, err
}

func (b *BillingSelfConsumption) BeforeSave() {
	if b.GraphHistory == nil {
		return
	}

	for k, node := range b.GraphHistory.Dict {
		if !node.Done {
			delete(b.GraphHistory.Dict, k)
		}
	}
}

type QueryGetBillingSelfConsumption struct {
	SelfconsumptionId string
	StartDate         time.Time
	EndDate           time.Time
}

type QueryGetBillingSelfConsumptionByCau struct {
	CauId         string
	DistributorId string
	StartDate     time.Time
	EndDate       time.Time
}

//go:generate mockery --case=snake --outpkg=mocks --output=../platform/mocks --name=BillingSelfConsumptionRepository
type BillingSelfConsumptionRepository interface {
	GetBySelfConsumptionBetweenDates(ctx context.Context, query QueryGetBillingSelfConsumption) (BillingSelfConsumption, error)
	Save(ctx context.Context, b BillingSelfConsumption) error
	GetSelfConsumptionByCau(ctx context.Context, query QueryGetBillingSelfConsumptionByCau) ([]BillingSelfConsumption, error)
}
