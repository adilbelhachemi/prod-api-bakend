package server

import (
	"errors"
	"log"
	"net/http"
)

type UpdateInventoryInput struct {
	ProductId string `json:"productId"`
	Delta     int    `json:"delta"`
}

func (s *Server) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	var input UpdateInventoryInput
	err := s.readJSON(w, r, &input)
	if err != nil {
		log.Printf("error - building json: %s \n", err)
		return
	}

	err = s.storage.UpdateInventory(input.ProductId, input.Delta)
	if err != nil {
		log.Printf("error - updating inventory: %s \n", err)
		s.errorJSON(w, errors.New("error updating inventory"), http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, nil)
}
