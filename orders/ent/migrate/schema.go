// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// OrdersColumns holds the columns for the "orders" table.
	OrdersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"created", "cancelled", "awaiting-payment", "complete"}},
		{Name: "expires_at", Type: field.TypeTime},
		{Name: "user_id", Type: field.TypeUUID},
		{Name: "version", Type: field.TypeUint, Default: 0},
		{Name: "product_order", Type: field.TypeUUID, Unique: true, Nullable: true},
	}
	// OrdersTable holds the schema information for the "orders" table.
	OrdersTable = &schema.Table{
		Name:       "orders",
		Columns:    OrdersColumns,
		PrimaryKey: []*schema.Column{OrdersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "orders_products_order",
				Columns:    []*schema.Column{OrdersColumns[7]},
				RefColumns: []*schema.Column{ProductsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ProductsColumns holds the columns for the "products" table.
	ProductsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "title", Type: field.TypeString},
		{Name: "price", Type: field.TypeInt},
		{Name: "version", Type: field.TypeUint, Default: 0},
	}
	// ProductsTable holds the schema information for the "products" table.
	ProductsTable = &schema.Table{
		Name:       "products",
		Columns:    ProductsColumns,
		PrimaryKey: []*schema.Column{ProductsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		OrdersTable,
		ProductsTable,
	}
)

func init() {
	OrdersTable.ForeignKeys[0].RefTable = ProductsTable
}