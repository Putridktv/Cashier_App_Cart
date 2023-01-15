package main

import (
	"cashierAppCart/api"
	"cashierAppCart/db"
	repo "cashierAppCart/repository"
)

// it'ss used for debugging
func main() {
	db := &db.JsonDB{}
	usersRepo := repo.NewUserRepository(db)
	sessionsRepo := repo.NewSessionsRepository(db)
	productsRepo := repo.NewProductRepository(db)
	cartsRepo := repo.NewCartRepository(db)

	mainAPI := api.NewAPI(usersRepo, sessionsRepo, productsRepo, cartsRepo)
	mainAPI.Start()
}
