package handlers

import (
	"net/http"
	"tracker-backend/internal/domain/services"
)

type SearchHandler struct {
	svc *services.SearchService
}

func NewSearchHandler(searchSvc *services.SearchService) *SearchHandler {
	return &SearchHandler{
		svc: searchSvc,
	}
}

func (h *SearchHandler) GlobalSearch(w http.ResponseWriter, r *http.Request) {

}
