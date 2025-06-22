package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"tracker-backend/internal/domain/services"
	"tracker-backend/internal/interfaces/rest/utils/response"
	"tracker-backend/internal/pkg/service"

	"github.com/go-chi/render"
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
	// get context keys
	ctx := r.Context()

	// parse query params
	query := r.URL.Query().Get("query")
	if query == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("search query is required"))
		return
	}

	slog.Debug("search query", slog.String("query", query))

	// parse limits
	limitTracks, err := strconv.Atoi(r.URL.Query().Get("limitTracks"))
	if err != nil {
		limitTracks = 10 // default value
	}

	generalLimit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		generalLimit = 5 // default value
	}

	// execute service function
	result, err := h.svc.GlobalSearch(ctx, query, limitTracks, generalLimit)
	if err != nil {
		if errors.Is(err, service.ErrInternal) {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	// return result
	render.JSON(w, r, result)
}
