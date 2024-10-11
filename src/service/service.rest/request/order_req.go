package service_rest_req

type Order struct {
	AppID       string                 `json:"app_id"`
	AppUser     string                 `json:"app_user"`
	AppTime     string                 `json:"app_time"`
	Amount      string                 `json:"amount"`
	AppTransID  string                 `json:"app_trans_id"`
	BankCode    string                 `json:"bank_code"`
	EmbedData   string                 `json:"embed_data"`
	Items       map[string]interface{} `json:"items"`
	Description string                 `json:"description"`
	Mac         string                 `json:"mac"`
}
