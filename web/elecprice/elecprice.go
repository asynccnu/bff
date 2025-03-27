package elecprice

import (
	elecpricev1 "github.com/asynccnu/be-api/gen/proto/elecprice/v1"
	"github.com/asynccnu/bff/errs"
	"github.com/asynccnu/bff/pkg/ginx"
	"github.com/asynccnu/bff/web"
	"github.com/asynccnu/bff/web/ijwt"
	"github.com/gin-gonic/gin"
)

type ElecPriceHandler struct {
	ElecPriceClient elecpricev1.ElecpriceServiceClient //注入的是grpc服务
	Administrators  map[string]struct{}                //这里注入的是管理员权限验证配置
}

func NewElecPriceHandler(elecPriceClient elecpricev1.ElecpriceServiceClient,
	administrators map[string]struct{}) *ElecPriceHandler {
	return &ElecPriceHandler{ElecPriceClient: elecPriceClient, Administrators: administrators}
}

func (h *ElecPriceHandler) RegisterRoutes(s *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	sg := s.Group("/elecprice")
	{
		sg.POST("/getAIDandName", authMiddleware, ginx.WrapClaimsAndReq(h.GetAIDandName))
		sg.POST("/setStandard", authMiddleware, ginx.WrapClaimsAndReq(h.SetStandard))
		sg.POST("/getRoomInfo", authMiddleware, ginx.WrapClaimsAndReq(h.GetRoomInfo))
		sg.POST("/getPrice", authMiddleware, ginx.WrapClaimsAndReq(h.GetPrice))
	}
}

// @Summary 获取楼栋和房间号
// @Description 通过区域获取楼栋和房间号
// @Tags 电费
// @Accept json
// @Produce json
// @Param request body elecprice.GetAIDandNameRequest true "设置电费提醒请求参数"
// @Success 200 {object} web.Response{msg=elecprice.GetAIDandNameResponse} "设置成功的返回信息"
// @Failure 500 {object} web.Response{msg=string} "系统异常"
// @Router /elecprice/getAIDandName [post]
func (h *ElecPriceHandler) GetAIDandName(ctx *gin.Context, req GetAIDandNameRequest, uc ijwt.UserClaims) (web.Response, error) {
	res, err := h.ElecPriceClient.GetAIDandName(ctx, &elecpricev1.GetAIDandNameRequest{
		AreaName: req.AreaName,
	})
	if err != nil {
		return web.Response{}, errs.ELECPRICE_SET_STANDARD_ERROR(err)
	}
	var architectureList []*Architecture
	for _, r := range res.ArchitectureList {
		architectureList = append(architectureList, &Architecture{
			ArchitectureName: r.ArchitectureName,
			ArchitectureID:   r.ArchitectureID,
		})
	}
	return web.Response{
		Data: GetAIDandNameResponse{
			ArchitectureList: architectureList,
		},
	}, nil
}

// @Summary 获取房间号和id
// @Description 根据房间号和空调/照明id
// @Tags 电费
// @Accept json
// @Produce json
// @Param request body elecprice.GetRoomInfoRequest true "获取楼栋信息请求参数"
// @Success 200 {object} web.Response{msg=elecprice.GetRoomInfoResponse} "获取成功的返回信息"
// @Failure 500 {object} web.Response{msg=string} "系统异常"
// @Router /elecprice/getRoomInfo [post]
func (h *ElecPriceHandler) GetRoomInfo(ctx *gin.Context, req GetRoomInfoRequest, uc ijwt.UserClaims) (web.Response, error) {
	res, err := h.ElecPriceClient.GetRoomInfo(ctx, &elecpricev1.GetRoomInfoRequest{
		ArchitectureID: req.ArchitectureID,
		Floor:          req.Floor,
	})
	if err != nil {
		return web.Response{}, errs.ELECPRICE_SET_STANDARD_ERROR(err)
	}
	var roomList []*Room
	for _, r := range res.RoomList {
		roomList = append(roomList, &Room{
			RoomID:   r.RoomID,
			RoomName: r.RoomName,
		})
	}
	return web.Response{
		Data: GetRoomInfoResponse{
			RoomList: roomList,
		},
	}, nil
}

// @Summary 获取电费
// @Description 根据房间号获取电费信息
// @Tags 电费
// @Accept json
// @Produce json
// @Param request body elecprice.GetPriceRequest true "获取电费请求参数"
// @Success 200 {object} web.Response{msg=elecprice.GetPriceResponse} "获取成功的返回信息"
// @Failure 500 {object} web.Response{msg=string} "系统异常"
// @Router /elecprice/getPrice [post]
func (h *ElecPriceHandler) GetPrice(ctx *gin.Context, req GetPriceRequest, uc ijwt.UserClaims) (web.Response, error) {
	res, err := h.ElecPriceClient.GetPrice(ctx, &elecpricev1.GetPriceRequest{
		RoomAircID:  req.RoomAircID,
		RoomLightID: req.RoomLightID,
	})
	if err != nil {
		return web.Response{}, errs.ELECPRICE_SET_STANDARD_ERROR(err)
	}
	return web.Response{
		Data: GetPriceResponse{
			Price: &Price{
				LightingRemainMoney:       res.Price.LightingRemainMoney,
				LightingYesterdayUseValue: res.Price.LightingYesterdayUseValue,
				LightingYesterdayUseMoney: res.Price.LightingYesterdayUseMoney,
				AirRemainMoney:            res.Price.AirRemainMoney,
				AirYesterdayUseValue:      res.Price.AirYesterdayUseValue,
				AirYesterdayUseMoney:      res.Price.AirYesterdayUseMoney,
			},
		},
	}, nil
}

// SetStandard 设置电费
// @Summary 设置电费提醒标准
// @Description 根据区域、楼栋和房间号设置电费提醒的金额标准
// @Tags 电费
// @Accept json
// @Produce json
// @Param request body elecprice.SetStandardRequest true "设置电费提醒请求参数"
// @Success 200 {object} web.Response{msg=string} "设置成功的返回信息"
// @Failure 500 {object} web.Response{msg=string} "系统异常"
// @Router /elecprice/setStandard [post]
func (h *ElecPriceHandler) SetStandard(ctx *gin.Context, req SetStandardRequest, uc ijwt.UserClaims) (web.Response, error) {

	_, err := h.ElecPriceClient.SetStandard(ctx, &elecpricev1.SetStandardRequest{
		StudentId: uc.StudentId,
		Money:     req.Money,
		Ids: &elecpricev1.SetStandardRequest_IDs{
			RoomAircID:  req.AirID,
			RoomLightID: req.LightID,
		},
	})
	if err != nil {
		return web.Response{}, errs.ELECPRICE_SET_STANDARD_ERROR(err)
	}

	return web.Response{
		Msg: "设置电费提醒标准成功!",
	}, nil
}

// 这个方法是用来检查是否是管理员的,虽然我觉得写的不好把管理员写死在代码里面了,但是先凑合着用,之后再改,具体怎么使用请随便去别的服务里面找一找
func (h *ElecPriceHandler) isAdmin(studentId string) bool {
	_, exists := h.Administrators[studentId]
	return exists
}
