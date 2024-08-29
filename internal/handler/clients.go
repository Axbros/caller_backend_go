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

var _ ClientsHandler = (*clientsHandler)(nil)

// ClientsHandler defining the handler interface
type ClientsHandler interface {
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

type clientsHandler struct {
	iDao dao.ClientsDao
}

// NewClientsHandler creating the handler interface
func NewClientsHandler() ClientsHandler {
	return &clientsHandler{
		iDao: dao.NewClientsDao(
			model.GetDB(),
			cache.NewClientsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create clients
// @Description submit information to create clients
// @Tags clients
// @accept json
// @Produce json
// @Param data body types.CreateClientsRequest true "clients information"
// @Success 200 {object} types.CreateClientsRespond{}
// @Router /api/v1/clients [post]
// @Security BearerAuth
func (h *clientsHandler) Create(c *gin.Context) {
	form := &types.CreateClientsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	clients := &model.Clients{}
	err = copier.Copy(clients, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateClients)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)

	recordID := h.iDao.IsExist(ctx, clients)

	if recordID > 0 {
		response.Success(c, gin.H{"id": recordID})
		return
	}
	err = h.iDao.Create(ctx, clients)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": clients.ID})
}

// DeleteByID delete a record by id
// @Summary delete clients
// @Description delete clients by id
// @Tags clients
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteClientsByIDRespond{}
// @Router /api/v1/clients/{id} [delete]
// @Security BearerAuth
func (h *clientsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getClientsIDFromPath(c)
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
// @Summary update clients
// @Description update clients information by id
// @Tags clients
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateClientsByIDRequest true "clients information"
// @Success 200 {object} types.UpdateClientsByIDRespond{}
// @Router /api/v1/clients/{id} [put]
// @Security BearerAuth
func (h *clientsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getClientsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateClientsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	clients := &model.Clients{}
	err = copier.Copy(clients, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDClients)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, clients)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get clients detail
// @Description get clients detail by id
// @Tags clients
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetClientsByIDRespond{}
// @Router /api/v1/clients/{id} [get]
// @Security BearerAuth
func (h *clientsHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getClientsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	clients, err := h.iDao.GetByID(ctx, id)
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

	data := &types.ClientsObjDetail{}
	err = copier.Copy(data, clients)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDClients)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = idStr

	response.Success(c, gin.H{"clients": data})
}

// List of records by query parameters
// @Summary list of clientss by query parameters
// @Description list of clientss by paging and conditions
// @Tags clients
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListClientssRespond{}
// @Router /api/v1/clients/list [post]
// @Security BearerAuth
func (h *clientsHandler) List(c *gin.Context) {
	form := &types.ListClientssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	clientss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertClientss(clientss)
	if err != nil {
		response.Error(c, ecode.ErrListClients)
		return
	}

	response.Success(c, gin.H{
		"clientss": data,
		"total":    total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete clientss
// @Description delete clientss by batch id
// @Tags clients
// @Param data body types.DeleteClientssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteClientssByIDsRespond{}
// @Router /api/v1/clients/delete/ids [post]
// @Security BearerAuth
func (h *clientsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteClientssByIDsRequest{}
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
// @Summary get clients by condition
// @Description get clients by condition
// @Tags clients
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetClientsByConditionRespond{}
// @Router /api/v1/clients/condition [post]
// @Security BearerAuth
func (h *clientsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetClientsByConditionRequest{}
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
	clients, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.ClientsObjDetail{}
	err = copier.Copy(data, clients)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDClients)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(clients.ID)

	response.Success(c, gin.H{"clients": data})
}

// ListByIDs list of records by batch id
// @Summary list of clientss by batch id
// @Description list of clientss by batch id
// @Tags clients
// @Param data body types.ListClientssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListClientssByIDsRespond{}
// @Router /api/v1/clients/list/ids [post]
// @Security BearerAuth
func (h *clientsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListClientssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	clientsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	clientss := []*types.ClientsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := clientsMap[id]; ok {
			record, err := convertClients(v)
			if err != nil {
				response.Error(c, ecode.ErrListClients)
				return
			}
			clientss = append(clientss, record)
		}
	}

	response.Success(c, gin.H{
		"clientss": clientss,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of clientss by last id and limit
// @Description list of clientss by last id and limit
// @Tags clients
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListClientssRespond{}
// @Router /api/v1/clients/list [get]
// @Security BearerAuth
func (h *clientsHandler) ListByLastID(c *gin.Context) {
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
	clientss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertClientss(clientss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDClients)
		return
	}

	response.Success(c, gin.H{
		"clientss": data,
	})
}

func getClientsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertClients(clients *model.Clients) (*types.ClientsObjDetail, error) {
	data := &types.ClientsObjDetail{}
	err := copier.Copy(data, clients)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(clients.ID)
	return data, nil
}

func convertClientss(fromValues []*model.Clients) ([]*types.ClientsObjDetail, error) {
	toValues := []*types.ClientsObjDetail{}
	for _, v := range fromValues {
		data, err := convertClients(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
