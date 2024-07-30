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

var _ CallHistoryHandler = (*callHistoryHandler)(nil)

// CallHistoryHandler defining the handler interface
type CallHistoryHandler interface {
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

type callHistoryHandler struct {
	iDao dao.CallHistoryDao
}

// NewCallHistoryHandler creating the handler interface
func NewCallHistoryHandler() CallHistoryHandler {
	return &callHistoryHandler{
		iDao: dao.NewCallHistoryDao(
			model.GetDB(),
			cache.NewCallHistoryCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create callHistory
// @Description submit information to create callHistory
// @Tags callHistory
// @accept json
// @Produce json
// @Param data body types.CreateCallHistoryRequest true "callHistory information"
// @Success 200 {object} types.CreateCallHistoryRespond{}
// @Router /api/v1/callHistory [post]
// @Security BearerAuth
func (h *callHistoryHandler) Create(c *gin.Context) {
	form := &types.CreateCallHistoryRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	callHistory := &model.CallHistory{}
	err = copier.Copy(callHistory, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateCallHistory)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, callHistory)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": callHistory.ID})
}

// DeleteByID delete a record by id
// @Summary delete callHistory
// @Description delete callHistory by id
// @Tags callHistory
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteCallHistoryByIDRespond{}
// @Router /api/v1/callHistory/{id} [delete]
// @Security BearerAuth
func (h *callHistoryHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getCallHistoryIDFromPath(c)
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
// @Summary update callHistory
// @Description update callHistory information by id
// @Tags callHistory
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateCallHistoryByIDRequest true "callHistory information"
// @Success 200 {object} types.UpdateCallHistoryByIDRespond{}
// @Router /api/v1/callHistory/{id} [put]
// @Security BearerAuth
func (h *callHistoryHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getCallHistoryIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateCallHistoryByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	callHistory := &model.CallHistory{}
	err = copier.Copy(callHistory, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDCallHistory)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, callHistory)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get callHistory detail
// @Description get callHistory detail by id
// @Tags callHistory
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetCallHistoryByIDRespond{}
// @Router /api/v1/callHistory/{id} [get]
// @Security BearerAuth
func (h *callHistoryHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getCallHistoryIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	callHistory, err := h.iDao.GetByID(ctx, id)
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

	data := &types.CallHistoryObjDetail{}
	err = copier.Copy(data, callHistory)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDCallHistory)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = idStr

	response.Success(c, gin.H{"callHistory": data})
}

// List of records by query parameters
// @Summary list of callHistorys by query parameters
// @Description list of callHistorys by paging and conditions
// @Tags callHistory
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListCallHistorysRespond{}
// @Router /api/v1/callHistory/list [post]
// @Security BearerAuth
func (h *callHistoryHandler) List(c *gin.Context) {
	form := &types.ListCallHistorysRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	callHistorys, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertCallHistorys(callHistorys)
	if err != nil {
		response.Error(c, ecode.ErrListCallHistory)
		return
	}

	response.Success(c, gin.H{
		"callHistorys": data,
		"total":        total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete callHistorys
// @Description delete callHistorys by batch id
// @Tags callHistory
// @Param data body types.DeleteCallHistorysByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteCallHistorysByIDsRespond{}
// @Router /api/v1/callHistory/delete/ids [post]
// @Security BearerAuth
func (h *callHistoryHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteCallHistorysByIDsRequest{}
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
// @Summary get callHistory by condition
// @Description get callHistory by condition
// @Tags callHistory
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetCallHistoryByConditionRespond{}
// @Router /api/v1/callHistory/condition [post]
// @Security BearerAuth
func (h *callHistoryHandler) GetByCondition(c *gin.Context) {
	form := &types.GetCallHistoryByConditionRequest{}
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
	callHistory, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.CallHistoryObjDetail{}
	err = copier.Copy(data, callHistory)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDCallHistory)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(callHistory.ID)

	response.Success(c, gin.H{"callHistory": data})
}

// ListByIDs list of records by batch id
// @Summary list of callHistorys by batch id
// @Description list of callHistorys by batch id
// @Tags callHistory
// @Param data body types.ListCallHistorysByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListCallHistorysByIDsRespond{}
// @Router /api/v1/callHistory/list/ids [post]
// @Security BearerAuth
func (h *callHistoryHandler) ListByIDs(c *gin.Context) {
	form := &types.ListCallHistorysByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	callHistoryMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	callHistorys := []*types.CallHistoryObjDetail{}
	for _, id := range form.IDs {
		if v, ok := callHistoryMap[id]; ok {
			record, err := convertCallHistory(v)
			if err != nil {
				response.Error(c, ecode.ErrListCallHistory)
				return
			}
			callHistorys = append(callHistorys, record)
		}
	}

	response.Success(c, gin.H{
		"callHistorys": callHistorys,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of callHistorys by last id and limit
// @Description list of callHistorys by last id and limit
// @Tags callHistory
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListCallHistorysRespond{}
// @Router /api/v1/callHistory/list [get]
// @Security BearerAuth
func (h *callHistoryHandler) ListByLastID(c *gin.Context) {
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
	callHistorys, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertCallHistorys(callHistorys)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDCallHistory)
		return
	}

	response.Success(c, gin.H{
		"callHistorys": data,
	})
}

func getCallHistoryIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertCallHistory(callHistory *model.CallHistory) (*types.CallHistoryObjDetail, error) {
	data := &types.CallHistoryObjDetail{}
	err := copier.Copy(data, callHistory)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(callHistory.ID)
	return data, nil
}

func convertCallHistorys(fromValues []*model.CallHistory) ([]*types.CallHistoryObjDetail, error) {
	toValues := []*types.CallHistoryObjDetail{}
	for _, v := range fromValues {
		data, err := convertCallHistory(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
