package server

import (
	"errors"
	"log"
	"net/http"
	"pratbacknd/internal/storage"
	"pratbacknd/internal/types"
)

func (s Server) GetCartUser(w http.ResponseWriter, r *http.Request) {
	user, err := s.currentUser(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cart, err := s.storage.GetCart(user.ID)
	if err != nil {
		if errors.Is(err, storage.ErrorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("error retreiving the cart %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, cart)
}

func (s Server) UpdateCartUser(w http.ResponseWriter, r *http.Request) {
	var input types.UpdateUserCartInput
	err := s.readJSON(w, r, &input)
	if err != nil {
		log.Printf("error - reading json: %s \n", err)
		s.errorJSON(w, errors.New("error reading userCart"), http.StatusBadRequest)
		return
	}

	currentUser, err := s.currentUser(w, r)
	if err != nil {
		log.Printf("error - retreiving current user: %s \n", err)
		s.errorJSON(w, nil, http.StatusForbidden)
		return
	}

	log.Printf("---> cart Input: %+v", input)

	cartUpdate, err := s.storage.CreateOrUpdateCart(currentUser.ID, input.ProductID, input.Delta)
	if err != nil {
		log.Printf("error - updating cart: %s \n", err)
		s.errorJSON(w, nil, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, cartUpdate)
}
