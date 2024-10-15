package service_rest_req

type CheckPay struct {
	AppTransID string `json:"app_trans_id"`
	UserID     string `json:"user_id"`
	TuitionID  string `json:"tuition_id"`
}
