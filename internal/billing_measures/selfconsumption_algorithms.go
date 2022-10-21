package billing_measures

import (
	"context"
)

func getCupsType(b BillingSelfConsumption) (map[string]PointSelfConsumptionType, map[string]string, int, int, map[string]string) {

	pointsTypes := map[string]PointSelfConsumptionType{}
	gdConfig := map[string]string{}
	ssaaConfig := map[string]string{}
	genTypesCount := 0
	conTypes := 0
	for _, p := range b.Points {
		pointsTypes[p.CUPS] = p.ServicePointType
		if p.ServicePointType == GdServicePointType {
			gdConfig[p.CUPS] = string(GdServicePointType)
			genTypesCount += 1
		}
		if p.ServicePointType == ConsumoServicePointType {
			conTypes += 1

		}
		if p.ServicePointType == SsaaServicePointType && p.CUPSgd != nil {
			ssaaConfig[p.CUPS] = *p.CUPSgd
		}
	}
	return pointsTypes, ssaaConfig, genTypesCount, conTypes, gdConfig
}

type SelfConsumptionConfA struct {
	Name string
	B    *BillingSelfConsumption `bson:"-"`
}

func NewSelfConsumptionConfA(b *BillingSelfConsumption) *SelfConsumptionConfA {
	return &SelfConsumptionConfA{
		Name: "CONFIGURATION_A",
		B:    b,
	}
}

func (algorithm SelfConsumptionConfA) ID() string {
	return algorithm.Name
}

func (algorithm SelfConsumptionConfA) Execute(_ context.Context) error {
	for i, _ := range algorithm.B.Curve {
		sumEHCR := 0.0
		sumEHEX := 0.0
		for c, _ := range algorithm.B.Curve[i].Points {
			point := &algorithm.B.Curve[i].Points[c]
			curve := &algorithm.B.Curve[i]
			vEHCR := point.AI - point.AE
			vEHEX := point.AE - point.AI

			if vEHCR < 0 {
				vEHCR = 0
			}
			if vEHEX < 0 {
				vEHEX = 0
			}
			point.EHCR = &vEHCR
			point.EHEX = &vEHEX

			sumEHCR = *point.EHCR + sumEHCR
			sumEHEX = *point.EHEX + sumEHEX

			curve.EHEX = &sumEHEX
			curve.EHCR = &sumEHCR
		}
	}
	return nil
}

type SelfConsumptionConfB struct {
	Name string
	B    *BillingSelfConsumption `bson:"-"`
}

func NewSelfConsumptionConfB(b *BillingSelfConsumption) *SelfConsumptionConfCAndB {
	return &SelfConsumptionConfCAndB{
		Name: "CONFIGURATION_B",
		B:    b,
	}
}
func (algorithm SelfConsumptionConfB) ID() string {
	return algorithm.Name
}

func (algorithm SelfConsumptionConfB) Execute(_ context.Context) error {
	pointsTypes, _, _, _, gdConfig := getCupsType(*algorithm.B)

	for i, _ := range algorithm.B.Curve {
		sumEHGN := 0.0
		sumEHSA := 0.0
		sumEHDC := 0.0
		sumEHEX := 0.0
		sumEHAU := 0.0
		sumEHCR := 0.0
		for c, _ := range algorithm.B.Curve[i].Points {
			point := &algorithm.B.Curve[i].Points[c]
			curve := &algorithm.B.Curve[i]
			typePoint := pointsTypes[point.CUPS]

			algorithm.setEHGN(&sumEHGN, point, curve, typePoint)
			algorithm.setEHCR(&sumEHCR, point, curve)
			algorithm.setEHEX(&sumEHGN, &sumEHEX, gdConfig, typePoint, point, curve)
			algorithm.setEHAU(&sumEHAU, &sumEHGN, &sumEHEX, gdConfig, typePoint, point, curve)
			algorithm.setEHSA(&sumEHSA, typePoint, point, curve)
			algorithm.setEHDC(&sumEHEX, &sumEHDC, &sumEHAU, curve)

		}
	}
	return nil
}

