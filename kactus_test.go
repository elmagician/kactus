package kactus_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	googlePubSub "cloud.google.com/go/pubsub"
	"github.com/DATA-DOG/go-txdb"
	"github.com/elmagician/godog"
	"github.com/elmagician/godog/colors"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
	"google.golang.org/api/option"

	"github.com/elmagician/kactus/features/definitions"
	"github.com/elmagician/kactus/features/interfaces/database"
	"github.com/elmagician/kactus/features/interfaces/fixtures"
	"github.com/elmagician/kactus/features/interfaces/picker"
	"github.com/elmagician/kactus/features/interfaces/pubsub"
	"github.com/elmagician/kactus/internal/logger"
)

const (
	pgSSLModeEnv  = "POSTGRES_SSL_MODE"
	pgUserEnv     = "POSTGRES_USER"
	pgPasswordEnv = "POSTGRES_PASSWORD"
	pgHostEnv     = "POSTGRES_HOST"
	pgPortEnv     = "POSTGRES_CUSTOM_PORT"
	pgDBNameEnv   = "POSTGRES_DB"

	gcpHostEnv        = "PUBSUB_HOST"
	gcpPortEnv        = "PUBSUB_PORT"
	gcpCredentialsEnv = "PUBSUB_TEST_CREDENTIALS"
	gcpFakeProjectEnv = "GCP_FAKE_PROJECT_ID"

	emulatorEnv = "PUBSUB_EMULATOR_HOST"
)

var (
	opts = godog.Options{
		Output:    colors.Colored(os.Stdout),
		Randomize: -1,
		Format:    "pretty", // can define default values
		Strict:    true,
		Paths:     []string{"test"},
	}

	pickerInstance *picker.Picker
	gcpClient      *pubsub.Google
	postgres       *database.Postgres
	fix            *fixtures.Fixtures
)

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
	txdb.Register("txdb", "postgres", dsn())
}

func TestMain(m *testing.M) {
	pflag.Parse()

	if len(pflag.Args()) > 0 {
		opts.Paths = pflag.Args()
	}

	status := godog.TestSuite{
		Name:                 "kactus",
		TestSuiteInitializer: installKactus,
		ScenarioInitializer:  installDefinitions,
		Options:              &opts,
	}.Run()

	// main test need to be run after godog tests for coverage
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func installDefinitions(s *godog.ScenarioContext) {
	definitions.InstallDebug(s)
	definitions.InstallPicker(s, pickerInstance)
	definitions.InstallFixtures(s, fix)
	definitions.InstallPostgres(s, postgres)
	definitions.InstallGooglePubsub(s, gcpClient)

	// Allows to define a visual step without real assertion.
	s.Step("^I assume that .*", func() error { return nil })
}

func installKactus(s *godog.TestSuiteContext) {
	_ = logger.SetDefault()

	pickerInstance = picker.New()
	pickerInstance.Reset()

	fix = fixtures.New(pickerInstance)
	fix.Reset()

	localDB, err := sql.Open("txdb", "identifier")
	if err != nil {
		panic(err)
	}

	postgres = database.NewPostgres(
		pickerInstance,
		database.PostgresInfo{Key: "database", DB: localDB},
		database.PostgresInfo{Key: "authent", DB: localDB},
	)
	postgres.Reset()

	cli, err := initPubsubClient(context.Background(), getEnv(gcpFakeProjectEnv, "test"))
	if err != nil {
		log.Fatal(err)
	}

	gcpClient = pubsub.NewGoogle(pickerInstance, pubsub.GoogleInfo{Client: cli, Key: "client"})
	gcpClient.Reset()

	s.AfterSuite(func() {
		_ = localDB.Close()
		_ = cli.Close()
	})
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultVal
}

func dsn() string {
	return "sslmode=" + getEnv(pgSSLModeEnv, "disable") +
		" user='" + getEnv(pgUserEnv, "root") +
		"' password='" + getEnv(pgPasswordEnv, "root") +
		"' host='" + getEnv(pgHostEnv, "localhost") +
		"' port='" + getEnv(pgPortEnv, "5432") +
		"' dbname='" + getEnv(pgDBNameEnv, "default") + "'"

}

func initPubsubClient(ctx context.Context, projectID string, opts ...option.ClientOption) (*googlePubSub.Client, error) {
	if os.Getenv(gcpHostEnv) != "" {
		emulatorHost := os.Getenv(gcpHostEnv)
		if os.Getenv(gcpPortEnv) != "" {
			emulatorHost += ":" + os.Getenv(gcpPortEnv)
		}
		if err := os.Setenv(emulatorEnv, emulatorHost); err != nil {
			return nil, err
		}
	}

	if os.Getenv(gcpCredentialsEnv) != "" {
		opts = append(opts, option.WithCredentialsFile(os.Getenv(gcpCredentialsEnv)))
	}

	return googlePubSub.NewClient(ctx, projectID, opts...)
}
