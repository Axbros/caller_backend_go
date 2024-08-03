package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/model"
)

func newGroupClientCache() *gotest.Cache {
	record1 := &model.GroupClient{}
	record1.ID = 1
	record2 := &model.GroupClient{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewGroupClientCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_groupClientCache_Set(t *testing.T) {
	c := newGroupClientCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.GroupClient)
	err := c.ICache.(GroupClientCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(GroupClientCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_groupClientCache_Get(t *testing.T) {
	c := newGroupClientCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.GroupClient)
	err := c.ICache.(GroupClientCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(GroupClientCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(GroupClientCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_groupClientCache_MultiGet(t *testing.T) {
	c := newGroupClientCache()
	defer c.Close()

	var testData []*model.GroupClient
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.GroupClient))
	}

	err := c.ICache.(GroupClientCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(GroupClientCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.GroupClient))
	}
}

func Test_groupClientCache_MultiSet(t *testing.T) {
	c := newGroupClientCache()
	defer c.Close()

	var testData []*model.GroupClient
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.GroupClient))
	}

	err := c.ICache.(GroupClientCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_groupClientCache_Del(t *testing.T) {
	c := newGroupClientCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.GroupClient)
	err := c.ICache.(GroupClientCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_groupClientCache_SetCacheWithNotFound(t *testing.T) {
	c := newGroupClientCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.GroupClient)
	err := c.ICache.(GroupClientCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewGroupClientCache(t *testing.T) {
	c := NewGroupClientCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewGroupClientCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewGroupClientCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
