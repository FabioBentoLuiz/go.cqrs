// Copyright 2016 Jet Basrawi. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package eventsourcing

import (
	"context"
	"encoding/json"
	"fmt"

	esDb "github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/direction"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/streamrevision"
	"github.com/gofrs/uuid"
)

// DomainRepository is the interface that all domain repositories should implement.
type DomainRepository interface {
	//Loads an aggregate of the given type and ID
	Load(aggregateTypeName string, aggregateID string) (AggregateRoot, error)

	//Saves the aggregate.
	Save(aggregate AggregateRoot, expectedVersion *int64) error
}

// GetEventStoreCommonDomainRepo is an implementation of the DomainRepository
// that uses GetEventStore for persistence
type GetEventStoreCommonDomainRepo struct {
	eventStore         *esDb.Client
	eventBus           EventBus
	streamNameDelegate StreamNamer
	aggregateFactory   AggregateFactory
	eventFactory       EventFactory
}

// NewCommonDomainRepository constructs a new CommonDomainRepository
func NewCommonDomainRepository(eventStore *esDb.Client, eventBus EventBus) (*GetEventStoreCommonDomainRepo, error) {
	if eventStore == nil {
		return nil, fmt.Errorf("nil Eventstore injected into repository")
	}

	if eventBus == nil {
		return nil, fmt.Errorf("nil EventBus injected into repository")
	}

	d := &GetEventStoreCommonDomainRepo{
		eventStore: eventStore,
		eventBus:   eventBus,
	}
	return d, nil
}

// SetAggregateFactory sets the aggregate factory that should be used to
// instantate aggregate instances
//
// Only one AggregateFactory can be registered at any one time.
// Any registration will overwrite the provious registration.
func (r *GetEventStoreCommonDomainRepo) SetAggregateFactory(factory AggregateFactory) {
	r.aggregateFactory = factory
}

// SetEventFactory sets the event factory that should be used to instantiate event
// instances.
//
// Only one event factory can be set at a time. Any subsequent registration will
// overwrite the previous factory.
func (r *GetEventStoreCommonDomainRepo) SetEventFactory(factory EventFactory) {
	r.eventFactory = factory
}

// SetStreamNameDelegate sets the stream name delegate
func (r *GetEventStoreCommonDomainRepo) SetStreamNameDelegate(delegate StreamNamer) {
	r.streamNameDelegate = delegate
}

// Load will load all events from a stream and apply those events to an aggregate
// of the type specified.
//
// The aggregate type and id will be passed to the configured StreamNamer to
// get the stream name.
func (r *GetEventStoreCommonDomainRepo) Load(aggregateType, id string) (AggregateRoot, error) {

	if r.aggregateFactory == nil {
		return nil, fmt.Errorf("the common domain repository has no Aggregate Factory")
	}

	if r.streamNameDelegate == nil {
		return nil, fmt.Errorf("the common domain repository has no stream name delegate")
	}

	if r.eventFactory == nil {
		return nil, fmt.Errorf("the common domain has no Event Factory")
	}

	aggregate := r.aggregateFactory.GetAggregate(aggregateType, id)
	if aggregate == nil {
		return nil, fmt.Errorf("the repository has no aggregate factory registered for aggregate type: %s", aggregateType)
	}

	streamName, err := r.streamNameDelegate.GetStreamName(aggregateType, id)
	if err != nil {
		return nil, err
	}

	events, err := r.eventStore.ReadStreamEvents(context.Background(), direction.Forwards, streamName, streamrevision.StreamRevisionStart, 1, false)
	if err != nil {
		return nil, fmt.Errorf("could not read events from stream %s", streamName)
	}

	for _, event := range events {
		evtNum := int64(event.EventNumber)
		em := NewEventMessage(id, event, &evtNum)
		aggregate.Apply(em, false)
		aggregate.IncrementVersion()
	}

	return aggregate, nil

}

// Save persists an aggregate
func (r *GetEventStoreCommonDomainRepo) Save(aggregate AggregateRoot, expectedVersion *int64) error {

	if r.streamNameDelegate == nil {
		return fmt.Errorf("the common domain repository has no stream name delagate")
	}

	resultEvents := aggregate.GetChanges()

	streamName, err := r.streamNameDelegate.GetStreamName(typeOf(aggregate), aggregate.AggregateID())
	if err != nil {
		return err
	}

	if len(resultEvents) > 0 {

		events := make([]messages.ProposedEvent, len(resultEvents))

		for k, v := range resultEvents {
			//TODO: There is no test for this code
			v.SetHeader("AggregateID", aggregate.AggregateID())
			//evs[k] = goes.NewEvent("", v.EventType(), v.Event(), v.GetHeaders())
			eventID, err := uuid.NewV4()
			if err != nil {
				return fmt.Errorf("could not generate UUID")
			}

			json, err := json.Marshal(v.Event())
			if err != nil {
				return fmt.Errorf("Error parsing %v", v.Event())
			}

			events[k] = messages.ProposedEvent{
				EventID:      eventID,
				EventType:    v.EventType(),
				ContentType:  "application/json",
				UserMetadata: nil,
				Data:         json,
			}
		}

		_, err := r.eventStore.AppendToStream(context.Background(), streamName, streamrevision.StreamRevisionAny, events)

		if err != nil {
			return fmt.Errorf("unexpected failure appending to stream %s. Error: %+v", streamName, err)
		}
	}

	aggregate.ClearChanges()

	for k, v := range resultEvents {
		if expectedVersion == nil {
			r.eventBus.PublishEvent(v)
		} else {
			ver := int64(*expectedVersion + int64(k) + 1)
			em := NewEventMessage(v.AggregateID(), v.Event(), &ver)
			r.eventBus.PublishEvent(em)
		}
	}

	return nil
}
