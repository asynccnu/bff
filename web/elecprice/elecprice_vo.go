package elecprice

type SetStandardRequest struct {
	Money    int64  `json:"money,omitempty"`    //金额
	Area     string `json:"area,omitempty"`     //区域
	Building string `json:"building,omitempty"` //建筑
	Room     string `json:"room,omitempty"`     //房间号
	LightID  string `json:"light_id,omitempty"` // 灯ID
	AirID    string `json:"air_id,omitempty"`   // 空调ID
}

type Price struct {
	LightingRemainMoney       string `json:"lighting_remain_money,omitempty"`
	LightingYesterdayUseValue string `json:"lighting_yesterday_use_value,omitempty"`
	LightingYesterdayUseMoney string `json:"lighting_yesterday_use_money,omitempty"`
	AirRemainMoney            string `json:"air_remain_money,omitempty"`
	AirYesterdayUseValue      string `json:"air_yesterday_use_value,omitempty"`
	AirYesterdayUseMoney      string `json:"air_yesterday_use_money,omitempty"`
}

type GetAIDandNameRequest struct {
	AreaName string `json:"area_name,omitempty"`
}

type Architecture struct {
	ArchitectureName string `json:"architecture_name,omitempty"`
	ArchitectureID   string `json:"architecture_id,omitempty"`
}

type GetAIDandNameResponse struct {
	ArchitectureList []*Architecture `json:"architecture_list,omitempty"`
}

type GetRoomInfoRequest struct {
	ArchitectureID string `json:"architecture_id,omitempty"`
	Floor          string `json:"floor,omitempty"`
}

type Room struct {
	RoomID   string `json:"room_id,omitempty"`
	RoomName string `json:"room_name,omitempty"`
}

type GetRoomInfoResponse struct {
	RoomList []*Room `json:"room_list,omitempty"`
}

type GetPriceRequest struct {
	RoomAircID  string `json:"room_airc_id,omitempty"`
	RoomLightID string `json:"room_light_id,omitempty"`
}

type GetPriceResponse struct {
	Price *Price `json:"price,omitempty"`
}
