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

var _ GroupCallHandler = (*groupCallHandler)(nil)

// GroupCallHandler defining the handler interface
type GroupCallHandler interface {
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

type groupCallHandler struct {
	iDao dao.GroupCallDao
}

// NewGroupCallHandler creating the handler interface
func NewGroupCallHandler() GroupCallHandler {
	return &groupCallHandler{
		iDao: dao.NewGroupCallDao(
			model.GetDB(),
			cache.NewGroupCallCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create groupCall
// @Description submit information to create groupCall
// @Tags groupCall
// @accept json
// @Produce json
// @Param data body types.CreateGroupCallRequest true "groupCall information"
// @Success 200 {object} types.CreateGroupCallRespond{}
// @Router /api/v1/groupCall [post]
// @Security BearerAuth
func (h *groupCallHandler) Create(c *gin.Context) {
	form := &types.CreateGroupCallRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	groupCall := &model.GroupCall{}
	err = copier.Copy(groupCall, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateGroupCall)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, groupCall)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": groupCall.ID})
}

// DeleteByID delete a record by id
// @Summary delete groupCall
// @Description delete groupCall by id
// @Tags groupCall
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteGroupCallByIDRespond{}
// @Router /api/v1/groupCall/{id} [delete]
// @Security BearerAuth
func (h *groupCallHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getGroupCallIDFromPath(c)
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
// @Summary update groupCall
// @Description update groupCall information by id
// @Tags groupCall
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateGroupCallByIDRequest true "groupCall information"
// @Success 200 {object} types.UpdateGroupCallByIDRespond{}
// @Router /api/v1/groupCall/{id} [put]
// @Security BearerAuth
func (h *groupCallHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getGroupCallIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateGroupCallByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	groupCall := &model.GroupCall{}
	err = copier.Copy(groupCall, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDGroupCall)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, groupCall)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get groupCall detail
// @Description get groupCall detail by id
// @Tags groupCall
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetGroupCallByIDRespond{}
// @Router /api/v1/groupCall/{id} [get]
// @Security BearerAuth
func (h *groupCallHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getGroupCallIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	groupCall, err := h.iDao.GetByID(ctx, id)
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

	data := &types.GroupCallObjDetail{}
	err = copier.Copy(data, groupCall)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDGroupCall)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = idStr

	response.Success(c, gin.H{"groupCall": data})
}

// List of records by query parameters
// @Summary list of groupCalls by query parameters
// @Description list of groupCalls by paging and conditions
// @Tags groupCall
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListGroupCallsRespond{}
// @Router /api/v1/groupCall/list [post]
// @Security BearerAuth
func (h *groupCallHandler) List(c *gin.Context) {
	form := &types.ListGroupCallsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	groupCalls, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertGroupCalls(groupCalls)
	if err != nil {
		response.Error(c, ecode.ErrListGroupCall)
		return
	}

	response.Success(c, gin.H{
		"groupCalls": data,
		"total":      total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete groupCalls
// @Description delete groupCalls by batch id
// @Tags groupCall
// @Param data body types.DeleteGroupCallsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteGroupCallsByIDsRespond{}
// @Router /api/v1/groupCall/delete/ids [post]
// @Security BearerAuth
func (h *groupCallHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteGroupCallsByIDsRequest{}
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
// @Summary get groupCall by condition
// @Description get groupCall by condition
// @Tags groupCall
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetGroupCallByConditionRespond{}
// @Router /api/v1/groupCall/condition [post]
// @Security BearerAuth
func (h *groupCallHandler) GetByCondition(c *gin.Context) {
	form := &types.GetGroupCallByConditionRequest{}
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
	groupCall, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.GroupCallObjDetail{}
	err = copier.Copy(data, groupCall)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDGroupCall)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(groupCall.ID)

	response.Success(c, gin.H{"groupCall": data})
}

// ListByIDs list of records by batch id
// @Summary list of groupCalls by batch id
// @Description list of groupCalls by batch id
// @Tags groupCall
// @Param data body types.ListGroupCallsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListGroupCallsByIDsRespond{}
// @Router /api/v1/groupCall/list/ids [post]
// @Security BearerAuth
func (h *groupCallHandler) ListByIDs(c *gin.Context) {
	form := &types.ListGroupCallsByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	groupCallMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	groupCalls := []*types.GroupCallObjDetail{}
	for _, id := range form.IDs {
		if v, ok := groupCallMap[id]; ok {
			record, err := convertGroupCall(v)
			if err != nil {
				response.Error(c, ecode.ErrListGroupCall)
				return
			}
			groupCalls = append(groupCalls, record)
		}
	}

	response.Success(c, gin.H{
		"groupCalls": groupCalls,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of groupCalls by last id and limit
// @Description list of groupCalls by last id and limit
// @Tags groupCall
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListGroupCallsRespond{}
// @Router /api/v1/groupCall/list [get]
// @Security BearerAuth
func (h *groupCallHandler) ListByLastID(c *gin.Context) {
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
	groupCalls, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertGroupCalls(groupCalls)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDGroupCall)
		return
	}

	response.Success(c, gin.H{
		"groupCalls": data,
	})
}

func getGroupCallIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertGroupCall(groupCall *model.GroupCall) (*types.GroupCallObjDetail, error) {
	data := &types.GroupCallObjDetail{}
	err := copier.Copy(data, groupCall)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(groupCall.ID)
	return data, nil
}

func convertGroupCalls(fromValues []*model.GroupCall) ([]*types.GroupCallObjDetail, error) {
	toValues := []*types.GroupCallObjDetail{}
	for _, v := range fromValues {
		data, err := convertGroupCall(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
