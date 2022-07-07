package handler

import "github.com/sentrionic/ecommerce/products/ent"

type ProductResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Price  int    `json:"price"`
	UserID string `json:"user_id"`
}

func serializeProductResponse(p *ent.Product) ProductResponse {
	return ProductResponse{
		ID:     p.ID.String(),
		Title:  p.Title,
		Price:  p.Price,
		UserID: p.UserID.String(),
	}
}
