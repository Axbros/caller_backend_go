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

var _ GroupClientHandler = (*groupClientHandler)(nil)

// GroupClientHandler defining the handler interface
type GroupClientHandler interface {
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

type groupClientHandler struct {
	iDao dao.GroupClientDao
}

// NewGroupClientHandler creating the handler interface
func NewGroupClientHandler() GroupClientHandler {
	return &groupClientHandler{
		iDao: dao.NewGroupClientDao(
			model.GetDB(),
			cache.NewGroupClientCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create groupClient
// @Description submit information to create groupClient
// @Tags groupClient
// @accept json
// @Produce json
// @Param data body types.CreateGroupClientRequest true "groupClient information"
// @Success 200 {object} types.CreateGroupClientRespond{}
// @Router /api/v1/groupClient [post]
// @Security BearerAuth
func (h *groupClientHandler) Create(c *gin.Context) {
	form := &types.CreateGroupClientRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	groupClient := &model.GroupClient{}
	err = copier.Copy(groupClient, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateGroupClient)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, groupClient)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": groupClient.ID})
}

// DeleteByID delete a record by id
// @Summary delete groupClient
// @Description delete groupClient by id
// @Tags groupClient
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteGroupClientByIDRespond{}
// @Router /api/v1/groupClient/{id} [delete]
// @Security BearerAuth
func (h *groupClientHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getGroupClientIDFromPath(c)
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
// @Summary update groupClient
// @Description update groupClient information by id
// @Tags groupClient
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateGroupClientByIDRequest true "groupClient information"
// @Success 200 {object} types.UpdateGroupClientByIDRespond{}
// @Router /api/v1/groupClient/{id} [put]
// @Security BearerAuth
func (h *groupClientHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getGroupClientIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateGroupClientByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	groupClient := &model.GroupClient{}
	err = copier.Copy(groupClient, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDGroupClient)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, groupClient)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get groupClient detail
// @Description get groupClient detail by id
// @Tags groupClient
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetGroupClientByIDRespond{}
// @Router /api/v1/groupClient/{id} [get]
// @Security BearerAuth
func (h *groupClientHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getGroupClientIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	groupClient, err := h.iDao.GetByID(ctx, id)
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

	data := &types.GroupClientObjDetail{}
	err = copier.Copy(data, groupClient)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDGroupClient)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = idStr

	response.Success(c, gin.H{"groupClient": data})
}

// List of records by query parameters
// @Summary list of groupClients by query parameters
// @Description list of groupClients by paging and conditions
// @Tags groupClient
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListGroupClientsRespond{}
// @Router /api/v1/groupClient/list [post]
// @Security BearerAuth
func (h *groupClientHandler) List(c *gin.Context) {
	form := &types.ListGroupClientsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	groupClients, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertGroupClients(groupClients)
	if err != nil {
		response.Error(c, ecode.ErrListGroupClient)
		return
	}

	response.Success(c, gin.H{
		"groupClients": data,
		"total":        total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete groupClients
// @Description delete groupClients by batch id
// @Tags groupClient
// @Param data body types.DeleteGroupClientsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteGroupClientsByIDsRespond{}
// @Router /api/v1/groupClient/delete/ids [post]
// @Security BearerAuth
func (h *groupClientHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteGroupClientsByIDsRequest{}
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
// @Summary get groupClient by condition
// @Description get groupClient by condition
// @Tags groupClient
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetGroupClientByConditionRespond{}
// @Router /api/v1/groupClient/condition [post]
// @Security BearerAuth
func (h *groupClientHandler) GetByCondition(c *gin.Context) {
	form := &types.GetGroupClientByConditionRequest{}
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
	groupClient, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.GroupClientObjDetail{}
	err = copier.Copy(data, groupClient)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDGroupClient)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(groupClient.ID)

	response.Success(c, gin.H{"groupClient": data})
}

// ListByIDs list of records by batch id
// @Summary list of groupClients by batch id
// @Description list of groupClients by batch id
// @Tags groupClient
// @Param data body types.ListGroupClientsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListGroupClientsByIDsRespond{}
// @Router /api/v1/groupClient/list/ids [post]
// @Security BearerAuth
func (h *groupClientHandler) ListByIDs(c *gin.Context) {
	form := &types.ListGroupClientsByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	groupClientMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	groupClients := []*types.GroupClientObjDetail{}
	for _, id := range form.IDs {
		if v, ok := groupClientMap[id]; ok {
			record, err := convertGroupClient(v)
			if err != nil {
				response.Error(c, ecode.ErrListGroupClient)
				return
			}
			groupClients = append(groupClients, record)
		}
	}

	response.Success(c, gin.H{
		"groupClients": groupClients,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of groupClients by last id and limit
// @Description list of groupClients by last id and limit
// @Tags groupClient
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListGroupClientsRespond{}
// @Router /api/v1/groupClient/list [get]
// @Security BearerAuth
func (h *groupClientHandler) ListByLastID(c *gin.Context) {
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
	groupClients, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertGroupClients(groupClients)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDGroupClient)
		return
	}

	response.Success(c, gin.H{
		"groupClients": data,
	})
}

func getGroupClientIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertGroupClient(groupClient *model.GroupClient) (*types.GroupClientObjDetail, error) {
	data := &types.GroupClientObjDetail{}
	err := copier.Copy(data, groupClient)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(groupClient.ID)
	return data, nil
}

func convertGroupClients(fromValues []*model.GroupClient) ([]*types.GroupClientObjDetail, error) {
	toValues := []*types.GroupClientObjDetail{}
	for _, v := range fromValues {
		data, err := convertGroupClient(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
