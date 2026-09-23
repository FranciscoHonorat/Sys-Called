package application

import "context"

type LogoutUseCase struct {
	sessions *Sessions
}

func NewLogoutUseCase(sessions *Sessions) *LogoutUseCase {
	return &LogoutUseCase{sessions: sessions}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	return uc.sessions.End(ctx, refreshToken)
}
