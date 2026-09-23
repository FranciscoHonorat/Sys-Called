package application

import (
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
)

type GetPublicKeysUseCase struct {
	keys out.SigningKeys
}

func NewGetPublicKeysUseCase(keys out.SigningKeys) *GetPublicKeysUseCase {
	return &GetPublicKeysUseCase{keys: keys}
}

func (uc *GetPublicKeysUseCase) Execute() []out.PublicKey {
	return uc.keys.PublicKeys()
}
