package entity

type Dialog struct {
	BaseModel
	OwnerId string     `gorm:"comment:用户id" json:"owner_id"`
	Type    DialogType `gorm:"comment:对话类型" json:"type"`
	GroupId uint       `gorm:"comment:群组id" json:"group_id"`
}

type BaseModel struct {
	ID        uint  `gorm:"primaryKey;autoIncrement;" json:"id"`
	CreatedAt int64 `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt int64 `gorm:"default:0;comment:删除时间" json:"deleted_at"`
}

type DialogType uint32

const (
	UserDialog = iota
	GroupDialog
)
