package service_rest_req

type Order struct {
	AppUser   string `json:"app_user"`
	Amount    string `json:"amount"`
	EmbedData string `json:"embed_data"`
	Items     []Item `json:"items"`
}

type Item struct {
	ItemID    string  `json:"itemid"`    // ID của sản phẩm
	ItemName  string  `json:"itemname"`  // Tên sản phẩm
	ItemPrice float64 `json:"itemprice"` // Giá của sản phẩm
}
