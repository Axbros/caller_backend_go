package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/cache"
	"caller/internal/model"
)

var _ UnanswerdCallDao = (*callLogDao)(nil)

// UnanswerdCallDao defining the dao interface
type UnanswerdCallDao interface {
	Create(ctx context.Context, table *model.UnanswerdCall) error
	CreateMultiple(ctx context.Context, table *[]model.UnanswerdCall) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.UnanswerdCall) error
	GetByID(ctx context.Context, id uint64) (*model.UnanswerdCall, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.UnanswerdCall, int64, error)

	DeleteByIDs(ctx context.Context, ids []uint64) error
	DeleteAll(ctx context.Context) error
	GetByCondition(ctx context.Context, condition *query.Conditions) ([]*model.UnanswerdCall, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.UnanswerdCall, error)
	GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.UnanswerdCall, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.UnanswerdCall) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.UnanswerdCall) error
	GetChildrenByUserID(ctx context.Context, UserID string) ([]model.GroupClient, error)
}

type callLogDao struct {
	db    *gorm.DB
	cache cache.UnanswerdCallCache // if nil, the cache is not used.
	sfg   *singleflight.Group      // if cache is nil, the sfg is not used.
}

// NewUnanswerdCallDao creating the dao interface
func NewUnanswerdCallDao(db *gorm.DB, xCache cache.UnanswerdCallCache) UnanswerdCallDao {
	if xCache == nil {
		return &callLogDao{db: db}
	}
	return &callLogDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *callLogDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a record, insert the record and the id value is written back to the table
func (d *callLogDao) Create(ctx context.Context, table *model.UnanswerdCall) error {
	return d.db.WithContext(ctx).Create(table).Error
}
func (d *callLogDao) CreateMultiple(ctx context.Context, table *[]model.UnanswerdCall) error {
	for _, call := range *table {
		var existingCall model.UnanswerdCall
		result := d.db.Where("mobile_number =? AND client_time =?", call.MobileNumber, call.ClientTime).First(&existingCall)
		if result.Error == nil {
			// 如果找到相同的记录，跳过当前数据的添加
			// fmt.Printf("已存在相同的 mobile_number: %s 和 client_time: %s 的记录，跳过添加\n", call.MobileNumber, call.ClientTime)
			continue
		} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果未找到相同记录，执行添加操作
			if err := d.db.Create(&call).Error; err != nil {
				return err
			}
		} else {
			// 其他错误情况
			return result.Error
		}
	}

	return nil

}

