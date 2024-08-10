package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/model"
)

func newUnanswerdCallCache() *gotest.Cache {
	record1 := &model.UnanswerdCall{}
	record1.ID = 1
	record2 := &model.UnanswerdCall{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewUnanswerdCallCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_callLogCache_Set(t *testing.T) {
	c := newUnanswerdCallCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UnanswerdCall)
	err := c.ICache.(UnanswerdCallCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(UnanswerdCallCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_callLogCache_Get(t *testing.T) {
	c := newUnanswerdCallCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UnanswerdCall)
	err := c.ICache.(UnanswerdCallCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UnanswerdCallCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(UnanswerdCallCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_callLogCache_MultiGet(t *testing.T) {
	c := newUnanswerdCallCache()
	defer c.Close()

	var testData []*model.UnanswerdCall
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UnanswerdCall))
	}

	err := c.ICache.(UnanswerdCallCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UnanswerdCallCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.UnanswerdCall))
	}
}

func Test_callLogCache_MultiSet(t *testing.T) {
	c := newUnanswerdCallCache()
	defer c.Close()

	var testData []*model.UnanswerdCall
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UnanswerdCall))
	}

	err := c.ICache.(UnanswerdCallCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_callLogCache_Del(t *testing.T) {
	c := newUnanswerdCallCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UnanswerdCall)
	err := c.ICache.(UnanswerdCallCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_callLogCache_SetCacheWithNotFound(t *testing.T) {
	c := newUnanswerdCallCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UnanswerdCall)
	err := c.ICache.(UnanswerdCallCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewUnanswerdCallCache(t *testing.T) {
	c := NewUnanswerdCallCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewUnanswerdCallCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewUnanswerdCallCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
