package elecprice

type CheckRequest struct {
	Area     string `form:"area,omitempty"`     //区域,例如:南湖学生宿舍(看不懂参数请看这个网页https://jnb.ccnu.edu.cn/MobileWebPayStandard_Vue/#/addRoom,前两个参数直接完全一致,后面一个参数只保留了门牌号)
	Building string `form:"building,omitempty"` //建筑,例如:南湖05栋
	Room     string `form:"room,omitempty"`     //房间号,例如:414
}

type SetStandardRequest struct {
	Money    int64  `json:"money,omitempty"`    //金额
	Area     string `json:"area,omitempty"`     //区域
	Building string `json:"building,omitempty"` //建筑
	Room     string `json:"room,omitempty"`     //房间号
}

type CheckResponse struct {
	Price *Price `json:"price,omitempty"` // 电费
}

type Price struct {
	LightingRemainMoney       string `json:"lighting_remain_money,omitempty"`
	LightingYesterdayUseValue string `json:"lighting_yesterday_use_value,omitempty"`
	LightingYesterdayUseMoney string `json:"lighting_yesterday_use_money,omitempty"`
	AirRemainMoney            string `json:"air_remain_money,omitempty"`
	AirYesterdayUseValue      string `json:"air_yesterday_use_value,omitempty"`
	AirYesterdayUseMoney      string `json:"air_yesterday_use_money,omitempty"`
}
