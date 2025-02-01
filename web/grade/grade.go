package grade

import (
	"context"
	"fmt"
	counterv1 "github.com/asynccnu/be-api/gen/proto/counter/v1"
	gradev1 "github.com/asynccnu/be-api/gen/proto/grade/v1"
	"github.com/asynccnu/bff/errs"
	"github.com/asynccnu/bff/pkg/ginx"
	"github.com/asynccnu/bff/pkg/logger"
	"github.com/asynccnu/bff/web"
	"github.com/asynccnu/bff/web/ijwt"
	"github.com/gin-gonic/gin"
)

type GradeHandler struct {
	GradeClient    gradev1.GradeServiceClient //注入的是grpc服务
	CounterClient  counterv1.CounterServiceClient
	Administrators map[string]struct{} //这里注入的是管理员权限验证配置
	l              logger.Logger
}

func NewGradeHandler(
	GradeClient gradev1.GradeServiceClient, //注入的是grpc服务
	CounterClient counterv1.CounterServiceClient,
	l logger.Logger,
	administrators map[string]struct{}) *GradeHandler {
	return &GradeHandler{
		GradeClient:    GradeClient,
		CounterClient:  CounterClient,
		Administrators: administrators,
		l:              l,
	}
}

func (h *GradeHandler) RegisterRoutes(s *gin.Engine, authMiddleware gin.HandlerFunc) {
	sg := s.Group("/grade")
	//这里有三类路由,分别是ginx.WrapClaimsAndReq()有参数且要验证
	sg.GET("/getGradeByTerm", authMiddleware, ginx.WrapClaimsAndReq(h.GetGradeByTerm))
	sg.GET("/getGradeScore", authMiddleware, ginx.WrapClaims(h.GetGradeScore))

}

// GradeByTerm 查询按学年和学期的成绩
// @Summary 查询按学年和学期的成绩
// @Description 根据学年号和学期号获取用户的成绩
// @Tags 成绩
// @Accept json
// @Produce json
// @Param xnm query int true "学年号（如 2023 表示 2023~2024 学年）"
// @Param xqm query int true "学期号（1 表示第一学期，2 表示第二学期）"
// @Success 200 {object} web.Response{data=GetGradeByTermResp} "成功返回学年和学期的成绩信息"
// @Failure 500 {object} web.Response "系统异常，获取失败"
// @Router /grade/getGradeByTerm [get]
func (h *GradeHandler) GetGradeByTerm(ctx *gin.Context, req GetGradeByTermReq, uc ijwt.UserClaims) (web.Response, error) {
	grades, err := h.GradeClient.GetGradeByTerm(ctx, &gradev1.GetGradeByTermReq{
		StudentId: uc.StudentId,
		Xnm:       req.Xnm,
		Xqm:       req.Xqm,
	})
	if err != nil {
		return web.Response{}, errs.GET_GRADE_BY_TERM_ERROR(err)
	}

	var resp GetGradeByTermResp
	for _, grade := range grades.Grades {
		resp.Grades = append(resp.Grades, Grade{
			Kcmc:                grade.Kcmc,                // 课程名
			Xf:                  grade.Xf,                  // 学分
			Jd:                  grade.Jd,                  //绩点
			Cj:                  grade.Cj,                  // 总成绩
			Kcxzmc:              grade.Kcxzmc,              // 课程性质名称 比如专业主干课程/通识必修课
			Kclbmc:              grade.Kclbmc,              // 课程类别名称，比如专业课/公共课
			Kcbj:                grade.Kcbj,                // 课程标记，比如主修/辅修
			RegularGradePercent: grade.RegularGradePercent, // 平时分占比
			RegularGrade:        grade.RegularGrade,        // 平时分分数
			FinalGradePercent:   grade.FinalGradePercent,   // 期末占比
			FinalGrade:          grade.FinalGrade,          // 期末分数
		})
	}

	//这里做了一个异步的增加用户的feedCount
	go func() {
		ct := context.Background()
		_, err := h.CounterClient.AddCounter(ct, &counterv1.AddCounterReq{StudentId: uc.StudentId})
		if err != nil {
			h.l.Error("增加用户feedCount失败:", logger.Error(err))
		}
	}()
	return web.Response{
		Msg:  fmt.Sprintf("获取%d~%d学年第%d学期成绩成功!", req.Xnm, req.Xnm+1, req.Xqm),
		Data: resp,
	}, nil
}

// GradeDetail 查询学分
// @Summary 查询学分
// @Description 查询学分
// @Tags 成绩
// @Accept json
// @Produce json
// @Success 200 {object} web.Response{data=GetGradeScoreResp} "成功返回学分"
// @Failure 500 {object} web.Response "系统异常，获取失败"
// @Router /grade/getGradeScore [get]
func (h *GradeHandler) GetGradeScore(ctx *gin.Context, uc ijwt.UserClaims) (web.Response, error) {
	// 调用 GradeClient 获取成绩数据
	score, err := h.GradeClient.GetGradeScore(ctx, &gradev1.GetGradeScoreReq{
		StudentId: uc.StudentId,
	})
	if err != nil {
		return web.Response{}, errs.GET_GRADE_SCORE_ERROR(err)
	}

	// 转换为目标结构体
	var resp GetGradeScoreResp
	for _, grade := range score.TypeOfGradeScore {
		typeOfGradeScore := TypeOfGradeScore{
			Kcxzmc:         grade.Kcxzmc,
			GradeScoreList: make([]*GradeScore, len(grade.GradeScoreList)),
		}

		for i := range grade.GradeScoreList {
			typeOfGradeScore.GradeScoreList[i] = &GradeScore{
				// 根据 GradeScore 的字段进行赋值
				Kcmc: grade.GradeScoreList[i].Kcmc,
				Xf:   grade.GradeScoreList[i].Xf,
			}
		}

		resp.TypeOfGradeScores = append(resp.TypeOfGradeScores, typeOfGradeScore)
	}

	return web.Response{
		Data: resp,
	}, nil
}
