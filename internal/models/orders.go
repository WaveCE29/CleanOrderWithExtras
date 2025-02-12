package models

type InputOrder struct {
	No                int    `json:"no"`
	PlatformProductId string `json:"platformProductId"`
	Qty               int    `json:"qty"`
	UnitPrice         int    `json:"unitPrice"`
	TotalPrice        int    `json:"totalPrice"`
}

type CleanedOrder struct {
	No         int    `json:"no"`
	ProductId  string `json:"productId"`
	MaterialId string `json:"materialId,omitempty"`
	ModelId    string `json:"modelId,omitempty"`
	Qty        int    `json:"qty"`
	UnitPrice  int    `json:"unitPrice"`
	TotalPrice int    `json:"totalPrice"`
}
