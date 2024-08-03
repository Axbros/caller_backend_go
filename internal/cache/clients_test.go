package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/model"
)

func newClientsCache() *gotest.Cache {
	record1 := &model.Clients{}
	record1.ID = 1
	record2 := &model.Clients{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewClientsCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_clientsCache_Set(t *testing.T) {
	c := newClientsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Clients)
	err := c.ICache.(ClientsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(ClientsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_clientsCache_Get(t *testing.T) {
	c := newClientsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Clients)
	err := c.ICache.(ClientsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ClientsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(ClientsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_clientsCache_MultiGet(t *testing.T) {
	c := newClientsCache()
	defer c.Close()

	var testData []*model.Clients
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Clients))
	}

	err := c.ICache.(ClientsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ClientsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Clients))
	}
}

func Test_clientsCache_MultiSet(t *testing.T) {
	c := newClientsCache()
	defer c.Close()

	var testData []*model.Clients
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Clients))
	}

	err := c.ICache.(ClientsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_clientsCache_Del(t *testing.T) {
	c := newClientsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Clients)
	err := c.ICache.(ClientsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_clientsCache_SetCacheWithNotFound(t *testing.T) {
	c := newClientsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Clients)
	err := c.ICache.(ClientsCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewClientsCache(t *testing.T) {
	c := NewClientsCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewClientsCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewClientsCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
