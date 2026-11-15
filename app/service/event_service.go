package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"eventhub/app/model"
	"eventhub/app/repository"
	"eventhub/helper"

	"github.com/gofiber/fiber/v2"
)

type EventService struct {
	eventRepository repository.EventRepository
}

func NewEventService(eventRepository repository.EventRepository) *EventService {
	return &EventService{eventRepository: eventRepository}
}

func (s *EventService) CreateEvent(c *fiber.Ctx) error {
	var req model.CreateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body request tidak valid")
	}

	if err := validateEventInput(req.Title, req.Quota); err != nil {
		return err
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

	if err := s.eventRepository.Create(c.UserContext(), event); err != nil {
		return mapRepoError(err)
	}

	location := fmt.Sprintf("/api/v1/events/%d", event.ID)

	return helper.Created(c, "event berhasil dibuat", toEventResponse(event), location)
}

func (s *EventService) GetAllEvents(c *fiber.Ctx) error {
	events, err := s.eventRepository.GetAll(c.UserContext())
	if err != nil {
		return mapRepoError(err)
	}

	responses := make([]model.EventResponse, 0, len(events))
	for i := range events {
		responses = append(responses, *toEventResponse(&events[i]))
	}

	return helper.Success(c, fiber.StatusOK, "daftar event", responses)
}

func (s *EventService) GetEventByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.BadRequest("id event tidak valid")
	}

	event, err := s.eventRepository.GetByID(c.UserContext(), id)
	if err != nil {
		return mapRepoError(err)
	}

	return helper.Success(c, fiber.StatusOK, "detail event", toEventResponse(event))
}

func (s *EventService) UpdateEvent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.BadRequest("id event tidak valid")
	}

	var req model.UpdateEventRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body request tidak valid")
	}

	if err := validateEventUpdate(&req); err != nil {
		return err
	}

	event, err := s.eventRepository.GetByID(c.UserContext(), id)
	if err != nil {
		return mapRepoError(err)
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

	if err := s.eventRepository.Update(c.UserContext(), event); err != nil {
		return mapRepoError(err)
	}

	return helper.Success(c, fiber.StatusOK, "event berhasil diperbarui", toEventResponse(event))
}

func (s *EventService) DeleteEvent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.BadRequest("id event tidak valid")
	}

	if err := s.eventRepository.Delete(c.UserContext(), id); err != nil {
		return mapRepoError(err)
	}

	return helper.NoContent(c)
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

func mapRepoError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return helper.NotFound("event tidak ditemukan")
	}
	return helper.Internal(err)
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
