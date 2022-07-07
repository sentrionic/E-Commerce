// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	entorder "github.com/sentrionic/ecommerce/orders/ent/order"
	"github.com/sentrionic/ecommerce/orders/ent/product"
	"github.com/sentrionic/ecommerce/orders/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	entorderMixin := schema.Order{}.Mixin()
	entorderMixinFields0 := entorderMixin[0].Fields()
	_ = entorderMixinFields0
	entorderFields := schema.Order{}.Fields()
	_ = entorderFields
	// entorderDescCreatedAt is the schema descriptor for created_at field.
	entorderDescCreatedAt := entorderMixinFields0[0].Descriptor()
	// entorder.DefaultCreatedAt holds the default value on creation for the created_at field.
	entorder.DefaultCreatedAt = entorderDescCreatedAt.Default.(func() time.Time)
	// entorderDescUpdatedAt is the schema descriptor for updated_at field.
	entorderDescUpdatedAt := entorderMixinFields0[1].Descriptor()
	// entorder.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	entorder.DefaultUpdatedAt = entorderDescUpdatedAt.Default.(func() time.Time)
	// entorder.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	entorder.UpdateDefaultUpdatedAt = entorderDescUpdatedAt.UpdateDefault.(func() time.Time)
	// entorderDescVersion is the schema descriptor for version field.
	entorderDescVersion := entorderFields[4].Descriptor()
	// entorder.DefaultVersion holds the default value on creation for the version field.
	entorder.DefaultVersion = entorderDescVersion.Default.(uint)
	// entorderDescID is the schema descriptor for id field.
	entorderDescID := entorderFields[0].Descriptor()
	// entorder.DefaultID holds the default value on creation for the id field.
	entorder.DefaultID = entorderDescID.Default.(func() uuid.UUID)
	productMixin := schema.Product{}.Mixin()
	productMixinFields0 := productMixin[0].Fields()
	_ = productMixinFields0
	productFields := schema.Product{}.Fields()
	_ = productFields
	// productDescCreatedAt is the schema descriptor for created_at field.
	productDescCreatedAt := productMixinFields0[0].Descriptor()
	// product.DefaultCreatedAt holds the default value on creation for the created_at field.
	product.DefaultCreatedAt = productDescCreatedAt.Default.(func() time.Time)
	// productDescUpdatedAt is the schema descriptor for updated_at field.
	productDescUpdatedAt := productMixinFields0[1].Descriptor()
	// product.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	product.DefaultUpdatedAt = productDescUpdatedAt.Default.(func() time.Time)
	// product.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	product.UpdateDefaultUpdatedAt = productDescUpdatedAt.UpdateDefault.(func() time.Time)
	// productDescTitle is the schema descriptor for title field.
	productDescTitle := productFields[1].Descriptor()
	// product.TitleValidator is a validator for the "title" field. It is called by the builders before save.
	product.TitleValidator = productDescTitle.Validators[0].(func(string) error)
	// productDescPrice is the schema descriptor for price field.
	productDescPrice := productFields[2].Descriptor()
	// product.PriceValidator is a validator for the "price" field. It is called by the builders before save.
	product.PriceValidator = productDescPrice.Validators[0].(func(int) error)
	// productDescVersion is the schema descriptor for version field.
	productDescVersion := productFields[3].Descriptor()
	// product.DefaultVersion holds the default value on creation for the version field.
	product.DefaultVersion = productDescVersion.Default.(uint)
	// productDescID is the schema descriptor for id field.
	productDescID := productFields[0].Descriptor()
	// product.DefaultID holds the default value on creation for the id field.
	product.DefaultID = productDescID.Default.(func() uuid.UUID)
}
