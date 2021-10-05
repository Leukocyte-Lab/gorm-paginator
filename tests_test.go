package paginator

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ory/dockertest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils/tests"
)

func setupMockDB(dialect string) *gorm.DB {
	mock, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer mock.Close()

	var dialector gorm.Dialector
	switch dialect {
	case "mysql":
		// TODO
	case "postgres":
		dialector = postgres.New(
			postgres.Config{
				DriverName: "postgres",
				Conn:       mock,
			})
	case "sqlserver":
		// TODO
	case "sqlite":
		// TODO
	default:
		dialector = tests.DummyDialector{}
	}

	db, err := gorm.Open(dialector, &gorm.Config{DryRun: true})
	if err != nil {
		panic(err)
	}

	return db
}

const ConnectionString = "host=localhost port=%s user=postgres dbname=postgres password=mypassword sslmode=disable"

var pool *dockertest.Pool

func setupResource(t *testing.T) *dockertest.Resource {
	var err error

	if pool == nil {
		pool, err = dockertest.NewPool("")
		if err != nil {
			t.Fatalf("Cloud not connect to docker: %s", err)
		}
	}

	resource, err := pool.Run("postgres", "13", []string{
		"POSTGRES_PASSWORD=mypassword",
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db, err := gorm.Open(postgres.Open(fmt.Sprintf(ConnectionString, resource.GetPort("5432/tcp"))), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return err
		}
		adb, _ := db.DB()
		return adb.Ping()
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	return resource
}

func setupTestDB(t *testing.T, res *dockertest.Resource) *gorm.DB {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(ConnectionString, res.GetPort("5432/tcp"))), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Migrate the schema
	db.AutoMigrate(&tests.User{})

	// Seeding Testing data
	for _, mockUser := range mockUsers {
		db.Create(&mockUser)
	}

	return db
}

func cleanResource(t *testing.T, resource *dockertest.Resource) {
	err := pool.Purge(resource)
	if err != nil {
		t.Error("Error when killing the container resource")
	}
}
