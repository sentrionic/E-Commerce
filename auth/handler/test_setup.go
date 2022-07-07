package handler

import (
	"database/sql"
	"database/sql/driver"
	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sentrionic/ecommerce/auth/ent"
	"github.com/sentrionic/ecommerce/auth/ent/enttest"
	_ "github.com/sentrionic/ecommerce/auth/ent/runtime"
	"modernc.org/sqlite"
	"os"
	"testing"
)

type sqliteDriver struct {
	*sqlite.Driver
}

func (d sqliteDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}
	c := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	if _, err := c.Exec("PRAGMA foreign_keys = on;", nil); err != nil {
		conn.Close()
		return nil, errors.Wrap(err, "failed to enable enable foreign keys")
	}
	return conn, nil
}

func init() {
	sql.Register("sqlite3", sqliteDriver{Driver: &sqlite.Driver{}})
}

type testHandler struct {
	client *ent.Client
	router *gin.Engine
}

func setupTest(t *testing.T) testHandler {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_URL", "asdasd")
	os.Setenv("SESSION_SECRET", "asdojashouidohasd")
	os.Setenv("Domain", "ecommerce.com")

	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")

	router := gin.Default()

	NewHandler(&Config{
		R:  router,
		DB: client,
	})

	return testHandler{
		router: router,
		client: client,
	}
}
