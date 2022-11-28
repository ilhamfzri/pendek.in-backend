package service

import (
	"context"
	"testing"

	"github.com/ilhamfzri/pendek.in/app/database"
	"github.com/ilhamfzri/pendek.in/app/logger"
)

var ctx = context.Background()
var db = database.NewDatabaseConnectionMock()
var log = new(logger.Logger)

func TestMain(m *testing.M) {
	m.Run()
}
