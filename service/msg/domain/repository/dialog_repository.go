package repository

import (
	"github.com/cossim/coss-server/service/msg/domain/entity"
)

type DialogRepository interface {
	CreateDialog(ownerID string, dialogType entity.DialogType, groupID uint) (*entity.Dialog, error)
	JoinDialog(dialogID uint, userID string) (*entity.DialogUser, error)
	GetUserDialogs(userID string) ([]uint, error)
	GetDialogsByIDs(dialogIDs []uint) ([]*entity.Dialog, error)
}