package application

import (
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetApplications(req *ApplicationListRequest) (*ApplicationListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	applications, err := s.repo.GetAll(limit, offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	var appResponses []*ApplicationResponse
	for _, app := range applications {
		appResponses = append(appResponses, toApplicationResponse(app))
	}

	return &ApplicationListResponse{
		Data:    appResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(appResponses)) < total,
	}, nil
}

func (s *Service) GetApplicationByID(id int64) (*ApplicationResponse, error) {
	app, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return toApplicationResponse(app), nil
}

func (s *Service) GetApplicationByCode(code string) (*ApplicationResponse, error) {
	app, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	return toApplicationResponse(app), nil
}

func (s *Service) CreateApplication(req *CreateApplicationRequest) (*ApplicationResponse, error) {
	app := &Application{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Icon:        req.Icon,
		URL:         req.URL,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}

	if err := s.repo.Create(app); err != nil {
		return nil, err
	}

	return toApplicationResponse(app), nil
}

func (s *Service) UpdateApplication(id int64, req *UpdateApplicationRequest) (*ApplicationResponse, error) {
	app, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Only update fields that are provided (not empty)
	if req.Name != "" {
		app.Name = req.Name
	}
	if req.Code != "" {
		app.Code = req.Code
	}
	if req.Description != "" {
		app.Description = req.Description
	}
	if req.Icon != "" {
		app.Icon = req.Icon
	}
	if req.URL != "" {
		app.URL = req.URL
	}
	if req.IsActive != nil {
		app.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		app.SortOrder = *req.SortOrder
	}

	if err := s.repo.Update(app); err != nil {
		return nil, err
	}

	return toApplicationResponse(app), nil
}

func (s *Service) DeleteApplication(id int64) error {
	return s.repo.Delete(id)
}

func (s *Service) GetPlanApplications(planID int64) ([]*PlanApplicationResponse, error) {
	return s.repo.GetPlanApplications(planID)
}

func (s *Service) AddApplicationsToPlan(planID int64, req *PlanApplicationRequest) error {
	return s.repo.AddApplicationsToPlan(planID, req.ApplicationIDs, req.IsIncluded)
}

func (s *Service) RemoveApplicationFromPlan(planID, applicationID int64) error {
	return s.repo.RemoveApplicationFromPlan(planID, applicationID)
}

// Helper function
func toApplicationResponse(app *Application) *ApplicationResponse {
	if app == nil {
		return nil
	}

	return &ApplicationResponse{
		ID:          app.ID,
		Name:        app.Name,
		Code:        app.Code,
		Description: app.Description,
		Icon:        app.Icon,
		URL:         app.URL,
		IsActive:    app.IsActive,
		SortOrder:   app.SortOrder,
		CreatedAt:   app.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   app.UpdatedAt.Format(time.RFC3339),
	}
}
