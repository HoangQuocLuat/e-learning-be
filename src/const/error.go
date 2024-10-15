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
	ElementErr_User      = "1-"
	ElementErr_Infor     = "2-"
	ElementErr_Class     = "3-"
	ElementErr_Schedules = "4-"
	ElementErr_Tuition   = "5-"
)

const (
	InternalError = "100"
	Unauthorized  = "99"
	GrpcError     = "98"

	AccountExist      = "1.1"
	WrongPassword     = "1.2"
	UserNotFound      = "1.3"
	UserNotActive     = "1.4"
	UnableCreateToken = "1.5"

	ClassExist = "2.1"
)

const (
	ServiceErr_Auth       = "1-"
	ServiceErr_E_Learning = "2-"
	ServiceErr_CronJob    = "3-"
)
