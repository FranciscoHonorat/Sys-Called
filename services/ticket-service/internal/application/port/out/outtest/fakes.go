package outtest

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/notification"
)

type EventStore struct {
	streams map[uuid.UUID][]event.Event
	actors  map[uuid.UUID][]string
	order   []uuid.UUID
}

func NewEventStore() *EventStore {
	return &EventStore{streams: make(map[uuid.UUID][]event.Event), actors: make(map[uuid.UUID][]string)}
}

func (s *EventStore) ActorsOf(aggregateID uuid.UUID) []string {
	return s.actors[aggregateID]
}

func (s *EventStore) Append(_ context.Context, aggregateID uuid.UUID, events []event.Event, expectedVersion int, actorID string) error {
	current := s.streams[aggregateID]
	if len(current) != expectedVersion {
		return domainErr.ErrConcurrencyConflict
	}
	if len(current) == 0 {
		s.order = append(s.order, aggregateID)
	}
	s.streams[aggregateID] = append(current, events...)
	for range events {
		s.actors[aggregateID] = append(s.actors[aggregateID], actorID)
	}
	return nil
}

func (s *EventStore) Load(_ context.Context, aggregateID uuid.UUID) ([]event.Event, error) {
	events, ok := s.streams[aggregateID]
	if !ok {
		return nil, domainErr.ErrEventStreamNotFound
	}
	return events, nil
}

func (s *EventStore) ListAggregateIDs(_ context.Context) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(s.order))
	copy(ids, s.order)
	return ids, nil
}

func (s *EventStore) ForgetAll() {
	s.streams = make(map[uuid.UUID][]event.Event)
	s.actors = make(map[uuid.UUID][]string)
	s.order = nil
}

type ResponsibleDirectory struct {
	Responsibles []string
	Names        map[string]string
	Err          error
}

func (d *ResponsibleDirectory) List(_ context.Context) ([]string, error) {
	return d.Responsibles, d.Err
}

func (d *ResponsibleDirectory) ListWithNames(_ context.Context) ([]out.Responsible, error) {
	responsibles := make([]out.Responsible, 0, len(d.Responsibles))
	for _, id := range d.Responsibles {
		responsibles = append(responsibles, out.Responsible{ID: id, Name: d.Names[id]})
	}
	return responsibles, d.Err
}

func (d *ResponsibleDirectory) Upsert(_ context.Context, id, name string) error {
	if d.Err != nil {
		return d.Err
	}
	if d.Names == nil {
		d.Names = make(map[string]string)
	}
	if _, exists := d.Names[id]; !exists {
		d.Responsibles = append(d.Responsibles, id)
	}
	d.Names[id] = name
	return nil
}

type TokenVerifier struct {
	Tokens map[string]actor.Actor
	Err    error
	Calls  int
}

func (v *TokenVerifier) Verify(_ context.Context, token string) (actor.Actor, error) {
	v.Calls++
	if v.Err != nil {
		return actor.Actor{}, v.Err
	}
	a, ok := v.Tokens[token]
	if !ok {
		return actor.Actor{}, errors.New("unknown token")
	}
	return a, nil
}

func Actor(id, role string) actor.Actor {
	a, err := actor.New(id, id, role)
	if err != nil {
		panic(err)
	}
	return a
}

type NotificationStore struct {
	Added   []notification.Notification
	Seen    map[string]time.Time
	AddErr  error
	ListErr error
}

func (s *NotificationStore) Add(_ context.Context, notifications []notification.Notification) error {
	if s.AddErr != nil {
		return s.AddErr
	}
	s.Added = append(s.Added, notifications...)
	return nil
}

func (s *NotificationStore) ListFor(_ context.Context, viewer actor.Actor, limit int) ([]notification.Notification, error) {
	if s.ListErr != nil {
		return nil, s.ListErr
	}
	var result []notification.Notification
	for i := len(s.Added) - 1; i >= 0 && len(result) < limit; i-- {
		n := s.Added[i]
		addressed := n.Audience == notification.ToUser(viewer.ID()) || n.Audience == notification.ToRole(viewer.Role())
		if addressed && n.ActorID != viewer.ID() {
			result = append(result, n)
		}
	}
	return result, nil
}

func (s *NotificationStore) SeenUntil(_ context.Context, userID string) (time.Time, error) {
	return s.Seen[userID], nil
}

func (s *NotificationStore) MarkSeen(_ context.Context, userID string, at time.Time) error {
	if s.Seen == nil {
		s.Seen = make(map[string]time.Time)
	}
	s.Seen[userID] = at
	return nil
}
