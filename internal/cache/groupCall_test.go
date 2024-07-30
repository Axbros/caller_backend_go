package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/model"
)

func newGroupCallCache() *gotest.Cache {
	record1 := &model.GroupCall{}
	record1.ID = 1
	record2 := &model.GroupCall{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewGroupCallCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_groupCallCache_Set(t *testing.T) {
	c := newGroupCallCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.GroupCall)
	err := c.ICache.(GroupCallCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(GroupCallCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_groupCallCache_Get(t *testing.T) {
	c := newGroupCallCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.GroupCall)
	err := c.ICache.(GroupCallCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(GroupCallCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(GroupCallCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_groupCallCache_MultiGet(t *testing.T) {
	c := newGroupCallCache()
	defer c.Close()

	var testData []*model.GroupCall
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.GroupCall))
	}

	err := c.ICache.(GroupCallCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(GroupCallCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.GroupCall))
	}
}

func Test_groupCallCache_MultiSet(t *testing.T) {
	c := newGroupCallCache()
	defer c.Close()

	var testData []*model.GroupCall
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.GroupCall))
	}

	err := c.ICache.(GroupCallCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_groupCallCache_Del(t *testing.T) {
	c := newGroupCallCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.GroupCall)
	err := c.ICache.(GroupCallCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_groupCallCache_SetCacheWithNotFound(t *testing.T) {
	c := newGroupCallCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.GroupCall)
	err := c.ICache.(GroupCallCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewGroupCallCache(t *testing.T) {
	c := NewGroupCallCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewGroupCallCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewGroupCallCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
