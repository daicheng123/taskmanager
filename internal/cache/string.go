package cache

/*
	string类型数据操作集合
*/

import (
	"context"
	"taskmanager/internal/cache/utils"
	"taskmanager/internal/consts"
	"taskmanager/pkg/logger"
	"taskmanager/pkg/store"
	"time"
)

type StringOperation struct {
	ctx context.Context
}

func NewStringOperation() *StringOperation {
	return &StringOperation{
		ctx: context.Background(),
	}
}

// Exists key是否存在
func (so StringOperation) Exists(key string) *utils.Result {
	rc, err := store.GetCacheOperator()
	if err != nil {
		logger.Error(err.Error())
		return utils.NewResult(nil, err)
	}
	return utils.NewResult(rc.Exists(so.ctx, key).Result())

}

// Get 获取单值
func (so StringOperation) Get(key string) *utils.Result {
	rc, err := store.GetCacheOperator()
	if err != nil {
		logger.Error(err.Error())
		return utils.NewResult(nil, err)
	}
	return utils.NewResult(rc.Get(so.ctx, key).Result())
}

// Set 设置单个值
func (so StringOperation) Set(key string, value interface{}, attrs ...*utils.Attr) *utils.Result {
	rc, err := store.GetCacheOperator()
	if err != nil {
		logger.Error(err.Error())
		return utils.NewResult(nil, err)
	}

	expire := utils.Attrs(attrs).Find(consts.AttrExpire)
	if expire == nil {
		expire = 0
	}

	return utils.NewResult(rc.Set(so.ctx, key, value, expire.(time.Duration)).Result())
}

//
////MGet 获取多个值 获取单值
//func (so StringOperation) MGet(keys ...string) *SliceResult {
//	co, err := store.GetCacheOperator()
//	return utils.NewStringResult(.(ths.ctx, keys...).Result())
//}
//
//func (so StringOperation) Set(key string, value interface{}, attrs ...*OperatorAttr) *StringResult {
//	expire := OperatorAttrs(attrs).Find(ATTR_EXPIRE)
//	if expire == nil {
//		expire = 0
//	}
//	if nx := OperatorAttrs(attrs).Find(ATTR_NX); nx != nil {
//		return NewStringResult(Redis().SetNX(ths.ctx, key, value, expire.(time.Duration)).Result())
//	}
//	return NewStringResult(Redis().Set(ths.ctx, key, value, expire.(time.Duration)).Result())
//}
