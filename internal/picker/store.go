package picker

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	// PersistentValue allows to store data in persistent store.
	PersistentValue DataScope = iota + 1
	// DisposableValue allows to store data in disposable store.
	DisposableValue
	// InstanceValue allows to store data in in instance store.
	// An instance is defined by a Key, a Kind and a Value (interface).
	// It can only store Kactus feature or internal instances.
	InstanceValue
)

type (
	// Store picked data.
	Store struct {
		// Disposable stores picked data that will be forgotten on resets.
		Disposable KeyValueStore

		// Persistent stores picked data for all suite duration.
		Persistent KeyValueStore

		// Instance stores picked data as instance values.
		Instance InstanceStore
	}

	// DataScope encourages to use provided constants for
	// data scope storage.
	DataScope int

	// KeyValueStore abstract management of values in store.
	KeyValueStore map[string]interface{}
)

// NewStore initialize picker store.
func NewStore() *Store {
	return &Store{
		Disposable: make(KeyValueStore),
		Persistent: make(KeyValueStore),
		Instance:   make(InstanceStore),
	}
}

// Reset stored Disposable values. It always returns nil but requires error for
// godog usage.
func (store *Store) Reset() {
	log.Debug("Resetting disposable store. /!\\ Persistent store is not reset.")

	store.Disposable = make(KeyValueStore)
}

// Pick store value under key. Providing a true bool to noReset will store
// value in persistent store.
func (store *Store) Pick(key string, value interface{}, scope DataScope) {
	var addTo KeyValueStore

	switch scope {
	case PersistentValue:
		log.Debug("Picking value to persistent store", zap.String("key", key), zap.Reflect("value", value))

		addTo = store.Persistent
	case DisposableValue:
		log.Debug("Picking value to disposable store", zap.String("key", key), zap.Reflect("value", value))

		addTo = store.Disposable
	case InstanceValue:
		val, ok := value.(InstanceItem)
		if !ok {
			tmpVal, ok := value.(*InstanceItem)
			if !ok {
				// Do not return error cause InstanceValue picking SHOULD NEVER be done
				// outside of library.
				log.Error("Picking instance value requires an InstanceItem value.")
				return
			}

			val = *tmpVal
		}

		store.Instance[key] = val

		return
	default:
		log.Error(fmt.Sprintf("Trying to pick to unknown scope: %v", scope))
		return
	}

	addTo[key] = value
}

// Get tries to retrieve stored values. If key is known in both
// Disposable and Persistent values, it will return the Disposable one.
func (store *Store) Get(key string) (value interface{}, exists bool) {
	log.Debug("Trying to retrieve value from Disposable store for: " + key)

	value, exists = store.Disposable[key]

	if !exists {
		log.Debug("Value is not disposable. Trying to retrieve value from Persistent store for: " + key)

		value, exists = store.Persistent[key]
	}

	log.Debug("Get results", zap.Reflect("value", value), zap.Bool("found", exists))

	return value, exists
}

// GetInstance retrieves known kactus instance from picker.
func (store *Store) GetInstance(key string) (kind InstanceKind, instance interface{}, exists bool) {
	var value InstanceItem

	value, exists = store.Instance[key]
	if !exists {
		return NoInstance, nil, false
	}

	if value.Kind == GCP {
		fmt.Println("getting GCP instance", key, "got value: ", value.Instance, fmt.Sprintf("pointer %p", value.Instance))
	}

	return value.Kind, value.Instance, true
}

// Del remove known key from Disposable or Persistent value.
func (store *Store) Del(key string, scope DataScope) {
	switch scope {
	case PersistentValue:
		log.Warn("Deleting persistent value for key: " + key)

		delete(store.Persistent, key)
	case DisposableValue:
		log.Debug("Deleting disposable value for key: " + key)

		delete(store.Disposable, key)
	case InstanceValue:
		log.Debug("Deleting instance value for key: " + key)

		delete(store.Instance, key)
	default:
		log.Error(fmt.Sprintf("Trying to delete from unknown scope: %v", scope))
	}
}