func (algorithm SelfConsumptionConfB) setEHGN(sumEHGN *float64, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve, typePoint PointSelfConsumptionType) {
	if typePoint == GdServicePointType {
		vEHGN := point.AE - point.AI
		if vEHGN < 0 {
			vEHGN = 0
		}
		point.EHGN = &vEHGN
	}
	sum := *point.EHGN + *sumEHGN
	sumEHGN = &sum
	curve.EHGN = sumEHGN
}
func (algorithm SelfConsumptionConfB) setEHCR(sumEHCR *float64, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	vEHCR := point.AI - point.AE
	if vEHCR < 0 {
		vEHCR = 0
	}
	//point.EHCR = &vEHCR
	//TODO:Se supone que va a ser asignado a la curva, se utiliza en el seteo de EHDC?
	//TODO:No se asigna al punto?
	sumEHCR = &vEHCR
	curve.EHCR = sumEHCR
}

func (algorithm SelfConsumptionConfB) setEHEX(sumEHGN, sumEHEX *float64, gdConfig map[string]string, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	vEHEX := point.AE - point.AI
	if vEHEX < 0 {
		vEHEX = 0
	}
	//TODO:NOS QUEDAREMOS SOLO CON LA ULTIMA?
	sumEHEX = &vEHEX
	curve.EHEX = sumEHEX
	if typePoint == GdServicePointType {
		countGd := 0
		for _, cp := range curve.Points {
			if _, ok := gdConfig[cp.CUPS]; !ok {
				countGd += 1
			}
		}
		porc := 0.0
		if sumEHGN != nil && *sumEHGN == 0 {
			porc = float64(1 / countGd)
		} else {
			porc = *sumEHGN / *point.EHGN
		}
		valEHEX := porc * *point.EHEX
		point.EHEX = &valEHEX
	}
}

func (algorithm SelfConsumptionConfB) setEHAU(sumEHAU, sumEHGN, sumEHEX *float64, gdConfig map[string]string, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	val := *sumEHGN
	if *sumEHEX > 0 {
		val -= *sumEHEX
	}
	sumEHAU = &val
	curve.EHAU = sumEHAU

	if typePoint == GdServicePointType {
		countGd := 0
		for _, cp := range curve.Points {
			if _, ok := gdConfig[cp.CUPS]; !ok {
				countGd += 1
			}
		}
		porc := 0.0
		if *sumEHGN == 0 {
			porc = float64(1 / countGd)
		} else {
			porc = *sumEHGN / *point.EHGN
		}
		vEHAU := porc * *point.EHAU
		point.EHAU = &vEHAU
	}
}

func (algorithm SelfConsumptionConfB) setEHSA(sumEHSA *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == GdServicePointType {
		vEHSA := point.AI - point.AE
		if vEHSA < 0 {
			vEHSA = 0
		}
		point.EHSA = &vEHSA
		sum := *point.EHSA + *sumEHSA
		sumEHSA = &sum
		curve.EHSA = sumEHSA
	}
}

func (algorithm SelfConsumptionConfB) setEHDC(sumEHEX, sumEHDC, sumEHAU *float64, curve *BillingSelfConsumptionCurve) {
	if *sumEHEX > 0 {
		val := *sumEHAU
		sumEHDC = &val
		curve.EHDC = sumEHDC
	} else {
		sum := *curve.EHCR + *curve.EHAU
		sumEHDC = &sum
		curve.EHDC = sumEHDC
		//TODO:PARA HACER ESTO TIENE QUE ESTAR EL VALOR DE EHCR EN LA CURVA NO?
	}
}

type SelfConsumptionConfCAndB struct {
	Name                   string
	CoeffiecientRepository ConsumCoefficientRepository
	B                      *BillingSelfConsumption `bson:"-"`
}

func NewSelfConsumptionConfCAndB(b *BillingSelfConsumption, cr ConsumCoefficientRepository) *SelfConsumptionConfCAndB {
	return &SelfConsumptionConfCAndB{
		CoeffiecientRepository: cr,
		Name:                   "CONFIGURATION_C_B",
		B:                      b,
	}
}

func (algorithm SelfConsumptionConfCAndB) ID() string {
	return algorithm.Name
}

