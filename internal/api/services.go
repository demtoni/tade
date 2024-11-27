package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/demtoni/tade/internal/database"
	manager "github.com/demtoni/tade/internal/manager/sdk"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ServiceRequest struct {
	Name     string                 `json:"name"`
	Months   int                    `json:"months"`
	Location string                 `json:"location"`
	Service  string                 `json:"service"`
	Prolong  bool                   `json:"prolong"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (r *ServiceRequest) Bind(_ *http.Request) error {
	switch {
	case r.Location == "":
		fallthrough
	case r.Service == "":
		fallthrough
	case len(r.Metadata) == 0:
		fallthrough
	case r.Name == "":
		return errors.New(ErrorEmptyField)
	case len(r.Name) > 72:
		return errors.New(ErrorLongServiceName)
	case r.Months <= 0:
		return errors.New(ErrorNegativePeriod)
	}
	return nil
}

type ServiceResponse struct {
	ID        int64                  `json:"id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	ExpiresAt int64                  `json:"expires_at,omitempty"`
	Location  string                 `json:"location,omitempty"`
	Service   string                 `json:"service,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

func (r *ServiceResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (s *Server) GetService(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		s.SendError(w, r, nil, http.StatusNotFound, ErrorServiceNotFound)
		return
	}

	service, err := s.queries.GetService(r.Context(), database.GetServiceParams{
		ID:     int64(id),
		UserID: u.ID,
	})
	if err != nil {
		s.SendError(w, r, nil, http.StatusNotFound, ErrorServiceNotFound)
		return
	}

	var meta map[string]interface{}

	switch service.Type {
	case "shadowsocks":
		meta, err = manager.GetShadowsocks(service.Address, fmt.Sprintf("%s%d", u.Name, service.CreatedAt))
		if err != nil {
			s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}
	}

	render.Render(w, r, &ServiceResponse{
		ID:        int64(id),
		Name:      service.Name,
		ExpiresAt: service.ExpiresAt,
		Location:  service.Name_2,
		Service:   service.Type,
		Metadata:  meta,
	})
}

func NewServiceListResponse(services *[]database.ListUserServicesRow) []render.Renderer {
	list := []render.Renderer{}
	for _, service := range *services {
		list = append(list, &ServiceResponse{
			ID:        service.ID,
			Name:      service.Name,
			ExpiresAt: service.ExpiresAt,
			Location:  service.Name_2,
			Service:   service.Type,
		})
	}
	return list
}

func (s *Server) ListUserServices(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)

	services, err := s.queries.ListUserServices(r.Context(), u.ID)
	if err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	render.RenderList(w, r, NewServiceListResponse(&services))
}

type LocationResponse struct {
	Location string   `json:"name"`
	Services []string `json:"services"`
}

func (r *LocationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewLocationListResponse(locations *[]database.ListLocationsRow) []render.Renderer {
	list := []render.Renderer{}
	for _, location := range *locations {
		list = append(list, &LocationResponse{
			location.Name, strings.Split(location.Services, ","),
		})
	}
	return list
}

func (s *Server) ListLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := s.queries.ListLocations(r.Context())
	if err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	render.RenderList(w, r, NewLocationListResponse(&locations))
}

func (s *Server) CreateService(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)

	data := &ServiceRequest{}
	if err := render.Bind(r, data); err != nil {
		s.SendError(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	price, err := s.queries.GetPrice(r.Context(), data.Service)
	if err != nil {
		s.SendError(w, r, nil, http.StatusNotFound, ErrorServiceUnknown)
		return
	}
	prolongPrice := price * int64(data.Months)
	if prolongPrice > u.Balance {
		s.SendError(w, r, nil, http.StatusForbidden, ErrorLowBalance)
		return
	}

	expiresAt := time.Now().AddDate(0, data.Months, 0).Unix()
	createdAt := time.Now().Unix()

	location, err := s.queries.GetLocation(r.Context(), database.GetLocationParams{
		Column1: sql.NullString{data.Service, true},
		Name:    data.Location,
	})
	if err != nil || location.ID == 0 {
		s.SendError(w, r, nil, http.StatusBadRequest, ErrorLocationNotSupported)
		return
	}

	switch data.Service {
	case "shadowsocks":
		name := fmt.Sprintf("%s%d", u.Name, createdAt)
		method := data.Metadata["method"].(string)
		plugin := data.Metadata["plugin"].(string)
		// TODO: this should be instead placed in a job queue
		if err := manager.DeployShadowsocks(location.Address, name, method, plugin); err != nil {
			s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}
	}

	if err := s.queries.UpdateBalance(r.Context(), database.UpdateBalanceParams{
		u.Balance - prolongPrice,
		u.ID,
	}); err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	prolong := int64(0)
	if data.Prolong {
		prolong = 1
	}

	id, err := s.queries.CreateService(r.Context(), database.CreateServiceParams{
		Name:         data.Name,
		Type:         data.Service,
		CreatedAt:    createdAt,
		ExpiresAt:    expiresAt,
		Prolong:      prolong,
		ProlongPrice: prolongPrice,
		UserID:       u.ID,
		LocationID:   location.ID,
	})
	if err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		if err := s.queries.UpdateBalance(r.Context(), database.UpdateBalanceParams{
			u.Balance,
			u.ID,
		}); err != nil {
			s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}
		return
	}

	render.Render(w, r, &ServiceResponse{ID: id})
}

func (s *Server) CheckServices() error {
	expired, err := s.queries.GetExpiredServices(context.TODO(), time.Now().Unix())
	if err != nil {
		return err
	}

	// TODO: process concurrently from worker pool
	for _, srv := range expired {
		u, err := s.queries.GetUser(context.TODO(), srv.UserID)
		if err != nil {
			return err
		}
		if srv.Prolong > 0 && u.Balance >= srv.ProlongPrice {
			if err := s.queries.ProlongService(context.TODO(), database.ProlongServiceParams{
				time.Now().Unix(), srv.ID,
			}); err != nil {
				return err
			}
			if err := s.queries.UpdateBalance(context.TODO(), database.UpdateBalanceParams{
				u.Balance - srv.ProlongPrice, u.ID,
			}); err != nil {
				return err
			}
			continue
		}
		if err := manager.DeleteShadowsocks(srv.Address, fmt.Sprintf("%s%d", u.Name, srv.CreatedAt)); err != nil {
			return err
		}
		if err := s.queries.DeleteService(context.TODO(), srv.ID); err != nil {
			return err
		}
	}
	return nil
}
