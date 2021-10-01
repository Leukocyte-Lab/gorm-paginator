package paginator

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