func (algorithm SelfConsumptionConfCAndB) setEHGN(ctx context.Context, sumEHGN *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	vEHGN := 0.0
	if typePoint == GdServicePointType {
		vEHGN = point.AE - point.AI
		if vEHGN < 0 {
			vEHGN = 0
		}
		sum := vEHGN + *sumEHGN
		sumEHGN = &sum
		curve.EHGN = sumEHGN
	}
	if typePoint == ConsumoServicePointType {
		//TODO: servicio temporal
		coeff, _ := algorithm.CoeffiecientRepository.Search(ctx, QueryConsumCoefficient{})
		vEHGN = point.AE * coeff
	}

	point.EHGN = &vEHGN
}

func (algorithm SelfConsumptionConfCAndB) setEHSA(_ context.Context, sumEHSA *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == GdServicePointType {
		vEHSA := point.AI - point.AE
		if vEHSA < 0 {
			vEHSA = 0
		}
		point.EHSA = &vEHSA
		sum := *point.EHSA + *sumEHSA
		sumEHSA = &sum
		curve.EHSA = sumEHSA
	}

}
func (algorithm SelfConsumptionConfCAndB) setEHDC(_ context.Context, sumEHDC *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == ConsumoServicePointType {
		point.EHDC = &point.AI
		sum := *point.EHDC + *sumEHDC
		sumEHDC = &sum
		curve.EHSA = sumEHDC
	}
}

func (algorithm SelfConsumptionConfCAndB) setEHEX(_ context.Context, sumEHEX *float64, sumEHGN *float64, genTypesCount int, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == ConsumoServicePointType {
		if point.EHGN != nil && point.EHDC != nil {
			vEHEX := *point.EHGN - *point.EHDC
			if vEHEX < 0 {
				vEHEX = 0
			}
			point.EHEX = &vEHEX
		}
		sum := *point.EHEX + *sumEHEX
		sumEHEX = &sum
		curve.EHEX = sumEHEX
	}
	if typePoint == GdServicePointType {
		porc := 0.0
		if *sumEHGN == 0 {
			porc = float64(1 / genTypesCount)
		} else {
			porc = *sumEHGN / *point.EHGN
		}
		vEHEX := porc * *sumEHEX
		point.EHEX = &vEHEX
	}
}

func (algorithm SelfConsumptionConfCAndB) setEHAU(_ context.Context, sumEHAU *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == ConsumoServicePointType && point.EHGN != nil && point.EHEX != nil {
		vEHAU := *point.EHGN - *point.EHEX
		point.EHAU = &vEHAU
		sum := *point.EHAU + *sumEHAU
		sumEHAU = &sum
		curve.EHEX = sumEHAU
	}
}
func (algorithm SelfConsumptionConfCAndB) setEHCR(_ context.Context, sumEHCR *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == ConsumoServicePointType && point.EHGN != nil && point.EHDC != nil && point.EHAU != nil {
		valueCond := *point.EHGN - *point.EHDC
		vEHCR := 0.0
		if valueCond >= 0 {
			vEHCR = 0.0
		} else {
			vEHCR = *point.EHDC - *point.EHAU
		}
		point.EHCR = &vEHCR
		sum := *point.EHCR + *sumEHCR
		sumEHCR = &sum
		curve.EHEX = sumEHCR
	}
}

func (algorithm SelfConsumptionConfCAndB) Execute(ctx context.Context) error {

	pointsTypes, _, genTypesCount, _, _ := getCupsType(*algorithm.B)
	for i, _ := range algorithm.B.Curve {
		sumEHGN := 0.0
		sumEHSA := 0.0
		sumEHDC := 0.0
		sumEHEX := 0.0
		sumEHAU := 0.0
		sumEHCR := 0.0
		for c, _ := range algorithm.B.Curve[i].Points {

			point := &algorithm.B.Curve[i].Points[c]
			curve := &algorithm.B.Curve[i]
			typePoint := pointsTypes[point.CUPS]

			algorithm.setEHGN(ctx, &sumEHGN, typePoint, point, curve)
			algorithm.setEHSA(ctx, &sumEHSA, typePoint, point, curve)
			algorithm.setEHDC(ctx, &sumEHDC, typePoint, point, curve)
			algorithm.setEHEX(ctx, &sumEHEX, &sumEHGN, genTypesCount, typePoint, point, curve)
			algorithm.setEHAU(ctx, &sumEHAU, typePoint, point, curve)
			algorithm.setEHCR(ctx, &sumEHCR, typePoint, point, curve)
		}
	}
	return nil
}

