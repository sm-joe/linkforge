package httpserver

import (
	"net/http"

	"github.com/sm-joe/linkforge/internal/link"
)

type ListHandler struct {
	service *link.Service
}

func NewListHandler(
	service *link.Service,
) *ListHandler {

	return &ListHandler{
		service: service,
	}
}

func (h *ListHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	links, err := h.service.ListLinks(
		r.Context(),
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"unable to list links",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		links,
	)
}
