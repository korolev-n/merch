package service

import (
	"context"

	"github.com/korolev-n/merch-auth/internal/repository"
)

type Shop interface {
	BuyItem(ctx context.Context, userID int, itemType string) error
}

type ShopService struct {
	ShopRepo repository.ShopRepository
}

func NewShopService(repo repository.ShopRepository) *ShopService {
	return &ShopService{
		ShopRepo: repo,
	}
}

func (s *ShopService) BuyItem(ctx context.Context, userID int, itemType string) error {
	item, err := s.ShopRepo.GetItemByType(ctx, itemType)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrItemNotFound
	}

	wallet, err := s.ShopRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if wallet.Balance < item.Price {
		return ErrInsufficientBalance
	}

	newBalance := wallet.Balance - item.Price
	if err := s.ShopRepo.UpdateWalletBalance(ctx, userID, newBalance); err != nil {
		return err
	}

	if err := s.ShopRepo.AddToUserInventory(ctx, userID, item.ID); err != nil {
		return err
	}

	return nil
}
