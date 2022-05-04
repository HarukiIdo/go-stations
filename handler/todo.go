package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/k0kubun/pp/v3"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var createTodoRequest model.CreateTODORequest

		// デコード
		err := json.NewDecoder(r.Body).Decode(&createTodoRequest)
		if err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)
			return
		}
		defer r.Body.Close()

		// Subjectが空の時、エラーを返す
		if createTodoRequest.Subject == "" {
			http.Error(
				w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)
			return
		}

		// SubJectが空でない時、Todoを新規作成
		todo, err := h.svc.CreateTODO(
			r.Context(),
			createTodoRequest.Subject,
			createTodoRequest.Description,
		)
		if err != nil {
			return
		}

		// エンコード
		m := &model.CreateTODOResponse{TODO: *todo}
		if err := json.NewEncoder(w).Encode(m); err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)
			return
		}

	case http.MethodPut:
		var updateTODORequest model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&updateTODORequest); err != nil {
			return
		}
		defer r.Body.Close()

		if updateTODORequest.ID == 0 || updateTODORequest.Subject == "" {
			http.Error(
				w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)
			return
		}
		todo, err := h.svc.UpdateTODO(r.Context(), updateTODORequest.ID, updateTODORequest.Subject, updateTODORequest.Description)
		if err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusNotFound),
				http.StatusBadRequest,
			)
			return
		}
		m := model.UpdateTODOResponse{TODO: *todo}
		if err := json.NewEncoder(w).Encode(m); err != nil {
			return
		}

	case http.MethodGet:
		var readTODORequest model.ReadTODORequest

		params := r.URL.Query()

		var err error
		var prevID int64 = 0
		prevIDParam := params.Get("prev_id")
		if prevIDParam != "" {
			prevID, err = strconv.ParseInt(prevIDParam, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			readTODORequest.PrevID = prevID
		}
		var size int64 = 0
		sizeParam := params.Get("size")
		if sizeParam != "" {
			size, err = strconv.ParseInt(sizeParam, 10, 64)
			if err != nil {
				return
			}
			readTODORequest.Size = size
		}
		pp.Printf("prevID:%v\n", prevID)
		pp.Printf("size:%v\n", size)

		// todosを取得
		todos, err := h.svc.ReadTODO(
			r.Context(),
			readTODORequest.PrevID,
			readTODORequest.Size,
		)
		if err != nil {
			log.Println(err)
			return
		}

		pp.Printf("len:%v\n", len(todos))
		pp.Println(todos)

		readTODOResponse := &model.ReadTODOResponse{
			TODOs: todos,
		}
		if err := json.NewEncoder(w).Encode(readTODOResponse); err != nil {
			return
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
