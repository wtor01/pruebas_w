package postgres

import (
	"bitbucket.org/sercide/data-ingestion/internal/inventory"
	"context"
	"fmt"
)

var validOrderColumnsDistributor = map[string]struct{}{
	"name":       struct{}{},
	"id":         struct{}{},
	"created_at": struct{}{},
	"created_by": struct{}{},
	"r1":         struct{}{},
}

var validColumnsDistributor = map[string]struct{}{
	"is_smarkia_active": struct{}{},
}

func (d Distributor) toDomain() inventory.Distributor {
	return inventory.Distributor{
		ID:        d.ID,
		Name:      d.Name,
		R1:        d.R1,
		CDOS:      d.CDOS,
		SmarkiaId: d.SmarkiaId,
		CreatedAt: d.CreatedAt,
		CreatedBy: d.CreatedBy,
		UpdatedAt: d.UpdatedAt,
	}
}

func (d InventoryPostgres) ListDistributors(ctx context.Context, search inventory.Search) ([]inventory.Distributor, int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	distributorsResult := make([]inventory.Distributor, 0, search.Limit)

	distributors := make([]Distributor, 0)

	paginate := NewPaginate(search.Limit, search.Offset)
	query := d.db.WithContext(ctxTimeout).
		Scopes(
			WithLike(search.Q, "name"),
			WithPaginate(paginate),
			WithOrder(validOrderColumnsDistributor, search.Sort),
			WithValues(validColumnsDistributor, search.Values),
		)
	if search.CurrentUser.ID != "" && !search.CurrentUser.IsAdmin {
		query.Where("id in (?)", search.CurrentUser.GetReadDistributors())
	}
	result, count := FetchAndCount(query.Find(&distributors))

	if result.Error != nil {
		return distributorsResult, 0, result.Error
	}

	for _, ds := range distributors {
		distributorsResult = append(distributorsResult, ds.toDomain())
	}

	return distributorsResult, int(count), nil
}

func (d InventoryPostgres) GetDistributorBySmarkiaID(ctx context.Context, id string) (inventory.Distributor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var distributor Distributor

	tx := d.db.WithContext(ctxTimeout).First(&distributor, "smarkia_id = ?", id)

	if tx.Error != nil {
		return inventory.Distributor{}, tx.Error
	}

	return distributor.toDomain(), nil
}

func (d InventoryPostgres) GetDistributorByID(ctx context.Context, id string) (inventory.Distributor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var distributor Distributor

	tx := d.db.WithContext(ctxTimeout).First(&distributor, "id = ?", id)

	if tx.Error != nil {
		return inventory.Distributor{}, tx.Error
	}

	return distributor.toDomain(), nil
}

func (d InventoryPostgres) GetDistributorByCdos(ctx context.Context, cdos string) (inventory.Distributor, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, d.dbTimeout)
	defer cancel()

	var distributor Distributor
	keyDistributorRedis := fmt.Sprintf("distributor:%v:%s", "cdos", cdos)

	if d.redisCache != nil {
		err := d.redisCache.Get(ctx, &distributor, keyDistributorRedis)
		if err == nil {
			return distributor.toDomain(), err
		}
	}

	tx := d.db.WithContext(ctxTimeout).First(&distributor, "cdos = ?", cdos)

	if tx.Error != nil {
		return inventory.Distributor{}, tx.Error
	}

	if d.redisCache != nil {
		err := d.redisCache.Set(ctx, keyDistributorRedis, distributor)
		if err != nil {
			fmt.Println(err)
		}
	}

	return distributor.toDomain(), nil
}
