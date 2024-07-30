package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/model"
)

func newDistributionCache() *gotest.Cache {
	record1 := &model.Distribution{}
	record1.ID = 1
	record2 := &model.Distribution{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewDistributionCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_distributionCache_Set(t *testing.T) {
	c := newDistributionCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Distribution)
	err := c.ICache.(DistributionCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(DistributionCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_distributionCache_Get(t *testing.T) {
	c := newDistributionCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Distribution)
	err := c.ICache.(DistributionCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(DistributionCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(DistributionCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_distributionCache_MultiGet(t *testing.T) {
	c := newDistributionCache()
	defer c.Close()

	var testData []*model.Distribution
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Distribution))
	}

	err := c.ICache.(DistributionCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(DistributionCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Distribution))
	}
}

func Test_distributionCache_MultiSet(t *testing.T) {
	c := newDistributionCache()
	defer c.Close()

	var testData []*model.Distribution
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Distribution))
	}

	err := c.ICache.(DistributionCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_distributionCache_Del(t *testing.T) {
	c := newDistributionCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Distribution)
	err := c.ICache.(DistributionCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_distributionCache_SetCacheWithNotFound(t *testing.T) {
	c := newDistributionCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Distribution)
	err := c.ICache.(DistributionCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDistributionCache(t *testing.T) {
	c := NewDistributionCache(&model.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewDistributionCache(&model.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewDistributionCache(&model.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
