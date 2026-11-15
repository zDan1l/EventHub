package service

import (
	"context"
	"strings"

	"eventhub/app/model"
	"eventhub/app/repository"
	"eventhub/helper"
)

type EventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepository repository.EventRepository) *EventService {
	return &EventService{eventRepository: eventRepository}
}

func (s *EventService) CreateEvent(ctx context.Context, req *model.CreateEventRequest) (*model.EventResponse, error) {
	if err := validateEventInput(req.Title, req.Quota); err != nil {
		return nil, err
	}

	event := &model.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Quota:       req.Quota,
		IsActive:    req.IsActive,
	}

	if err := s.eventRepository.Create(ctx, event); err != nil {
		return nil, err
	}

	return toEventResponse(event), nil
}

func (s *EventService) GetAllEvents(ctx context.Context) ([]model.EventResponse, error) {
	events, err := s.eventRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]model.EventResponse, 0, len(events))
	for i := range events {
		responses = append(responses, *toEventResponse(&events[i]))
	}

	return responses, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id int) (*model.EventResponse, error) {
	event, err := s.eventRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toEventResponse(event), nil
}

func (s *EventService) UpdateEvent(ctx context.Context, id int, req *model.UpdateEventRequest) (*model.EventResponse, error) {
	if err := validateEventUpdate(req); err != nil {
		return nil, err
	}

	event, err := s.eventRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.StartTime != nil {
		event.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		event.EndTime = *req.EndTime
	}
	if req.Quota != nil {
		event.Quota = *req.Quota
	}
	if req.IsActive != nil {
		event.IsActive = *req.IsActive
	}

	if err := s.eventRepository.Update(ctx, event); err != nil {
		return nil, err
	}

	return toEventResponse(event), nil
}

func (s *EventService) DeleteEvent(ctx context.Context, id int) error {
	return s.eventRepository.Delete(ctx, id)
}

func validateEventInput(title string, quota int) error {
	fields := map[string]string{}

	if strings.TrimSpace(title) == "" {
		fields["title"] = "judul event wajib diisi"
	}

	if quota <= 0 {
		fields["quota"] = "kuota event harus lebih dari 0"
	}

	if len(fields) > 0 {
		return helper.Validation(fields)
	}

	return nil
}

func validateEventUpdate(req *model.UpdateEventRequest) error {
	fields := map[string]string{}

	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		fields["title"] = "judul event wajib diisi"
	}

	if req.Quota != nil && *req.Quota <= 0 {
		fields["quota"] = "kuota event harus lebih dari 0"
	}

	if len(fields) > 0 {
		return helper.Validation(fields)
	}

	return nil
}

func toEventResponse(event *model.Event) *model.EventResponse {
	return &model.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Quota:       event.Quota,
		IsActive:    event.IsActive,
		CreatedAt:   event.CreatedAt,
	}
}
