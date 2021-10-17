package testdata

type AuthReq struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

type AuthRes struct {
	Success string `body:"text"`
	Error   error  `body:"text"  code:"403"`
}

type User struct {
	UserId   uint64 `json:"userId"`
	Username string `json:"username"`
}

type ListReq struct {
	PageSize int `query:"pageSize"`
	PageNo   int `query:"pageNo"`
}

type ListRespBase struct {
	Tot int `json:"tot"`
}

type ListUserResp struct {
	ListRespBase
	List []*User `json:"list"`
}

type FuzzyQueryUser struct {
	BlurUsername *string `query:"blurUsername"`
}
