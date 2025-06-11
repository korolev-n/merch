package service

import (
	"context"

	"github.com/korolev-n/merch/internal/domain"
	"github.com/korolev-n/merch/internal/repository"
)

type Info interface {
	GetUserInfo(ctx context.Context, userID int) (*domain.InfoResponse, error)
}

type InfoService struct {
	Repo repository.InfoRepository
}

func NewInfoService(repo repository.InfoRepository) *InfoService {
	return &InfoService{Repo: repo}
}

func (s *InfoService) GetUserInfo(ctx context.Context, userID int) (*domain.InfoResponse, error) {
	balance, err := s.Repo.GetWalletBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.Repo.GetUserInventory(ctx, userID)
	if err != nil {
		return nil, err
	}

	received, err := s.Repo.GetReceivedCoins(ctx, userID)
	if err != nil {
		return nil, err
	}

	sent, err := s.Repo.GetSentCoins(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.InfoResponse{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: domain.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}