type SelfConsumptionConfDandE1 struct {
	Name                   string
	CoeffiecientRepository ConsumCoefficientRepository
	B                      *BillingSelfConsumption `bson:"-"`
}

func NewSelfConsumptionConfDandE1(b *BillingSelfConsumption, cr ConsumCoefficientRepository) *SelfConsumptionConfDandE1 {
	return &SelfConsumptionConfDandE1{
		CoeffiecientRepository: cr,
		Name:                   "CONFIGURATION_D_E1",
		B:                      b,
	}
}

func (algorithm SelfConsumptionConfDandE1) ID() string {
	return algorithm.Name
}

func (algorithm SelfConsumptionConfDandE1) setEHGN(ctx context.Context, sumEHGN *float64, AIssaa *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	vEHGN := 0.0

	if typePoint == GdServicePointType && AIssaa != nil {
		vEHGN = point.AE - *AIssaa
		if vEHGN < 0 {
			vEHGN = 0
		}
		point.EHGN = &vEHGN
		sum := vEHGN + *sumEHGN
		sumEHGN = &sum
		curve.EHGN = sumEHGN

	}
	if typePoint == ConsumoServicePointType {
		//TODO: servicio temporal
		coeff, _ := algorithm.CoeffiecientRepository.Search(ctx, QueryConsumCoefficient{})
		vEHGN = point.AE * coeff
	}
	if vEHGN < 0 {
		vEHGN = 0
	}
	point.EHGN = &vEHGN
}

func (algorithm SelfConsumptionConfDandE1) setEHSA(_ context.Context, sumEHSA *float64, AIssaa *float64, sumAEgd *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == GdServicePointType && AIssaa != nil {
		vEHSA := *AIssaa - *sumAEgd
		if vEHSA < 0 {
			vEHSA = 0
		}
		point.EHSA = &vEHSA
	}

	if point.EHSA != nil {
		sum := *point.EHSA + *sumEHSA
		sumEHSA = &sum
		curve.EHSA = sumEHSA
	}
}

func (algorithm SelfConsumptionConfDandE1) setEHDC(_ context.Context, sumEHDC *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == GdServicePointType {
		vEHSA := point.AI
		if vEHSA < 0 {
			vEHSA = 0
		}
		point.EHSA = &vEHSA
	}

	if point.EHSA != nil {
		sum := *point.EHSA + *sumEHDC
		sumEHDC = &sum
		curve.EHSA = sumEHDC
	}
}

func (algorithm SelfConsumptionConfDandE1) setEHEX(_ context.Context, sumEHEX *float64, sumEHGN *float64, genTypesCount int, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == ConsumoServicePointType {
		if point.EHGN != nil && point.EHDC != nil {
			vEHEX := *point.EHGN - *point.EHDC
			if vEHEX < 0 {
				vEHEX = 0
			}
			point.EHEX = &vEHEX
			sum := *point.EHEX + *sumEHEX
			sumEHEX = &sum
			curve.EHEX = sumEHEX
		}

	}
	if typePoint == GdServicePointType {
		porc := 0.0
		if *sumEHGN == 0 {
			porc = float64(1 / genTypesCount)
		} else {
			if point.EHGN != nil && *point.EHGN != 0 {
				porc = *sumEHGN / *point.EHGN
			}
		}
		vEHEX := porc * *sumEHEX
		point.EHEX = &vEHEX
	}
}
func (algorithm SelfConsumptionConfDandE1) setEHAU(_ context.Context, sumEHAU *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == ConsumoServicePointType && point.EHGN != nil && point.EHEX != nil {
		vEHAU := *point.EHGN - *point.EHEX
		point.EHAU = &vEHAU
		sum := *point.EHAU + *sumEHAU
		sumEHAU = &sum
		curve.EHEX = sumEHAU
	}
}

