package picker_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/elmagician/kactus/internal/logger"
	"github.com/elmagician/kactus/internal/matchers"
	"github.com/elmagician/kactus/internal/picker"
	. "github.com/elmagician/kactus/internal/test"
	"github.com/elmagician/kactus/internal/types"
)

func init() {
	picker.NoLog()
	matchers.NoLog()
	types.NoLog()
}

func TestUnit_Store_NewStore(t *testing.T) {
	Convey("I should be able to create a new empty store ", t, func() {
		store := picker.NewStore()
		So(store.Disposable, ShouldBeEmpty)
		So(store.Persistent, ShouldBeEmpty)
	})
}

func TestUnit_Store_Pick(t *testing.T) {
	Convey("Given I wish to pick value, ", t, func() {

		store := picker.NewStore()

		Convey("As disposable, it should be stored in Disposable space", func() {
			store.Pick("test", "test", picker.DisposableValue)
			So(store.Persistent, ShouldBeEmpty)
			So(store.Disposable, ShouldNotBeEmpty)
			So(store.Disposable["test"], ShouldEqual, "test")
		})

		Convey("As persistent, it should be stored in Persistent space", func() {
			store.Pick("test", "test", picker.PersistentValue)
			So(store.Disposable, ShouldBeEmpty)
			So(store.Persistent, ShouldNotBeEmpty)
			So(store.Persistent["test"], ShouldEqual, "test")
		})

		Convey("In unknown space, it should log an error", func() {
			core, logs := observer.New(zapcore.DebugLevel)
			logger.Set(zap.New(core))
			So(logs.Len(), ShouldEqual, 0)
			picker.ResetLog()

			store.Pick("test", "test", 0)

			So(store.Disposable, ShouldBeEmpty)
			So(store.Persistent, ShouldBeEmpty)

			So(logs.Len(), ShouldEqual, 1)
			entry := logs.TakeAll()[0]
			So(entry.Message, ShouldEqual, "Trying to pick to unknown scope: 0")
			So(entry.LoggerName, ShouldContainSubstring, ".picker")
		})

	})
}

func TestUnit_Store_Get(t *testing.T) {
	Convey("Given I wish to retrieve value, ", t, func() {

		store := picker.NewStore()

		store.Pick("persistent", "pers", picker.PersistentValue)
		store.Pick("disposable", "disp", picker.DisposableValue)
		store.Pick("persistentOverWrite", "base", picker.PersistentValue)
		store.Pick("persistentOverWrite", "overwritten", picker.DisposableValue)
		store.Pick("pi", 3.14, picker.PersistentValue)
		store.Pick("lol", true, picker.DisposableValue)

		Convey("picked in a single scope, I should retrieve its value", func() {
			val, ok := store.Get("persistent")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "pers")

			val, ok = store.Get("disposable")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "disp")

			val, ok = store.Get("pi")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, 3.14)

			val, ok = store.Get("lol")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, true)
		})

		Convey("picked in disposable and persistent, I should return disposable value", func() {
			val, ok := store.Get("persistentOverWrite")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "overwritten")
		})

		Convey("unpicked, I should return nil, false tuple", func() {
			val, ok := store.Get("unknown")
			So(ok, ShouldBeFalse)
			So(val, ShouldBeNil)
		})

	})
}

func TestUnit_Store_Reset(t *testing.T) {
	Convey("When Reseting store, ", t, func() {

		store := picker.NewStore()

		store.Pick("persistent", "pers", picker.PersistentValue)
		store.Pick("disposable", "disp", picker.DisposableValue)
		store.Pick("persistentOverWrite", "base", picker.PersistentValue)
		store.Pick("persistentOverWrite", "overwritten", picker.DisposableValue)
		store.Pick("pi", 3.14, picker.PersistentValue)
		store.Pick("lol", true, picker.DisposableValue)

		store.Reset()

		Convey("only disposable values should be cleared", func() {
			val, ok := store.Get("persistent")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "pers")

			_, ok = store.Get("disposable")
			So(ok, ShouldBeFalse)

			val, ok = store.Get("pi")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, 3.14)

			_, ok = store.Get("lol")
			So(ok, ShouldBeFalse)

			val, ok = store.Get("persistentOverWrite")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "base")
		})
	})
}

