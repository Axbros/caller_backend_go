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

var _ SmsHandler = (*smsHandler)(nil)

// SmsHandler defining the handler interface
type SmsHandler interface {
	Create(c *gin.Context)
	DeleteByMachineIdAndAddress(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
}

type smsHandler struct {
	iDao dao.SmsDao
}

// NewSmsHandler creating the handler interface
func NewSmsHandler() SmsHandler {
	return &smsHandler{
		iDao: dao.NewSmsDao(
			model.GetDB(),
			cache.NewSmsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create sms
// @Description submit information to create sms
// @Tags sms
// @accept json
// @Produce json
// @Param data body types.CreateSmsRequest true "sms information"
// @Success 200 {object} types.CreateSmsRespond{}
// @Router /api/v1/sms [post]
// @Security BearerAuth
func (h *smsHandler) Create(c *gin.Context) {
	form := &types.CreateSmsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	sms := &model.Sms{}
	err = copier.Copy(sms, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateSms)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, sms)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": sms.ID})
}

// DeleteByID delete a record by id
// @Summary delete sms
// @Description delete sms by id
// @Tags sms
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteSmsByIDRespond{}
// @Router /api/v1/sms/{id} [delete]
// @Security BearerAuth
func (h *smsHandler) DeleteByMachineIdAndAddress(c *gin.Context) {
	machineID := getMachineIDFromPath(c)
	address := getAddressFromPath(c)
	if machineID == "" {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByMachineIdAndAddress(ctx, machineID, address)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("machineID", machineID), logger.Any("address", address), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update sms
// @Description update sms information by id
// @Tags sms
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateSmsByIDRequest true "sms information"
// @Success 200 {object} types.UpdateSmsByIDRespond{}
// @Router /api/v1/sms/{id} [put]
// @Security BearerAuth
func (h *smsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getSmsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateSmsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	sms := &model.Sms{}
	err = copier.Copy(sms, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDSms)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, sms)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get sms detail
// @Description get sms detail by id
// @Tags sms
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetSmsByIDRespond{}
// @Router /api/v1/sms/{id} [get]
// @Security BearerAuth
func (h *smsHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getSmsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	sms, err := h.iDao.GetByID(ctx, id)
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

	data := &types.SmsObjDetail{}
	err = copier.Copy(data, sms)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDSms)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = idStr

	response.Success(c, gin.H{"sms": data})
}

// List of records by query parameters
// @Summary list of smss by query parameters
// @Description list of smss by paging and conditions
// @Tags sms
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListSmssRespond{}
// @Router /api/v1/sms/list [post]
// @Security BearerAuth
func (h *smsHandler) List(c *gin.Context) {
	form := &types.ListSmssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	smss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertSmss(smss)
	if err != nil {
		response.Error(c, ecode.ErrListSms)
		return
	}

	response.Success(c, gin.H{
		"smss":  data,
		"total": total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete smss
// @Description delete smss by batch id
// @Tags sms
// @Param data body types.DeleteSmssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteSmssByIDsRespond{}
// @Router /api/v1/sms/delete/ids [post]
// @Security BearerAuth
func (h *smsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteSmssByIDsRequest{}
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
// @Summary get sms by condition
// @Description get sms by condition
// @Tags sms
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetSmsByConditionRespond{}
// @Router /api/v1/sms/condition [post]
// @Security BearerAuth
func (h *smsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetSmsByConditionRequest{}
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
	sms, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.SmsObjDetail{}
	err = copier.Copy(data, sms)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDSms)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(sms.ID)

	response.Success(c, gin.H{"sms": data})
}

// ListByIDs list of records by batch id
// @Summary list of smss by batch id
// @Description list of smss by batch id
// @Tags sms
// @Param data body types.ListSmssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListSmssByIDsRespond{}
// @Router /api/v1/sms/list/ids [post]
// @Security BearerAuth
func (h *smsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListSmssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	smsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	smss := []*types.SmsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := smsMap[id]; ok {
			record, err := convertSms(v)
			if err != nil {
				response.Error(c, ecode.ErrListSms)
				return
			}
			smss = append(smss, record)
		}
	}

	response.Success(c, gin.H{
		"smss": smss,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of smss by last id and limit
// @Description list of smss by last id and limit
// @Tags sms
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListSmssRespond{}
// @Router /api/v1/sms/list [get]
// @Security BearerAuth
func (h *smsHandler) ListByLastID(c *gin.Context) {
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
	smss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertSmss(smss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDSms)
		return
	}

	response.Success(c, gin.H{
		"smss": data,
	})
}

func getSmsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func getMachineIDFromPath(c *gin.Context) string {
	idStr := c.Param("machine_id")
	return idStr
}

func getAddressFromPath(c *gin.Context) string {
	idStr := c.Param("address")
	return idStr
}

func convertSms(sms *model.Sms) (*types.SmsObjDetail, error) {
	data := &types.SmsObjDetail{}
	err := copier.Copy(data, sms)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here
	data.ID = utils.Uint64ToStr(sms.ID)
	return data, nil
}

func convertSmss(fromValues []*model.Sms) ([]*types.SmsObjDetail, error) {
	toValues := []*types.SmsObjDetail{}
	for _, v := range fromValues {
		data, err := convertSms(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
