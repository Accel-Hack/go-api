package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Accel-Hack/go-api/internal/app/server/parser"
	"github.com/Accel-Hack/go-api/internal/app/usercase/sample"
	"github.com/gorilla/mux"
)

type InternalSampleHandler struct {
	Usecase sample.Usecase
	Logger  *slog.Logger
}

func (h *InternalSampleHandler) Get(w http.ResponseWriter, r *http.Request) {
	parseID := parser.QueryUUID().Required().Key("id")

	query := r.URL.Query()
	id, err := parseID(query)
	if err != nil {
		h.Logger.Error("parse id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sample, err := h.Usecase.Get(r.Context(), id)
	if err != nil {
		h.Logger.Error("get sample", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(sample); err != nil {
		h.Logger.Error("encode sample to JSON", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *InternalSampleHandler) Search(w http.ResponseWriter, r *http.Request) {
	var (
		parseName   = parser.QueryString().Required().Key("name")
		parseLimit  = parser.QueryInt().OrNil().Key("limit")
		parseOffset = parser.QueryInt().OrNil().Key("offset")
	)

	query := r.URL.Query()
	name, err := parseName(query)
	if err != nil {
		h.Logger.Error("parse name", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limit, err := parseLimit(query)
	if err != nil {
		h.Logger.Error("parse limit", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	offset, err := parseOffset(query)
	if err != nil {
		h.Logger.Error("parse offset", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sample, err := h.Usecase.Search(r.Context(), name, limit, offset)
	if err != nil {
		h.Logger.Error("get sample", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(sample); err != nil {
		h.Logger.Error("encode sample to JSON", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *InternalSampleHandler) Add(w http.ResponseWriter, r *http.Request) {
	var (
		parseName       = parser.QueryString().Required().Key("name")
		parseBirthday   = parser.QueryTime().Required().Key("birthday")
		parseIsJapanese = parser.QueryBool().Required().Key("is_japanese")
	)

	query := r.URL.Query()
	name, err := parseName(query)
	if err != nil {
		h.Logger.Error("parse name", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	birthday, err := parseBirthday(query)
	if err != nil {
		h.Logger.Error("parse birthday", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isJapanese, err := parseIsJapanese(query)
	if err != nil {
		h.Logger.Error("parse is_japanese", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := h.Usecase.Add(r.Context(), sample.AddQuery{
		Name:       name,
		Birthday:   birthday,
		IsJapanese: isJapanese,
	})
	if err != nil {
		h.Logger.Error("add new sample", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	}); err != nil {
		h.Logger.Error("encode sample to JSON", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *InternalSampleHandler) Edit(w http.ResponseWriter, r *http.Request) {
	var (
		parseID         = parser.QueryUUID().Required().Key("id")
		parseName       = parser.QueryString().OrNil().Key("name")
		parseBirthday   = parser.QueryTime().OrNil().Key("birthday")
		parseIsJapanese = parser.QueryBool().OrNil().Key("is_japanese")
	)

	query := r.URL.Query()
	id, err := parseID(query)
	if err != nil {
		h.Logger.Error("parse id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name, err := parseName(query)
	if err != nil {
		h.Logger.Error("parse name", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	birthday, err := parseBirthday(query)
	if err != nil {
		h.Logger.Error("parse birthday", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isJapanese, err := parseIsJapanese(query)
	if err != nil {
		h.Logger.Error("parse is_japanese", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Usecase.Edit(r.Context(), sample.UpdateQuery{
		ID:         id,
		Name:       name,
		Birthday:   birthday,
		IsJapanese: isJapanese,
	}); err != nil {
		h.Logger.Error("edit sample", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	}); err != nil {
		h.Logger.Error("encode sample to JSON", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *InternalSampleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	parseID := parser.QueryUUID().Required().Key("id")

	query := r.URL.Query()
	uid, err := parseID(query)
	if err != nil {
		h.Logger.Error("parse id", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Usecase.Delete(r.Context(), uid); err != nil {
		h.Logger.Error("delete sample", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Route registers routes to mux.Router.
//
//	GET    /sample
//	POST   /sample
//	PUT    /sample
//	DELETE /sample
//	GET    /samples
func (h *InternalSampleHandler) Route(mux *mux.Router) {
	h.Logger.Info(`expose GET "/sample"`)
	mux.HandleFunc("/sample", h.Get).Methods(http.MethodGet)
	h.Logger.Info(`expose PUT "/sample"`)
	mux.HandleFunc("/sample", h.Add).Methods(http.MethodPut)
	h.Logger.Info(`expose POST "/sample"`)
	mux.HandleFunc("/sample", h.Edit).Methods(http.MethodPost)
	h.Logger.Info(`expose DELETE "/sample"`)
	mux.HandleFunc("/sample", h.Delete).Methods(http.MethodDelete)
	h.Logger.Info(`expose GET "/samples"`)
	mux.HandleFunc("/samples", h.Search).Methods(http.MethodGet)
}
