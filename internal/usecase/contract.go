package usecase

import (
	"context"

	"github.com/RizqiSugiarto/coding-test/internal/dto"
)

type Auth interface {
	Login(ctx context.Context, req dto.LoginRequestDTO) (*dto.LoginResponseDTO, error)
}
