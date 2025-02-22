package feed

type FeedEvent struct {
	Id           int64             `json:"id"`
	Title        string            `json:"title"`
	Type         string            `json:"type"`
	Content      string            `json:"content"`
	CreatedAt    int64             `json:"created_at"`
	ExtendFields map[string]string `json:"extend_fields"`
}

type MuxiOfficialMSG struct {
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	ExtendFields map[string]string `json:"extend_fields"`
	PublicTime   int64             `json:"public_time"`
	Id           string            `json:"id"`
}

type GetFeedEventsResp struct {
	ReadEvents   []FeedEvent `json:"read_events"`
	UnreadEvents []FeedEvent `json:"unread_events"`
}

type ClearFeedEventReq struct {
	FeedId int64 `json:"feed_id,omitempty"`
}

type ReadFeedEventReq struct {
	FeedId int64 `json:"feed_id"`
}

type ChangeFeedAllowListReq struct {
	Grade          bool `json:"grade"`
	Muxi           bool `json:"muxi"`
	Holiday        bool `json:"holiday"`
	AirConditioner bool `json:"air_conditioner"`
	Light          bool `json:"light"`
}

type GetFeedAllowListResp struct {
	Grade          bool `json:"grade"`
	Muxi           bool `json:"muxi"`
	Holiday        bool `json:"holiday"`
	AirConditioner bool `json:"air_conditioner"`
	Light          bool `json:"light"`
}
type ChangeElectricityStandardReq struct {
	ElectricityStandard bool `json:"electricity_standard"`
}

type SaveFeedTokenReq struct {
	Token string `json:"token"`
}
type RemoveFeedTokenReq struct {
	Token string `json:"token"`
}

type PublicMuxiOfficialMSGReq struct {
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	ExtendFields map[string]string `json:"extend_fields,omitempty"`
	LaterTime    int64             `json:"later_time"`
}

type PublicMuxiOfficialMSGResp struct {
	Title        string            `json:"title"`
	Content      string            `json:"content"`
	PublicTime   string            `json:"public_time"`
	ExtendFields map[string]string `json:"extend_fields,omitempty"`
	Id           string            `json:"id"`
}

type StopMuxiOfficialMSGReq struct {
	Id string `json:"id"`
}

type GetToBePublicMuxiOfficialMSGResp struct {
	MSGList []MuxiOfficialMSG `json:"msg_list"`
}
