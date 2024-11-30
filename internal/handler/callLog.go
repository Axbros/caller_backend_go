package handler

import (
	"errors"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
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

var _ UnanswerdCallHandler = (*callLogHandler)(nil)

// UnanswerdCallHandler defining the handler interface
type UnanswerdCallHandler interface {
	Create(c *gin.Context)
	MultipleCreate(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	DeleteAll(c *gin.Context)
	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
	GetByUserID(c *gin.Context)
	ReadOfflineMessage(c *gin.Context)
}

type callLogHandler struct {
	iDao dao.UnanswerdCallDao
	dDao dao.DistributionDao
}

// NewUnanswerdCallHandler creating the handler interface
func NewUnanswerdCallHandler() UnanswerdCallHandler {
	return &callLogHandler{
		iDao: dao.NewUnanswerdCallDao(
			model.GetDB(),
			cache.NewUnanswerdCallCache(model.GetCacheType()),
		),
		dDao: dao.NewDistributionDao(
			model.GetDB(),
			cache.NewDistributionCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create callLog
// @Description submit information to create callLog
// @Tags callLog
// @accept json
// @Produce json
// @Param data body types.CreateUnanswerdCallRequest true "callLog information"
// @Success 200 {object} types.CreateUnanswerdCallRespond{}
// @Router /api/v1/callLog [post]
// @Security BearerAuth
func (h *callLogHandler) Create(c *gin.Context) {
	form := &types.CreateUnanswerdCallRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	callLog := &model.UnanswerdCall{}
	err = copier.Copy(callLog, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUnanswerdCall)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, callLog)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": callLog.ID})
}

// DeleteByID delete a record by id
// @Summary delete callLog
// @Description delete callLog by id
// @Tags callLog
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUnanswerdCallByIDRespond{}
// @Router /api/v1/callLog/{id} [delete]
// @Security BearerAuth
func (h *callLogHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUnanswerdCallIDFromPath(c)
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
// @Summary update callLog
// @Description update callLog information by id
// @Tags callLog
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUnanswerdCallByIDRequest true "callLog information"
// @Success 200 {object} types.UpdateUnanswerdCallByIDRespond{}
// @Router /api/v1/callLog/{id} [put]
// @Security BearerAuth
func (h *callLogHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUnanswerdCallIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateUnanswerdCallByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	callLog := &model.UnanswerdCall{}
	err = copier.Copy(callLog, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDUnanswerdCall)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, callLog)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get callLog detail
// @Description get callLog detail by id
// @Tags callLog
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUnanswerdCallByIDRespond{}
// @Router /api/v1/callLog/{id} [get]
// @Security BearerAuth
func (h *callLogHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getUnanswerdCallIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	callLog, err := h.iDao.GetByID(ctx, id)
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

	data := &types.UnanswerdCallObjDetail{}
	err = copier.Copy(data, callLog)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUnanswerdCall)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = idStr

	response.Success(c, gin.H{"callLog": data})
}

// List of records by query parameters
// @Summary list of callLogs by query parameters
// @Description list of callLogs by paging and conditions
// @Tags callLog
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUnanswerdCallsRespond{}
// @Router /api/v1/callLog/list [post]
// @Security BearerAuth
func (h *callLogHandler) List(c *gin.Context) {
	form := &types.ListUnanswerdCallsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	callLogs, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUnanswerdCalls(callLogs)
	if err != nil {
		response.Error(c, ecode.ErrListUnanswerdCall)
		return
	}

	response.Success(c, gin.H{
		"callLogs": data,
		"total":    total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete callLogs
// @Description delete callLogs by batch id
// @Tags callLog
// @Param data body types.DeleteUnanswerdCallsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteUnanswerdCallsByIDsRespond{}
// @Router /api/v1/callLog/delete/ids [post]
// @Security BearerAuth
func (h *callLogHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteUnanswerdCallsByIDsRequest{}
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

func (h *callLogHandler) DeleteAll(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteAll(ctx)
	if err != nil {
		logger.Error("Delete All error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c)
}

// GetByCondition get a record by condition
// @Summary get callLog by condition
// @Description get callLog by condition
// @Tags callLog
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUnanswerdCallByConditionRespond{}
// @Router /api/v1/callLog/condition [post]
// @Security BearerAuth
func (h *callLogHandler) GetByCondition(c *gin.Context) {
	form := &types.GetUnanswerdCallByConditionRequest{}
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
	callLog, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &[]types.UnanswerdCallObjDetail{}
	err = copier.Copy(data, callLog)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUnanswerdCall)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	// data.ID = utils.Uint64ToStr(callLog)

	response.Success(c, gin.H{"callLog": data})
}

// ListByIDs list of records by batch id
// @Summary list of callLogs by batch id
// @Description list of callLogs by batch id
// @Tags callLog
// @Param data body types.ListUnanswerdCallsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListUnanswerdCallsByIDsRespond{}
// @Router /api/v1/callLog/list/ids [post]
// @Security BearerAuth
func (h *callLogHandler) ListByIDs(c *gin.Context) {
	form := &types.ListUnanswerdCallsByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	callLogMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	callLogs := []*types.UnanswerdCallObjDetail{}
	for _, id := range form.IDs {
		if v, ok := callLogMap[id]; ok {
			record, err := convertUnanswerdCall(v)
			if err != nil {
				response.Error(c, ecode.ErrListUnanswerdCall)
				return
			}
			callLogs = append(callLogs, record)
		}
	}

	response.Success(c, gin.H{
		"callLogs": callLogs,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of callLogs by last id and limit
// @Description list of callLogs by last id and limit
// @Tags callLog
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListUnanswerdCallsRespond{}
// @Router /api/v1/callLog/list [get]
// @Security BearerAuth
func (h *callLogHandler) ListByLastID(c *gin.Context) {
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
	callLogs, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUnanswerdCalls(callLogs)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDUnanswerdCall)
		return
	}

	response.Success(c, gin.H{
		"callLogs": data,
	})
}
func (h *callLogHandler) GetByUserID(c *gin.Context) {
	callLogType := c.Query("type")

	userID, _, _ := getUnanswerdCallIDFromPath(c)
	children, _ := h.iDao.GetChildrenByUserID(c, userID)
	res := []*model.UnanswerdCall{}
	var records []*model.UnanswerdCall
	var childrenList []int
	for i := 0; i < len(children); i++ {
		childrenList = append(childrenList, children[i].ClientID)
		if callLogType == "keypad" {
			records, _ = h.iDao.GetByCondition(c, &query.Conditions{
				Columns: []query.Column{
					{
						Name:  "client_id",
						Value: children[i].ClientID,
					},
				},
			})
		} else {
			records, _ = h.iDao.GetByCondition(c, &query.Conditions{
				Columns: []query.Column{
					{
						Name:  "client_id",
						Value: children[i].ClientID,
					},
					{
						Name:  "type",
						Value: callLogType,
					},
				},
			})
		}

		for _, record := range records {
			res = append(res, record)
		}
	}
	response.Success(c, gin.H{
		"children": childrenList,
		"data":     res,
	})
}

func (h *callLogHandler) MultipleCreate(c *gin.Context) {
	recordType := c.Param("type")
	// var machineId string
	// machineId = c.Request.Header.Get("machine_id")
	// body := c.Request.Body
	json := make(map[string]string) //注意该结构接受的内容
	c.BindJSON(&json)
	// fmt.Printf("%v", &json)
	unansweredList := []model.UnanswerdCall{}
	for _, value := range json {
		record := model.UnanswerdCall{

			MobileNumber: value,
			Type:         recordType,
		}
		unansweredList = append(unansweredList, record)
	}
	h.iDao.CreateMultiple(c, &unansweredList)

	response.Success(c, gin.H{
		"status": "ok",
	})
}

func getUnanswerdCallIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertUnanswerdCall(callLog *model.UnanswerdCall) (*types.UnanswerdCallObjDetail, error) {
	data := &types.UnanswerdCallObjDetail{}
	err := copier.Copy(data, callLog)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(callLog.ID)
	return data, nil
}

func convertUnanswerdCalls(fromValues []*model.UnanswerdCall) ([]*types.UnanswerdCallObjDetail, error) {
	toValues := []*types.UnanswerdCallObjDetail{}
	for _, v := range fromValues {
		data, err := convertUnanswerdCall(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}

func (h *callLogHandler) ReadOfflineMessage(c *gin.Context) {
	userID := c.Param("user_id")
	setOfflineMsgUnread(c, userID, "false")
	response.Success(c, gin.H{
		"status": "ok",
	})
}