// DeleteByID delete a record by id
func (d *callLogDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UnanswerdCall{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a record by id
func (d *callLogDao) UpdateByID(ctx context.Context, table *model.UnanswerdCall) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *callLogDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.UnanswerdCall) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.ClientMachineCode != "" {
		update["client_id"] = table.ClientMachineCode
	}
	if table.MobileNumber != "" {
		update["mobile_number"] = table.MobileNumber
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *callLogDao) GetByID(ctx context.Context, id uint64) (*model.UnanswerdCall, error) {
	// no cache
	if d.cache == nil {
		record := &model.UnanswerdCall{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache or database
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, model.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.UnanswerdCall{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				// if data is empty, set not found cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, model.ErrRecordNotFound) {
					err = d.cache.SetCacheWithNotFound(ctx, id)
					if err != nil {
						return nil, err
					}
					return nil, model.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			err = d.cache.Set(ctx, id, table, cache.UnanswerdCallExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.UnanswerdCall)
		if !ok {
			return nil, model.ErrRecordNotFound
		}
		return table, nil
	} else if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// fail fast, if cache error return, don't request to db
	return nil, err
}

// GetByColumns get paging records by column information,
// Note: query performance degrades when table rows are very large because of the use of offset.
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	size: lines per page
//	sort: sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: search for a male over 20 years of age
//
//	params = &query.Params{
//	    Page: 0,
//	    Size: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *callLogDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.UnanswerdCall, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.UnanswerdCall{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.UnanswerdCall{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// DeleteByIDs delete records by batch id
func (d *callLogDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.UnanswerdCall{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.deleteCache(ctx, id)
	}

	return nil
}

func (d *callLogDao) DeleteAll(ctx context.Context) error {
	err := d.db.WithContext(ctx).Where("id > (?)", 0).Delete(&model.UnanswerdCall{}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetByCondition get a record by condition
// query conditions:
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like, in
//	value: column value, if exp=in, multiple values are separated by commas
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: find a male aged 20
//
//	condition = &query.Conditions{
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *callLogDao) GetByCondition(ctx context.Context, c *query.Conditions) ([]*model.UnanswerdCall, error) {
	queryStr, args, err := c.ConvertToGorm()
	if err != nil {
		return nil, err
	}

	table := []*model.UnanswerdCall{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).Find(&table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs get records by batch id
func (d *callLogDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.UnanswerdCall, error) {
	// no cache
	if d.cache == nil {
		var records []*model.UnanswerdCall
		err := d.db.WithContext(ctx).Where("id IN (?)", ids).Find(&records).Error
		if err != nil {
			return nil, err
		}
		itemMap := make(map[uint64]*model.UnanswerdCall)
		for _, record := range records {
			itemMap[record.ID] = record
		}
		return itemMap, nil
	}

	// get form cache or database
	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedIDs []uint64
	for _, id := range ids {
		_, ok := itemMap[id]
		if !ok {
			missedIDs = append(missedIDs, id)
			continue
		}
	}

	// get missed data
	if len(missedIDs) > 0 {
		// find the id of an active placeholder, i.e. an id that does not exist in database
		var realMissedIDs []uint64
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if errors.Is(err, cacheBase.ErrPlaceholder) {
				continue
			}
			realMissedIDs = append(realMissedIDs, id)
		}

		if len(realMissedIDs) > 0 {
			var missedData []*model.UnanswerdCall
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				for _, data := range missedData {
					itemMap[data.ID] = data
				}
				err = d.cache.MultiSet(ctx, missedData, cache.UnanswerdCallExpireTime)
				if err != nil {
					return nil, err
				}
			} else {
				for _, id := range realMissedIDs {
					_ = d.cache.SetCacheWithNotFound(ctx, id)
				}
			}
		}
	}

	return itemMap, nil
}

// GetByLastID get paging records by last id and limit
func (d *callLogDao) GetByLastID(ctx context.Context, lastID uint64, limit int, sort string) ([]*model.UnanswerdCall, error) {
	page := query.NewPage(0, limit, sort)

	records := []*model.UnanswerdCall{}
	err := d.db.WithContext(ctx).Order(page.Sort()).Limit(page.Size()).Where("id < ?", lastID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// CreateByTx create a record in the database using the provided transaction
func (d *callLogDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.UnanswerdCall) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *callLogDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.UnanswerdCall{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *callLogDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.UnanswerdCall) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
func (d *callLogDao) GetChildrenByUserID(ctx context.Context, UserID string) ([]model.GroupClient, error) {
	var distributionRecord model.Distribution
	err := d.db.WithContext(ctx).Model(&model.Distribution{}).Where("user_id =?", UserID).First(&distributionRecord).Error
	if err != nil {
		return nil, err
	}

	var groupCallRecord model.GroupCall
	err = d.db.WithContext(ctx).Model(&model.GroupCall{}).Where("id =?", distributionRecord.GroupCallID).First(&groupCallRecord).Error
	if err != nil {
		return nil, err
	}

	var groupClientRecord []model.GroupClient
	err = d.db.WithContext(ctx).Model(&model.GroupClient{}).Where("group_name =?", groupCallRecord.GroupName).Find(&groupClientRecord).Error
	if err != nil {
		return nil, err
	}

	return groupClientRecord, nil
}