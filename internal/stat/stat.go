package stat

import (
	"linkshorter/configs"
	"linkshorter/pkg/middleware"
	"linkshorter/pkg/response"
	"net/http"
	"time"
)

const (
	GroupByDay    = "day"
	GroupByMounth = "mounth"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

type StatHandler struct {
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	router.Handle("GET /stat", middleware.IsAuth(handler.GetStat(), deps.Config))
}

func (h *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02 ", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalide from param", http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalide to param", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMounth {
			http.Error(w, "Invalide by param ", http.StatusBadRequest)
			return
		}

		stats := h.StatRepository.GetStats(by, from, to)

		response.Json(w, stats, http.StatusOK)
	}
}