func (algorithm SelfConsumptionConfDandE1) setEHCR(_ context.Context, sumEHCR *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == ConsumoServicePointType && point.EHGN != nil && point.EHDC != nil && point.EHAU != nil {
		valueCond := *point.EHGN - *point.EHDC
		vEHCR := 0.0
		if valueCond >= 0 {
			vEHCR = 0.0
		} else {
			vEHCR = *point.EHDC - *point.EHAU
		}
		point.EHCR = &vEHCR
		sum := *point.EHCR + *sumEHCR
		sumEHCR = &sum
		curve.EHEX = sumEHCR
	}
}

func (algorithm SelfConsumptionConfDandE1) Execute(ctx context.Context) error {
	pointsTypes, ssaaConfig, genTypesCount, _, _ := getCupsType(*algorithm.B)
	for i, _ := range algorithm.B.Curve {
		sumEHGN := 0.0
		sumEHSA := 0.0
		sumEHDC := 0.0
		sumEHEX := 0.0
		sumEHAU := 0.0
		sumEHCR := 0.0
		sumAEgd := 0.0
		for c, _ := range algorithm.B.Curve[i].Points {

			point := &algorithm.B.Curve[i].Points[c]
			curve := &algorithm.B.Curve[i]
			typePoint := pointsTypes[point.CUPS]

			vAIssaa := new(float64)
			cupsSsaa, ssaaCupsFound := ssaaConfig[point.CUPS]
			if typePoint == GdServicePointType && ssaaCupsFound {
				for _, curvePoint := range curve.Points {
					if curvePoint.CUPS == cupsSsaa {
						vAIssaa = &curvePoint.AI
					}
				}
			}

			algorithm.setEHGN(ctx, &sumEHGN, vAIssaa, typePoint, point, curve)
			algorithm.setEHSA(ctx, &sumEHSA, vAIssaa, &sumAEgd, typePoint, point, curve)
			algorithm.setEHDC(ctx, &sumEHDC, typePoint, point, curve)
			algorithm.setEHEX(ctx, &sumEHEX, &sumEHGN, genTypesCount, typePoint, point, curve)
			algorithm.setEHAU(ctx, &sumEHAU, typePoint, point, curve)
			algorithm.setEHCR(ctx, &sumEHCR, typePoint, point, curve)
		}
	}
	return nil
}

type SelfConsumptionConfE1 struct {
	Name string
	B    *BillingSelfConsumption `bson:"-"`
}

func NewSelfConsumptionConfE1(b *BillingSelfConsumption, cr ConsumCoefficientRepository) *SelfConsumptionConfE1 {
	return &SelfConsumptionConfE1{
		Name: "CONFIGURATION_E1",
		B:    b,
	}
}

func (algorithm SelfConsumptionConfE1) ID() string {
	return algorithm.Name
}

func (algorithm SelfConsumptionConfE1) setEHGN(_ context.Context, sumEHGN *float64, AIssaa *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	vEHGN := 0.0
	if typePoint == GdServicePointType && AIssaa != nil {
		vEHGN = point.AE - *AIssaa
		if vEHGN < 0 {
			vEHGN = 0
		}
		point.EHGN = &vEHGN

		sum := vEHGN + *sumEHGN
		sumEHGN = &sum
		curve.EHGN = sumEHGN
	}
}

