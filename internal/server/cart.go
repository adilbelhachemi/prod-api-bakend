package server

import (
	"errors"
	"log"
	"net/http"
	"pratbacknd/internal/storage"
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

}
