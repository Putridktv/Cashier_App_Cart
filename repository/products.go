package repository

import (
	"cashierAppCart/db"
	"cashierAppCart/model"
	"encoding/json"
)

type ProductRepository struct {
	db db.DB
}

func NewProductRepository(db db.DB) ProductRepository {
	return ProductRepository{db}
}

func (u *ProductRepository) ReadProducts() ([]model.Product, error) {
	records, err := u.db.Load("products") //membaca fungsi load dngan membuat var
	if err != nil {
		return nil, err
	}

	var listProducts []model.Product                     //var untuk menampung
	err = json.Unmarshal([]byte(records), &listProducts) //conv agar dapat ditampilkan
	if err != nil {
		return nil, err
	} else {
		return listProducts, nil
	}
}

func (u *ProductRepository) ResetProducts() error {
	err := u.db.Reset("products", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}