func (algorithm SelfConsumptionConfE1) setEHCR(_ context.Context, sumEHCR *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint) {
	if typePoint == ConsumoServicePointType {
		vEHCR := point.AI - point.AE
		point.EHCR = &vEHCR
		sumEHCR = point.EHCR
	}
}
func (algorithm SelfConsumptionConfE1) setEHEX(_ context.Context, sumEHEX *float64, sumEHGN *float64, genTypesCount int, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == FronteraServicePointType {
		vEHEX := point.AE - point.AI
		if vEHEX < 0 {
			vEHEX = 0
		}
		sumEHEX = &vEHEX
		curve.EHEX = sumEHEX
		point.EHEX = &vEHEX
	}
	if typePoint == GdServicePointType {
		porc := 0.0
		if *sumEHGN == 0 {
			porc = float64(1 / genTypesCount)
		} else {
			if point.EHGN != nil {
				porc = *sumEHGN / *point.EHGN
			}
		}
		vEHEX := porc * *sumEHEX
		point.EHEX = &vEHEX
	}
}
func (algorithm SelfConsumptionConfE1) setEHAU(_ context.Context, sumEHAU *float64, sumEHEX *float64, sumEHGN *float64, gdConfig map[string]string, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == FronteraServicePointType {
		vEHAU := 0.0
		if *sumEHEX > 0 {
			vEHAU = *sumEHGN - *sumEHEX
		} else {
			vEHAU = *sumEHGN
		}
		point.EHAU = &vEHAU
		sumEHAU = point.EHAU
		curve.EHAU = sumEHAU
	}

	if typePoint == GdServicePointType {
		porc := 0.0
		countGd := 0
		for _, cp := range curve.Points {
			if _, ok := gdConfig[cp.CUPS]; !ok {
				countGd += 1
			}
		}
		if *sumEHGN == 0 && countGd != 0 {
			porc = float64(1 / countGd)
		} else {
			porc = *sumEHGN / *point.EHGN
		}

		vEHAU := porc * *sumEHEX
		point.EHAU = &vEHAU
	}
}

func (algorithm SelfConsumptionConfE1) setEHSA(_ context.Context, sumEHSA *float64, AIssaa *float64, sumAEgd *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint, curve *BillingSelfConsumptionCurve) {
	if typePoint == GdServicePointType {
		if typePoint == GdServicePointType && AIssaa != nil {
			vEHSA := *AIssaa - *sumAEgd
			if vEHSA < 0 {
				vEHSA = 0
			}
			point.EHSA = &vEHSA
			sum := *point.EHSA + *sumEHSA
			sumEHSA = &sum
			curve.EHSA = sumEHSA
		}
	}
}
func (algorithm SelfConsumptionConfE1) setEHDC(_ context.Context, sumEHAU *float64, sumEHEX *float64, sumEHCR *float64, typePoint PointSelfConsumptionType, point *BillingSelfConsumptionCurvePoint) {
	if typePoint == FronteraServicePointType {
		if *sumEHEX > 0 {
			point.EHDC = sumEHAU
		} else {
			v := *sumEHCR + *sumEHAU
			point.EHDC = &v
		}
		sumEHCR = point.EHCR
	}
}

func (algorithm SelfConsumptionConfE1) Execute(ctx context.Context) error {
	pointsTypes, ssaaConfig, genTypesCount, _, gdConfig := getCupsType(*algorithm.B)
	for i, _ := range algorithm.B.Curve {
		sumEHGN := 0.0
		sumEHSA := 0.0
		sumEHEX := 0.0
		sumEHAU := 0.0
		sumEHCR := 0.0
		sumAEgd := 0.0
		for c, _ := range algorithm.B.Curve[i].Points {
			point := &algorithm.B.Curve[i].Points[c]
			curve := &algorithm.B.Curve[i]
			typePoint := pointsTypes[point.CUPS]

			vAIssaa := new(float64)
			cupsSsaa, ssaaCupsFound := ssaaConfig[point.CUPS]
			if typePoint == GdServicePointType && ssaaCupsFound {
				for _, curvePoint := range curve.Points {
					if curvePoint.CUPS == cupsSsaa {
						vAIssaa = &curvePoint.AI
					}
				}
			}
			algorithm.setEHGN(ctx, &sumEHGN, vAIssaa, typePoint, point, curve)
			algorithm.setEHCR(ctx, &sumEHCR, typePoint, point)
			algorithm.setEHEX(ctx, &sumEHEX, &sumEHGN, genTypesCount, typePoint, point, curve)
			algorithm.setEHAU(ctx, &sumEHAU, &sumEHEX, &sumEHGN, gdConfig, typePoint, point, curve)
			algorithm.setEHSA(ctx, &sumEHSA, vAIssaa, &sumAEgd, typePoint, point, curve)
			algorithm.setEHDC(ctx, &sumEHAU, &sumEHEX, &sumEHCR, typePoint, point)
		}
	}
	return nil
}