func TestUnit_Store_Del(t *testing.T) {
	Convey("When I want to remove picked value from store, ", t, func() {
		core, logs := observer.New(zapcore.DebugLevel)
		logger.Set(zap.New(core))
		So(logs.Len(), ShouldEqual, 0)
		picker.ResetLog()

		store := picker.NewStore()

		store.Pick("persistent", "pers", picker.PersistentValue)
		store.Pick("disposable", "disp", picker.DisposableValue)
		store.Pick("persistentOverWrite", "base", picker.PersistentValue)
		store.Pick("persistentOverWrite", "overwritten", picker.DisposableValue)
		store.Pick("pi", 3.14, picker.PersistentValue)
		store.Pick("lol", true, picker.DisposableValue)

		Convey("I should be able to forgive persistent values and log a warning", func() {
			store.Del("persistent", picker.PersistentValue)
			_, ok := store.Get("persistent")
			So(ok, ShouldBeFalse)

			So(logs.Len(), ShouldEqual, 1)
			entry := logs.TakeAll()[0]
			So(entry.Message, ShouldEqual, "Deleting persistent value for key: persistent")
			So(entry.Level, ShouldBeEquivalent, zap.WarnLevel)
			So(entry.LoggerName, ShouldContainSubstring, ".picker")

			_, ok = store.Get("pi")
			So(ok, ShouldBeTrue)
		})

		Convey("I should be able to forgive disposable values", func() {
			store.Del("disposable", picker.DisposableValue)
			_, ok := store.Get("disposable")
			So(ok, ShouldBeFalse)

			_, ok = store.Get("lol")
			So(ok, ShouldBeTrue)
		})

		Convey("I should not impact other scope if key exists in multiple spaces.", func() {
			store.Del("persistentOverWrite", picker.DisposableValue)
			val, ok := store.Get("persistentOverWrite")
			So(ok, ShouldBeTrue)
			So(val, ShouldEqual, "base")
		})

		Convey("I should log an error if provided scope does not exists", func() {
			store.Del("persistent", 666)

			So(logs.Len(), ShouldEqual, 1)
			entry := logs.TakeAll()[0]
			So(entry.Message, ShouldEqual, "Trying to delete from unknown scope: 666")
			So(entry.Level, ShouldBeEquivalent, zap.ErrorLevel)
			So(entry.LoggerName, ShouldContainSubstring, ".picker")
		})
	})
}

func TestUnit_Store_GetInstance(t *testing.T) {
	Convey("When I try to get instance", t, func() {
		instanceName := "someInstance"
		instanceValue := "instanceValue"
		instanceKind := picker.GCP

		store := picker.Store{
			Instance: picker.InstanceStore{
				instanceName: picker.InstanceItem{
					Kind:     picker.GCP,
					Instance: instanceValue,
				},
				"otherInstance":
				picker.InstanceItem{
					Kind:     picker.Fixture,
					Instance: "instanceValue",
				},
			},
		}

		Convey("should get instance if it's exist in store", func() {
			kind, instance, exists := store.GetInstance(instanceName)

			So(exists, ShouldBeTrue)
			So(kind, ShouldEqual, instanceKind)
			So(instance, ShouldEqual, instanceValue)
		})

		Convey("should get NoInstance kind if instance does not exist", func() {
			kind, instance, exists := store.GetInstance("foo")

			So(exists, ShouldBeFalse)
			So(kind, ShouldEqual, picker.NoInstance)
			So(instance, ShouldBeNil)
		})
	})
}
