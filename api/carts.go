package api

import (
	"cashierAppCart/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (api *API) AddCart(w http.ResponseWriter, r *http.Request) {
	// Get username context to struct model.Cart.
	username := fmt.Sprintf("%s", r.Context().Value("username"))

	r.ParseForm()

	// Check r.Form with key product, if not found then return response code 400 and message "Request Product Not Found".

	if len(r.Form) == 0 {
		errorResponse := model.ErrorResponse{Error: "Request Product Not Found"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(convErr)
		return
	} else {
		w.WriteHeader(200)
	}

	var list []model.Product
	for _, formList := range r.Form {
		for _, v := range formList {
			item := strings.Split(v, ",")
			p, _ := strconv.ParseFloat(item[2], 64)
			q, _ := strconv.ParseFloat(item[3], 64)
			total := p * q
			list = append(list, model.Product{
				Id:       item[0],
				Name:     item[1],
				Price:    item[2],
				Quantity: item[3],
				Total:    total,
			})
		}
	}

	// Add data field Name, Cart and TotalPrice with struct model.Cart.
	var totalCart float64
	for _, v := range list {
		totalCart += v.Total
	}

	cart := model.Cart{}
	cart.Name = username
	cart.Cart = list
	cart.TotalPrice = totalCart

	_, err := api.cartsRepo.CartUserExist(cart.Name)
	if err != nil {
		api.cartsRepo.AddCart(cart)
	} else {
		api.cartsRepo.UpdateCart(cart)
	}
	api.dashboardView(w, r)

}
