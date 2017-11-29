package models

import (
	log "github.com/Sirupsen/logrus"
	"github.com/up-finder/silk.web/app/db"
	"github.com/up-finder/silk.web/app/json"
	"gopkg.in/redis.v3"
)

const (
	UUID_QUEUE_SIZE = 10000
	UUID_QUEUE_NAME = "silk:events:uuids:list" //queue in redis with uuids of the models, used for limiting the amount of information stored in redis about the saved entities
	UUID_HASH_NAME  = "silk:events:uuids:hash" //the hash in redis with the last UUID_QUEUE_SIZE saved items
)

type EventStatus int

const (
	EventSaved EventStatus = iota
	EventNotSaved
	FailedToObtainEventStatus
)

// The event is simply a wrapper around json, that can be saved to filedb
type Event struct {
	data json.JSON
}

// Constructor
var NewEvent = func(data json.JSON) Model {
	return &Event{data: data}
}

// Serializes event into string, required for saving in filedb
func (e *Event) Serialize() string {
	return string(e.data)
}

// Gets the status for the event with specific uuid
// Returns EventSaved if the events was saved among the last UUID_QUEUE_SIZE events, EventNotSaved o/w
// and FailedToObtainEventStatus if Redis is unavailable
var GetEventStatus = func(uuid string) EventStatus {
	_, err := db.Redis.HGet(UUID_HASH_NAME, uuid).Result()
	if err != nil && err != redis.Nil {
		log.Errorf("Models: GetEventStatus %v : error getting uuid from redis hash: %v", uuid, err)
		return FailedToObtainEventStatus
	}
	if err == redis.Nil {
		return EventNotSaved
	} else {
		return EventSaved
	}
}

// Performs check, that the entity with that uuid was not saved among the last UUID_QUEUE_SIZE events
// And calls filedb save method to save data
// Returns true if save succeeded, returns files if event with this id was saved mong the last UUID_QUEUE_SIZE events
// or if Redis is unavailable
func (e *Event) Save() bool {
	uuid, ok := e.data.Get("uuid")
	if !ok {
		log.Errorf("Models: Event %v : error finding uuid field", e.data)
		return false
	}
	if uuid == "" {
		log.Errorf("Models: Event %v : uuid field must be not blank", e.data)
		return false
	}
	err := db.Redis.Ping().Err()
	if err != nil {
		log.Errorf("Models: Saving Event %v : error connecting to Redis: %v", e.data, err)
		return false
	}
	_, err = db.Redis.HGet(UUID_HASH_NAME, uuid).Result()
	if err != redis.Nil {
		log.Errorf("Models: Event %v : uuid '%v' already exists", e.data, uuid)
		return false
	}
	db.Redis.LPush(UUID_QUEUE_NAME, uuid)
	size, err := db.Redis.LLen(UUID_QUEUE_NAME).Result()
	if size > UUID_QUEUE_SIZE {
		toDel, err := db.Redis.RPop(UUID_QUEUE_NAME).Result()
		if err != nil {
			log.Errorf("Models: Event %v : error deleting an element from redis queue: %v", e.data, err)
			return false
		}
		err = db.Redis.HDel(toDel).Err()
		if err != nil {
			log.Errorf("Models: Event %v : error deleting an element from redis hash: %v", e.data, err)
			return false
		}
	}
	err = db.Redis.HSet(UUID_HASH_NAME, uuid, e.Serialize()).Err()
	if err != nil {
		log.Errorf("Models: Event %v : error setting redis hash element: %v", e.data, err)
		return false
	}
	db.File.Save(e)
	return true
}
