package initializers

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	//postgres://hmzmzsiv:FbyRQCACS46OMJPhaMXDW3SI5OQDpd1h@satao.db.elephantsql.com/hmzmzsiv
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Unable to connect to DB")
	}
}
