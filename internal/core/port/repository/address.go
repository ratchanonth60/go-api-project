package repository

import (
	"project-api/internal/core/entity"
	"project-api/internal/core/port/utils"
)

type IAddressRepository interface {
	utils.BaseInterface[entity.Address]
}
