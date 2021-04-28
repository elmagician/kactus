package picker

const (
	// NoInstance is returned when
	// failed to retrieve instance from store.
	NoInstance InstanceKind = iota

	// Postgres kind refers to postgres.DB instance
	Postgres

	// GCP kind refers to pubsub.GCP instance
	GCP

	// REST refers to http.?? instance
	REST

	// Fixture refers to fixtures.Fixtures instance
	Fixture
)

type (
	// InstanceStore represents the structure used
	// by Store to manage kactus tools instances.
	InstanceStore map[string]InstanceItem

	// InstanceItem describe a kactus tool instance.
	InstanceItem struct {
		Kind     InstanceKind
		Instance interface{}
	}

	// InstanceKind encourages user to pass by
	// instance constant when
	// providing instance kind.
	InstanceKind int
)
