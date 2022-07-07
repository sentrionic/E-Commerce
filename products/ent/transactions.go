package ent

import (
	"context"
	"fmt"
	gen "github.com/sentrionic/ecommerce/products/ent/product"
)

func UpdateProductTx(tx *Tx, prev *Product) error {
	// We begin the update operation:
	n, err := tx.Product.Update().
		// We limit our update to only work on the correct record and version:
		Where(gen.ID(prev.ID), gen.Version(prev.Version-1)).
		SetVersion(prev.Version).
		SetPrice(prev.Price).
		SetTitle(prev.Title).
		SetNillableOrderID(prev.OrderID).
		Save(context.Background())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	// SaveX returns the number of affected records. If this value is
	// different from 1 the record must have been changed by another
	// process.
	if n != 1 {
		return fmt.Errorf("update failed: product id=%s updated by another process", prev.ID)
	}
	return nil
}
