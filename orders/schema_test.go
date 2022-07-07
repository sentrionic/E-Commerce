package main

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/common/order"
	"github.com/sentrionic/ecommerce/orders/ent"
	"github.com/sentrionic/ecommerce/orders/ent/enttest"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSchema_Product(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("it implements optimistic concurrency control for the product", func(t *testing.T) {
		ctx := context.Background()
		client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")

		// Create an instance of a product
		product, err := client.Product.Create().
			SetTitle("Product").
			SetPrice(20).
			Save(ctx)
		assert.NoError(t, err)

		// Read another copy of the same user.
		productCopy := client.Product.GetX(ctx, product.ID)
		productCopy.Price = 25
		productCopy.Version = 1

		// Open a new transaction:
		tx, err := client.Tx(ctx)
		assert.NoError(t, err)

		// Try to update the record once. This should succeed.
		err = ent.UpdateProductTx(tx, productCopy)
		assert.NoError(t, err)

		product.Price = 15
		product.Version = 1
		// Try to update the record a second time. This should fail.
		err = ent.UpdateProductTx(tx, product)
		assert.Error(t, err)
	})
}

func TestSchema_Order(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("it implements optimistic concurrency control for the order", func(t *testing.T) {
		ctx := context.Background()
		client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")

		product, err := client.Product.Create().
			SetTitle("Product").
			SetPrice(20).
			Save(ctx)
		assert.NoError(t, err)

		// Create an instance of a order
		ord, err := client.Order.
			Create().
			SetProduct(product).
			SetUserID(uuid.New()).
			SetStatus(order.Created).
			SetExpiresAt(time.Now().Add(time.Minute * 15)).
			Save(ctx)
		assert.NoError(t, err)

		// Read another copy of the same user.
		orderCopy := client.Order.GetX(ctx, ord.ID)
		orderCopy.Status = order.AwaitingPayment
		orderCopy.Version = 1

		// Open a new transaction:
		tx, err := client.Tx(ctx)
		assert.NoError(t, err)

		// Try to update the record once. This should succeed.
		err = ent.UpdateOrderTx(tx, orderCopy)
		assert.NoError(t, err)

		ord.Status = order.Cancelled
		ord.Version = 1
		// Try to update the record a second time. This should fail.
		err = ent.UpdateOrderTx(tx, ord)
		assert.Error(t, err)
	})
}
