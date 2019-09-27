package api


type KaKaoUser struct {
	ID			int64	`json:"id"`
	//Properties	string	`json:"properties"`
	Properties	KaKaoUserProperties	`json:"properties"`
	HasEmail 	bool 	`json:"kakao_account.has_email"`
	Email		string 	`json:"kakao_account.email"`
	HasAge		bool 	`json:"kakao_account.has_age_range"`
	Age 		string 	`json:"kakao_aacount.age_range"`
	HasBirthday bool 	`json:"kakao_account.has_birthday"`
	Birthday	string	`json:"kakao_account.birthday"`
	HasGender	bool	`json:"kakao_account.has_gender"`
	Gender		string	`json:"kakao_account.gender"`
}

type KaKaoUserProperties struct {
	Nickname		string	`json:"nickname"`
	ProfileImage	string	`json:"profile_image"`
	ThumbnailImage	string	`json:"thumbnail_image"`
	//CustomField1	string	`json:"custom_field1"`
}

type KaKaoAccessTokenInfo struct {
	ID		int64	`json:"id"`
	ExpiresInMillis	int64	`json:"expiresInMillis"`
	AppID	int64	`json:"appId"`
}

const (

	KaKaoUserInfoAPIURL = "https://kapi.kakao.com/v2/user/me"
	KaKaoLogoutAPIURL = "https://kapi.kakao.com/v1/user/logout"
	KaKaoAccessTokenInfoURL = "https://kapi.kakao.com/v1/user/access_token_info"

)
