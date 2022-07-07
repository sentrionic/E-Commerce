package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"entgo.io/ent/dialect"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sentrionic/ecommerce/common/token"
	"github.com/sentrionic/ecommerce/products/ent"
	"github.com/sentrionic/ecommerce/products/ent/enttest"
	_ "github.com/sentrionic/ecommerce/products/ent/runtime"
	"github.com/sentrionic/ecommerce/products/publishers"
	"github.com/sentrionic/ecommerce/products/utils"
	"github.com/stretchr/testify/assert"
	"modernc.org/sqlite"
	"os"
	"testing"
	"time"
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
	client        *ent.Client
	router        *gin.Engine
	mockPublisher *publishers.MockProductPublisher
}

func setupTest(t *testing.T) testHandler {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_URL", "asdasd")
	os.Setenv("SESSION_SECRET", "asdojashouidohasd")
	os.Setenv("NATS_CLIENT_ID", "id")
	os.Setenv("NATS_URL", "url")
	os.Setenv("NATS_CLUSTER_ID", "cluster")

	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")

	mockPublisher := new(publishers.MockProductPublisher)

	router := gin.Default()

	cfg, err := utils.LoadConfig(context.Background())
	assert.NoError(t, err)

	NewHandler(&Config{
		R:      router,
		DB:     client,
		P:      mockPublisher,
		Config: cfg,
	})

	return testHandler{
		router:        router,
		client:        client,
		mockPublisher: mockPublisher,
	}
}

func setupCookie(t *testing.T, userId uuid.UUID) string {
	accessToken, err := token.CreateToken(userId, time.Minute, os.Getenv("SESSION_SECRET"))
	assert.NoError(t, err)
	return fmt.Sprintf("session=%s; Path=/; Max-Age=3600; HttpOnly", accessToken)
}

func addProduct(t *testing.T, ctx context.Context, client *ent.Client, userId uuid.UUID) *ent.Product {
	product, err := client.Product.Create().
		SetTitle("Test Title").
		SetPrice(20).
		SetUserID(userId).
		Save(ctx)
	assert.NoError(t, err)
	return product
}
