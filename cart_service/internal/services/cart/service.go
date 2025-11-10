package cartService

import cartRepo "cart_service/internal/repo/cart"

type service struct {
	repo cartRepo.Repo
}

type Service interface {
}

func NewService(repo cartRepo.Repo) Service {

	return &service{
		repo: repo,
	}
}
