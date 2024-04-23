package adapters

import (
	"context"
	"fmt"
	"github.com/cossim/coss-server/internal/group/cache"
	"github.com/cossim/coss-server/internal/group/domain/group"
	ptime "github.com/cossim/coss-server/pkg/utils/time"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
)

type BaseModel struct {
	ID        uint32 `gorm:"primaryKey;autoIncrement;"`
	CreatedAt int64  `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt int64  `gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt int64  `gorm:"default:0;comment:删除时间"`
}

type GroupModel struct {
	BaseModel
	Type            uint   `gorm:"default:0;comment:群聊类型(0=私密群, 1=公开群)"`
	Status          uint   `gorm:"comment:群聊状态(0=正常状态, 1=锁定状态, 2=删除状态)"`
	MaxMembersLimit int    `gorm:"comment:群聊人数限制"`
	CreatorID       string `gorm:"type:varchar(64);comment:创建者id"`
	Name            string `gorm:"comment:群聊名称"`
	Avatar          string `gorm:"default:'';comment:头像（群）"`
}

func (bm *BaseModel) TableName() string {
	fmt.Println("table name")
	return "groups"
}

func (bm *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := ptime.Now()
	bm.CreatedAt = now
	bm.UpdatedAt = now
	return nil
}

func (bm *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	bm.UpdatedAt = ptime.Now()
	return nil
}

func (m *GroupModel) FromEntity(e *group.Group) error {
	if m == nil {
		m = &GroupModel{}
	}
	if err := e.Validate(); err != nil {
		return err
	}
	m.ID = e.ID
	m.CreatedAt = e.CreatedAt
	m.Type = uint(e.Type)
	m.Status = uint(e.Status)
	m.MaxMembersLimit = e.MaxMembersLimit
	m.CreatorID = e.CreatorID
	m.Name = e.Name
	m.Avatar = e.Avatar
	return nil
}

func (m *GroupModel) ToEntity() (*group.Group, error) {
	if m == nil {
		return nil, errors.New("group model is nil")
	}
	return &group.Group{
		ID:              m.ID,
		CreatedAt:       m.CreatedAt,
		Type:            group.Type(m.Type),
		Status:          group.Status(m.Status),
		MaxMembersLimit: m.MaxMembersLimit,
		CreatorID:       m.CreatorID,
		Name:            m.Name,
		Avatar:          m.Avatar,
	}, nil
}

const mySQLDeadlockErrorCode = 1213

var _ group.Repository = &MySQLGroupRepository{}

func NewMySQLGroupRepository(db *gorm.DB, cache cache.GroupCache) *MySQLGroupRepository {
	return &MySQLGroupRepository{
		db:    db,
		cache: cache,
	}
}

type MySQLGroupRepository struct {
	db    *gorm.DB
	cache cache.GroupCache
}

func (m *MySQLGroupRepository) Automigrate() error {
	return m.db.AutoMigrate(&GroupModel{})
}

func (m *MySQLGroupRepository) UpdateFields(ctx context.Context, id uint32, fields map[string]interface{}) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Unscoped().Updates(fields).Error
}

func (m *MySQLGroupRepository) Get(ctx context.Context, id uint32) (*group.Group, error) {
	cachedGroup, err := m.cache.GetGroup(ctx, id)
	if err == nil && cachedGroup != nil {
		return cachedGroup, nil

	}

	model := &GroupModel{}
	if err := m.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}

	entity, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	if err := m.cache.SetGroup(ctx, entity); err != nil {
		log.Println("Error caching group:", err)
	}

	return model.ToEntity()
}

func (m *MySQLGroupRepository) Update(ctx context.Context, group *group.Group, updateFn func(h *group.Group) (*group.Group, error)) error {
	model := &GroupModel{}
	if err := model.FromEntity(group); err != nil {
		return err
	}
	if err := m.db.WithContext(ctx).Save(model).Error; err != nil {
		return err
	}

	if updateFn == nil {
		return nil
	}

	_, err := updateFn(group)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQLGroupRepository) Create(ctx context.Context, group *group.Group, createFn func(h *group.Group) (*group.Group, error)) error {
	for {
		_, err := m.create(ctx, group, createFn)

		if val, ok := errors.Cause(err).(*mysql.MySQLError); ok && val.Number == mySQLDeadlockErrorCode {
			continue
		}

		return err
	}
}

func (m *MySQLGroupRepository) create(ctx context.Context, group *group.Group, createFn func(h *group.Group) (*group.Group, error)) (*group.Group, error) {
	model := &GroupModel{}
	if err := model.FromEntity(group); err != nil {
		return nil, err
	}

	if err := m.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	entity, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	return createFn(entity)
}

func (m *MySQLGroupRepository) Delete(ctx context.Context, id uint32) error {
	err := m.db.WithContext(ctx).Delete(&GroupModel{}, id).Error
	return err
}

// Find retrieves groups based on the provided query.
func (m *MySQLGroupRepository) Find(ctx context.Context, query group.Query) ([]*group.Group, error) {
	if m.cache == nil {
		return m.findWithoutCache(ctx, query)
	}

	cacheKey := m.generateCacheKey(query)
	cachedGroups, err := m.cache.GetGroups(ctx, cacheKey)
	if err == nil && cachedGroups != nil {
		return cachedGroups, nil
	}

	groups, err := m.findWithoutCache(ctx, query)
	if err != nil {
		return nil, err
	}

	if err := m.cache.SetGroup(ctx, groups...); err != nil {
		log.Println("Error caching groups:", err)
	}

	return groups, nil
}

// findWithoutCache executes the query directly without using cache.
func (m *MySQLGroupRepository) findWithoutCache(ctx context.Context, query group.Query) ([]*group.Group, error) {
	// Perform query directly on the database
	var models []*GroupModel

	db := m.db.WithContext(ctx)

	// Apply filters
	if len(query.ID) > 0 {
		db = db.Where("id IN (?)", query.ID)
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if len(query.UserID) > 0 {
		db = db.Joins("JOIN group_members ON groups.id = group_members.group_id").
			Where("group_members.user_id IN (?)", query.UserID)
	}
	if query.CreateAt != nil {
		db = db.Where("created_at >= ?", query.CreateAt)
	}
	if query.UpdateAt != nil {
		db = db.Where("updated_at >= ?", query.UpdateAt)
	}
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	if query.Offset > 0 {
		db = db.Offset(query.Offset)
	}

	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}

	// Convert models to entities
	groups := make([]*group.Group, len(models))
	for i, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		groups[i] = entity
	}

	return groups, nil
}

// generateCacheKey generates a cache key based on the query parameters.
func (m *MySQLGroupRepository) generateCacheKey(query group.Query) []uint32 {
	// Example: concatenate query parameters to form a cache key
	// This is a simplistic approach; you might need to customize it based on your specific requirements
	cacheKey := make([]uint32, 0, len(query.ID))
	cacheKey = append(cacheKey, query.ID...)
	return cacheKey
}