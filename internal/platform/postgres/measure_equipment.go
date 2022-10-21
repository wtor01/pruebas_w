package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
)

var validOrderColumnsMeasureEquipment = map[string]struct{}{
	"id":            struct{}{},
	"type":          struct{}{},
	"technology":    struct{}{},
	"brand":         struct{}{},
	"model":         struct{}{},
	"serial_number": struct{}{},
}

func (d Meter) toDomain() inventory.MeasureEquipment {
	return inventory.MeasureEquipment{
		ID:                d.ID,
		SerialNumber:      d.SerialNumber,
		Technology:        d.Technology,
		Type:              d.Type,
		Brand:             d.Brand,
		Model:             d.Model,
		ActiveConstant:    d.ActiveConstant,
		ReactiveConstant:  d.ReactiveConstant,
		MaximeterConstant: d.MaximeterConstant,
		CreatedAt:         d.CreatedAt,
		CreatedBy:         d.CreatedByID,
		UpdatedAt:         d.UpdatedAt,
	}
}

func (d InventoryPostgres) ListMeters(ctx context.Context, search inventory.Search) ([]inventory.MeasureEquipment, int, error) {

	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	domainResult := make([]inventory.MeasureEquipment, 0)

	equipments := make([]Meter, 0)

	paginate := NewPaginate(search.Limit, search.Offset)

	result, count := FetchAndCount(d.db.WithContext(ctxTimeout).
		Scopes(
			WithLike(search.Q, "brand", "model", "serial_number"),
			WithPaginate(paginate),
			WithOrder(validOrderColumnsMeasureEquipment, search.Sort),
		).Where("distributor_id = ?", search.DistributorID).
		Find(&equipments))

	if result.Error != nil {
		return domainResult, 0, result.Error
	}

	for _, ds := range equipments {
		domainResult = append(domainResult, ds.toDomain())
	}

	return domainResult, int(count), nil

}

func (d InventoryPostgres) GetMeasureEquipmentByID(ctx context.Context, distributorId string, id string) (inventory.MeasureEquipment, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var equipment Meter

	tx := d.db.WithContext(ctxTimeout).First(&equipment, "distributor_id = ? AND id = ?", distributorId, id)

	if tx.Error != nil {
		return inventory.MeasureEquipment{}, tx.Error
	}

	return equipment.toDomain(), nil
}

func (d InventoryPostgres) GetMeasureEquipmentBySmarkiaID(ctx context.Context, id string) (inventory.MeasureEquipment, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var equipment Meter

	tx := d.db.WithContext(ctxTimeout).First(&equipment, "smarkia_id = ?", id)

	if tx.Error != nil {
		return inventory.MeasureEquipment{}, tx.Error
	}

	return equipment.toDomain(), nil
}

func (d InventoryPostgres) GroupMetersByType(ctx context.Context, distributorId string) (map[string]int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	rows, err := d.db.WithContext(ctxTimeout).Raw(" SELECT meters.type as t, COUNT(*) as count FROM meters WHERE distributor_id = ? AND meters.deleted_at IS NULL GROUP BY meters.type", distributorId).Rows()

	if err != nil {
		return make(map[string]int), err
	}
	defer rows.Close()

	result := make(map[string]int)
	type Result struct {
		T     string
		Count int
	}
	for rows.Next() {
		var r Result
		err = d.db.ScanRows(rows, &r)
		if err != nil {
			return make(map[string]int), err
		}
		result[r.T] = r.Count

	}

	return result, nil
}
