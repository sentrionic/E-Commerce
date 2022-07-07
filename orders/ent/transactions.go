package ent

import (
	"context"
	"fmt"
	ogen "github.com/sentrionic/ecommerce/orders/ent/order"
	pgen "github.com/sentrionic/ecommerce/orders/ent/product"
)

func UpdateProductTx(tx *Tx, prev *Product) error {
	// We begin the update operation:
	n, err := tx.Product.Update().
		// We limit our update to only work on the correct record and version:
		Where(pgen.ID(prev.ID), pgen.Version(prev.Version-1)).
		SetVersion(prev.Version).
		SetPrice(prev.Price).
		SetTitle(prev.Title).
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

func UpdateOrderTx(tx *Tx, prev *Order) error {
	// We begin the update operation:
	n, err := tx.Order.Update().
		// We limit our update to only work on the correct record and version:
		Where(ogen.ID(prev.ID), ogen.Version(prev.Version-1)).
		SetVersion(prev.Version).
		SetStatus(prev.Status).
		Save(context.Background())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	// SaveX returns the number of affected records. If this value is
	// different from 1 the record must have been changed by another
	// process.
	if n != 1 {
		return fmt.Errorf("update failed: order id=%s updated by another process", prev.ID)
	}
	return nil
}
