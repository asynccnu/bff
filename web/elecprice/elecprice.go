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

func (h *ElecPriceHandler) RegisterRoutes(s *gin.Engine, authMiddleware gin.HandlerFunc) {
	sg := s.Group("/elecprice")
	sg.GET("/check", authMiddleware, ginx.WrapClaimsAndReq(h.Check))
	sg.POST("/setStandard", authMiddleware, ginx.WrapClaimsAndReq(h.SetStandard))

}

// Check 查询电费
// @Summary 查询电费
// @Description 根据区域、楼栋和房间号查询电费信息
// @Tags 电费
// @Param area query string true "区域，例如 '东区学生宿舍'"
// @Param building query string true "楼栋，例如 '1号楼'"
// @Param room query string true "房间号，例如 '101'"
// @Produce json
// @Success 200 {object} web.Response{data=CheckResponse} "成功返回电费信息"
// @Failure 500 {object} web.Response{msg=string} "系统异常"
// @Router /elecprice/check [get]
func (h *ElecPriceHandler) Check(ctx *gin.Context, req CheckRequest, uc ijwt.UserClaims) (web.Response, error) {
	checkresponse, err := h.ElecPriceClient.Check(ctx, &elecpricev1.CheckRequest{
		Place: &elecpricev1.Place{
			Area:     req.Area,
			Building: req.Building,
			Room:     req.Room,
		},
	})
	if err != nil {
		return web.Response{}, errs.ELECPRICE_CHECK_ERROR(err)
	}

	return web.Response{
		Data: CheckResponse{
			Price: &Price{
				LightingRemainMoney:       checkresponse.Price.LightingRemainMoney,
				LightingYesterdayUseValue: checkresponse.Price.LightingYesterdayUseValue,
				LightingYesterdayUseMoney: checkresponse.Price.LightingYesterdayUseMoney,
				AirRemainMoney:            checkresponse.Price.AirRemainMoney,
				AirYesterdayUseValue:      checkresponse.Price.AirYesterdayUseValue,
				AirYesterdayUseMoney:      checkresponse.Price.AirYesterdayUseMoney,
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
// @Param request body SetStandardRequest true "设置电费提醒请求参数"
// @Success 200 {object} web.Response{msg=string} "设置成功的返回信息"
// @Failure 500 {object} web.Response{msg=string} "系统异常"
// @Router /elecprice/setStandard [post]
func (h *ElecPriceHandler) SetStandard(ctx *gin.Context, req SetStandardRequest, uc ijwt.UserClaims) (web.Response, error) {

	_, err := h.ElecPriceClient.SetStandard(ctx, &elecpricev1.SetStandardRequest{
		StudentId: uc.StudentId,
		Money:     req.Money,
		Place: &elecpricev1.Place{
			Area:     req.Area,
			Building: req.Building,
			Room:     req.Room,
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
