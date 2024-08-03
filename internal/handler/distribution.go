package handler

import (
	"errors"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/cache"
	"caller/internal/dao"
	"caller/internal/ecode"
	"caller/internal/model"
	"caller/internal/types"
)

var _ DistributionHandler = (*distributionHandler)(nil)

// DistributionHandler defining the handler interface
type DistributionHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
}

type distributionHandler struct {
	iDao dao.DistributionDao
}

// NewDistributionHandler creating the handler interface
func NewDistributionHandler() DistributionHandler {
	return &distributionHandler{
		iDao: dao.NewDistributionDao(
			model.GetDB(),
			cache.NewDistributionCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create distribution
// @Description submit information to create distribution
// @Tags distribution
// @accept json
// @Produce json
// @Param data body types.CreateDistributionRequest true "distribution information"
// @Success 200 {object} types.CreateDistributionRespond{}
// @Router /api/v1/distribution [post]
// @Security BearerAuth
func (h *distributionHandler) Create(c *gin.Context) {
	form := &types.CreateDistributionRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	distribution := &model.Distribution{}
	err = copier.Copy(distribution, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateDistribution)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, distribution)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": distribution.ID})
}

// DeleteByID delete a record by id
// @Summary delete distribution
// @Description delete distribution by id
// @Tags distribution
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteDistributionByIDRespond{}
// @Router /api/v1/distribution/{id} [delete]
// @Security BearerAuth
func (h *distributionHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getDistributionIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update distribution
// @Description update distribution information by id
// @Tags distribution
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateDistributionByIDRequest true "distribution information"
// @Success 200 {object} types.UpdateDistributionByIDRespond{}
// @Router /api/v1/distribution/{id} [put]
// @Security BearerAuth
func (h *distributionHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getDistributionIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateDistributionByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	distribution := &model.Distribution{}
	err = copier.Copy(distribution, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDDistribution)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, distribution)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get distribution detail
// @Description get distribution detail by id
// @Tags distribution
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetDistributionByIDRespond{}
// @Router /api/v1/distribution/{id} [get]
// @Security BearerAuth
func (h *distributionHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getDistributionIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	distribution, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.DistributionObjDetail{}
	err = copier.Copy(data, distribution)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDDistribution)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = idStr

	response.Success(c, gin.H{"distribution": data})
}

// List of records by query parameters
// @Summary list of distributions by query parameters
// @Description list of distributions by paging and conditions
// @Tags distribution
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListDistributionsRespond{}
// @Router /api/v1/distribution/list [post]
// @Security BearerAuth
func (h *distributionHandler) List(c *gin.Context) {
	form := &types.ListDistributionsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	distributions, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertDistributions(distributions)
	if err != nil {
		response.Error(c, ecode.ErrListDistribution)
		return
	}

	response.Success(c, gin.H{
		"distributions": data,
		"total":         total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete distributions
// @Description delete distributions by batch id
// @Tags distribution
// @Param data body types.DeleteDistributionsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteDistributionsByIDsRespond{}
// @Router /api/v1/distribution/delete/ids [post]
// @Security BearerAuth
func (h *distributionHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteDistributionsByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.DeleteByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByCondition get a record by condition
// @Summary get distribution by condition
// @Description get distribution by condition
// @Tags distribution
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetDistributionByConditionRespond{}
// @Router /api/v1/distribution/condition [post]
// @Security BearerAuth
func (h *distributionHandler) GetByCondition(c *gin.Context) {
	form := &types.GetDistributionByConditionRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	err = form.Conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	distribution, err := h.iDao.GetByCondition(ctx, &form.Conditions)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByCondition not found", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByCondition error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.DistributionObjDetail{}
	err = copier.Copy(data, distribution)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDDistribution)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(distribution.ID)

	response.Success(c, gin.H{"distribution": data})
}

// ListByIDs list of records by batch id
// @Summary list of distributions by batch id
// @Description list of distributions by batch id
// @Tags distribution
// @Param data body types.ListDistributionsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListDistributionsByIDsRespond{}
// @Router /api/v1/distribution/list/ids [post]
// @Security BearerAuth
func (h *distributionHandler) ListByIDs(c *gin.Context) {
	form := &types.ListDistributionsByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	distributionMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	distributions := []*types.DistributionObjDetail{}
	for _, id := range form.IDs {
		if v, ok := distributionMap[id]; ok {
			record, err := convertDistribution(v)
			if err != nil {
				response.Error(c, ecode.ErrListDistribution)
				return
			}
			distributions = append(distributions, record)
		}
	}

	response.Success(c, gin.H{
		"distributions": distributions,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of distributions by last id and limit
// @Description list of distributions by last id and limit
// @Tags distribution
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListDistributionsRespond{}
// @Router /api/v1/distribution/list [get]
// @Security BearerAuth
func (h *distributionHandler) ListByLastID(c *gin.Context) {
	lastID := utils.StrToUint64(c.Query("lastID"))
	if lastID == 0 {
		lastID = math.MaxInt32
	}
	limit := utils.StrToInt(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	sort := c.Query("sort")

	ctx := middleware.WrapCtx(c)
	distributions, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertDistributions(distributions)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDDistribution)
		return
	}

	response.Success(c, gin.H{
		"distributions": data,
	})
}

func getDistributionIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertDistribution(distribution *model.Distribution) (*types.DistributionObjDetail, error) {
	data := &types.DistributionObjDetail{}
	err := copier.Copy(data, distribution)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(distribution.ID)
	return data, nil
}

func convertDistributions(fromValues []*model.Distribution) ([]*types.DistributionObjDetail, error) {
	toValues := []*types.DistributionObjDetail{}
	for _, v := range fromValues {
		data, err := convertDistribution(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
