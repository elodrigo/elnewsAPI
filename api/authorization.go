package api

import (
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type AccessTokens struct {
	KaKaoAccessTokenInfo
}

// 아이디 번호가 없을시 대신 0을 리턴시켜서 타입 맞추기
func Authorization(r *http.Request) (int64, error) {


	var err error
	var client *http.Client
	var accessTokenInfoResp *http.Response

	var accessTokens AccessTokens

	// 세션 가져오기
	session, err := Store.Get(r, "session")
	if err != nil {
		return 0, err
	}

	// 세션 정보 가져오기
	loginType := session.Values["login_type"]
	token := session.Values["token"]

	if token == nil {
		err = errors.New("token is empty")
		return 0, err
	}
	t, _ := reflect.ValueOf(token).Interface().(oauth2.Token)

	switch loginType {
	case "kakao":

		// 클라이언트 오브젝트에다가 토큰으로 설정 해놓기
		client = OAuthConfKaKao.Client(oauth2.NoContext, &t)

		// token을 카카오 api url로 보내 응답 받기
		accessTokenInfoResp, err = client.Get(KaKaoAccessTokenInfoURL)
		if err != nil {
			return 0, err
		}
		defer accessTokenInfoResp.Body.Close()

	default:

		// 클라이언트 오브젝트에다가 토큰으로 설정 해놓기
		client = OAuthConfKaKao.Client(oauth2.NoContext, &t)

		// token을 카카오 api url로 보내 응답 받기
		accessTokenInfoResp, err = client.Get(KaKaoAccessTokenInfoURL)
		if err != nil {
			return 0, err
		}
		defer accessTokenInfoResp.Body.Close()

	}

	// 가저온 정보 json 바디에서 읽어 들이기
	accessTokenInfo, err := ioutil.ReadAll(accessTokenInfoResp.Body)
	if err != nil {
		return 0, err
	}

	// AccessToken struct 에다가 json Unmarshal 하기
	err = json.Unmarshal(accessTokenInfo, &accessTokens.KaKaoAccessTokenInfo)
	if err != nil {
		return 0, err
	}

	// 0
	if accessTokens.ExpiresInMillis <= 0 {
		err = errors.New("expired: AccessToken is expired")
		return 0, err
	}

	// 유저의 번호로된 아이디와 함께 리턴
	return accessTokens.ID, err

}
