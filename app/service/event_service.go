package service

import (
	"context"

	"eventhub/app/model"
	"eventhub/app/repository"
)

type EventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepository repository.EventRepository) *EventService {
	return &EventService{eventRepository: eventRepository}
}

func (s *EventService) CreateEvent(ctx context.Context, req *model.CreateEventRequest) (*model.EventResponse, error) {
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

func (s *EventService) GetEventByID(ctx context.Context, id int64) (*model.EventResponse, error) {
	event, err := s.eventRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toEventResponse(event), nil
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
