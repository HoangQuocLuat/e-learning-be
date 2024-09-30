package kafka_event

type SendMailEvent struct {
	Emails []string `json:"emails"`
}

func (a *SendMailEvent) GetId() string {
	return "id_event"
}
