package service

import (
	"errors"
	"market/internal/model"
	"market/internal/repository"
)

var (
	ErrPermissionDenied = errors.New("you have no access")
	ErrNoProduct        = errors.New("product doesn't exists")
	ErrProductExists    = errors.New("product already exists")
)

type ProductService struct {
	productRepo repository.ProductRepo
	userRepo    repository.UserRepo
}

func NewProductService(productRepo repository.ProductRepo, userRepo repository.UserRepo) *ProductService {
	return &ProductService{productRepo: productRepo, userRepo: userRepo}
}

func (s *ProductService) Create(product model.Product) (int, error) {
	id, err := s.productRepo.Create(product)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ProductService) GetAll(q model.ProductQueryInput) ([]model.Product, error) {
	return s.productRepo.GetAll(q)
}
func (s *ProductService) GetProductsByUserID(userID int, q model.ProductQueryInput) ([]model.Product, error) {
	return s.productRepo.GetProductsByUserID(userID, q)
}

func (s *ProductService) GetProductsByCategory(productType string, q model.ProductQueryInput) ([]model.Product, error) {
	return s.productRepo.GetProductsByCategory(productType, q)
}

func (s *ProductService) GetByID(productID int) (model.Product, error) {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (s *ProductService) Update(userID, productID int, input model.UpdateProductInput) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		if err := input.Validate(); err != nil {
			return err
		}
		return s.productRepo.Update(productID, input)
	}

	return ErrPermissionDenied
}

func (s *ProductService) IncreaseViewsCounter(productID int) error {
	views := 1
	input := model.UpdateProductInput{
		Views: &views,
	}
	if err := input.Validate(); err != nil {
		return err
	}
	return s.productRepo.Update(productID, input)
}

func (s *ProductService) Delete(userID, productID int) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Role == model.ADMIN || user.ID == userID {
		return s.productRepo.Delete(productID)
	}

	return ErrPermissionDenied
}
