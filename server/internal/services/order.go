package services

import (
	"fmt"
	"yt-solutions-server/config"
	"yt-solutions-server/internal/core"
)

type OrderRepo interface {
	ProcessOrder(userId, price int, key *core.Key) error
	ProcessRenewal(ownerId, renewal int, key, expiration string) error
}

type OrderConfig struct {
	MaxFree      int
	MaxLifetime  int
	LicensePrice int
	RenewalPrice int
}

type OrderService struct {
	UserRepo
	KeyRepo
	OrderRepo
	OrderConfig
}

func NewOrderService(u UserRepo, k KeyRepo, o OrderRepo, cfg config.OrdersConfig) *OrderService {
	return &OrderService{
		UserRepo:  u,
		KeyRepo:   k,
		OrderRepo: o,
		OrderConfig: OrderConfig{
			MaxFree:      cfg.MaxFree,
			MaxLifetime:  cfg.MaxLifetime,
			LicensePrice: cfg.Price,
			RenewalPrice: cfg.Renewal,
		},
	}
}

func (o OrderService) GetAvailability() (int, int, int, error) {
	price := o.LicensePrice

	free, err := o.KeyRepo.GetAllFree()
	if err != nil {
		return 0, 0, 0, err
	}

	lifetime, err := o.KeyRepo.GetAllLifetime()
	if err != nil {
		return 0, 0, 0, err
	}

	if o.MaxFree-len(free) > 0 {
		price = 0
	}

	return o.MaxFree - len(free), o.MaxLifetime - len(lifetime), price, nil
}

func (o OrderService) GetRenewalInfo() int {
	return o.OrderConfig.RenewalPrice
}

func (o OrderService) ProcessOrder(userId int) (bool, *core.Key, error) {

	var key core.Key

	user, err := o.UserRepo.GetUserById(userId)
	if err != nil {
		return false, nil, nil
	}

	free, lifetime, _, err := o.GetAvailability()
	if err != nil {
		return false, nil, err
	}

	key.SelectKeyType(free, lifetime)

	if key.Type != core.TypeFree {
		if user.Balance < float64(o.OrderConfig.LicensePrice) {
			return false, nil, nil
		}
	}

	key.GenerateKey(user)
	key.SetExpireDate()
	key.SetOwner(user.ID)

	price := o.selectKeyPrice(key.Type)

	if err = o.OrderRepo.ProcessOrder(user.ID, price, &key); err != nil {
		return false, nil, err
	}

	return true, &key, nil
}

func (o OrderService) selectKeyPrice(keyType string) int {
	if keyType == core.TypeFree {
		return 0
	}

	return o.LicensePrice
}

func (o OrderService) ProcessRenewal(ownerId int) (bool, error) {
	owner, err := o.UserRepo.GetUserById(ownerId)
	if err != nil {
		return false, err
	}

	if !owner.HasKey {
		return false, nil
	}

	if owner.Balance < float64(o.RenewalPrice) {
		return false, err
	}

	key, err := o.KeyRepo.GetKeyByOwnerId(owner.ID)
	if err != nil {
		return false, err
	}
	key.UpdateExpiration()

	fmt.Println(key.Expire)

	if err = o.OrderRepo.ProcessRenewal(owner.ID, o.RenewalPrice, key.Key, key.Expire); err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}
