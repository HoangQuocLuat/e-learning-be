package src_const

const (
	StatusCodeInvalid = 601

	StatusCodeMongoCreateError = 701
	StatusCodeMongoReadError   = 702
	StatusCodeMongoUpdateError = 703
	StatusCodeMongoDeleteError = 704
)

const (
	CreateErr     = "1-"
	ReadErr       = "2-"
	UpdateErr     = "3-"
	DeleteErr     = "4-"
	InvalidErr    = "5-"
	ExistedErr    = "6-"
	NotExistedErr = "7-"
)

const (
	ElementErr_Account = "1-"
)

const (
	InternalError = "100"
	Unauthorized  = "99"
	GrpcError     = "98"

	AccountExist      = "1.1"
	WrongPassword     = "1.2"
	AccountNotFound   = "1.3"
	AccountNotActive  = "1.4"
	UnableCreateToken = "1.5"
)

const (
	ServiceErr_Auth       = "1"
	ServiceErr_E_Learning = "2"
)
