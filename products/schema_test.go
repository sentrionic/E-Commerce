package main

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sentrionic/ecommerce/products/ent"
	"github.com/sentrionic/ecommerce/products/ent/enttest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")

	t.Run("it implements optimistic concurrency control", func(t *testing.T) {
		// Create an instance of a product
		product, err := client.Product.Create().
			SetTitle("Product").
			SetPrice(20).
			SetUserID(uuid.New()).
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
