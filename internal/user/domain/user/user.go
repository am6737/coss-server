package user

type User struct {
	ID           string
	CossID       string
	Email        string
	Tel          string
	NickName     string
	Avatar       string
	PublicKey    string
	Password     string
	LastIp       string
	LineIp       string
	CreatedIp    string
	Signature    string
	LineAt       int64
	LastAt       int64
	Status       UserStatus
	EmailVerity  uint
	Bot          uint
	SecretBundle string
	CreatedAt    int64
	UpdatedAt    int64
	DeletedAt    int64
}

type UserStatus uint

const (
	// UserStatusNormal 正常状态
	UserStatusNormal UserStatus = iota + 1
	// UserStatusDisable 禁用状态
	UserStatusDisable
	// UserStatusDeleted 删除状态
	UserStatusDeleted
	// UserStatusLock 锁定状态
	UserStatusLock
)

// UserActivateStatus 用户激活状态
type UserActivateStatus uint

const (
	NotActivated UserActivateStatus = iota //未激活
	Activated                              //已经激活
)

type UserLogin struct {
	ID          uint
	DriveID     uint
	CreatedAt   int64
	UpdatedAt   int64
	DeletedAt   int64
	UserId      string
	LoginCount  uint
	LastAt      int64
	Token       string
	DriverId    string
	DriverToken string
	DriverType  string
	Platform    string
	ClientIP    string
	Rid         string
}

//type BaseModel struct {
//	ID        uint
//	CreatedAt int64
//	UpdatedAt int64
//	DeletedAt int64
//}
