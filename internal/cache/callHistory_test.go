package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/model"
)

func newCallHistoryCache() *gotest.Cache {
	record1 := &model.CallHistory{}
	record1.ID = 1
	record2 := &model.CallHistory{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewCallHistoryCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_callHistoryCache_Set(t *testing.T) {
	c := newCallHistoryCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.CallHistory)
	err := c.ICache.(CallHistoryCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(CallHistoryCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_callHistoryCache_Get(t *testing.T) {
	c := newCallHistoryCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.CallHistory)
	err := c.ICache.(CallHistoryCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CallHistoryCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(CallHistoryCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_callHistoryCache_MultiGet(t *testing.T) {
	c := newCallHistoryCache()
	defer c.Close()

	var testData []*model.CallHistory
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.CallHistory))
	}

	err := c.ICache.(CallHistoryCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(CallHistoryCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.CallHistory))
	}
}

func Test_callHistoryCache_MultiSet(t *testing.T) {
	c := newCallHistoryCache()
	defer c.Close()

	var testData []*model.CallHistory
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.CallHistory))
	}

	err := c.ICache.(CallHistoryCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_callHistoryCache_Del(t *testing.T) {
	c := newCallHistoryCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.CallHistory)
	err := c.ICache.(CallHistoryCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_callHistoryCache_SetCacheWithNotFound(t *testing.T) {
	c := newCallHistoryCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.CallHistory)
	err := c.ICache.(CallHistoryCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCallHistoryCache(t *testing.T) {
	c := NewCallHistoryCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewCallHistoryCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewCallHistoryCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
