package inventory

import (
	"errors"
	"github.com/satori/go.uuid"
	"time"
)

type Distributor struct {
	ID        uuid.UUID
	Name      string
	R1        string
	CDOS      string
	SmarkiaId string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
}

func NewDistributor(name, r1, smarkiaId string) (Distributor, error) {
	if name == "" || smarkiaId == "" || r1 == "" {
		return Distributor{}, errors.New("fields cannot be empty")
	}
	id := uuid.NewV4()

	return Distributor{
		ID:        id,
		Name:      name,
		R1:        r1,
		SmarkiaId: smarkiaId,
		CreatedAt: time.Now(),
		CreatedBy: "",
		UpdatedAt: time.Time{},
	}, nil
}

func ValidateDistributorForm(dt Distributor) error {

	if dt.Name == "" || dt.SmarkiaId == "" || dt.R1 == "" {
		return errors.New("fields cannot be empty")
	}
	return nil
}
