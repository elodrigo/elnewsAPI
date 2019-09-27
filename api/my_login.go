package api

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/kakao"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

var (

	myKey = []byte("key_here")
	Store = sessions.NewCookieStore(myKey)
	DomainName = "domain_here"

)

type LoginType struct {
	KaKao loginTypeSecond `json:"kakao"`
}

type loginTypeSecond struct {
	ClientID 		string `json:"client_id"`
	ClientSecret 	string `json:"client_secret"`
	RedirectURL 	string `json:"redirect_url"`
}

type LoginURL struct {
	LoginURL string `json:"login_url"`
}

type User struct {
	KaKaoUser
}



var OAuthConfKaKao *oauth2.Config

func init() {

	var loginType LoginType

	configFile, err := os.Open("config/my_login.json")
	if err != nil {
		log.Println(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)

	err = jsonParser.Decode(&loginType)
	if err != nil {
		log.Println(err)
	}

	OAuthConfKaKao = &oauth2.Config {
		ClientID: loginType.KaKao.ClientID,
		ClientSecret: loginType.KaKao.ClientSecret,
		RedirectURL: loginType.KaKao.RedirectURL,
		Endpoint: kakao.Endpoint,
	}

	gob.Register(oauth2.Token{})

}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {

	var client *http.Client
	//var status int

	var authUser User

	var userInfoResp *http.Response

	// 세션 가져오기
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Println(err)
	}

	// 세션 정보 가져오기
	loginType := session.Values["login_type"]
	token := session.Values["token"]

	if token == nil {
		log.Println("session is empty")
		return
	}

	t, _ := reflect.ValueOf(token).Interface().(oauth2.Token)

	switch loginType {
	case "kakao":
		// 클라이언트 오브젝트에다가 토큰으로 설정 해놓기
		client = OAuthConfKaKao.Client(oauth2.NoContext, &t)

		// 유저 인포 정보 카카오 url을 통해 가져오기
		userInfoResp, err = client.Get(KaKaoUserInfoAPIURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer userInfoResp.Body.Close()

	default:
		client = OAuthConfKaKao.Client(oauth2.NoContext, &t)

		// 유저 인포 정보 카카오 url을 통해 가져오기
		userInfoResp, err = client.Get(KaKaoUserInfoAPIURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer userInfoResp.Body.Close()

	}

	// 가져온 유저정보 json 바디에서 읽어 들이기
	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// authUser Struct에다가 json unmarshal 하기
	err = json.Unmarshal(userInfo, &authUser.KaKaoUser)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(authUser.KaKaoUser.ID)
	//session.Values["id"] = authUser.KaKaoUser.ID

	//fmt.Fprintf(w, string(userInfo))
	err = json.NewEncoder(w).Encode(authUser.KaKaoUser)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusAccepted)

}

func GetLoginURLJson(w http.ResponseWriter, r *http.Request) {

	loginType := r.FormValue("login_type")

	myLoginURL := GetLoginURL(loginType, "")

	loginURL := &LoginURL {

		LoginURL: myLoginURL,

	}

	err := json.NewEncoder(w).Encode(loginURL)
	if err != nil {
		log.Println(err)
	}

}

func GetLoginURL(loginType string, state string) string {

	var url string

	switch loginType {
	case "kakao":
		url = OAuthConfKaKao.AuthCodeURL(state)

	default:
		url = OAuthConfKaKao.AuthCodeURL(state)
	}

	return url
}

func Authenticate(w http.ResponseWriter, r *http.Request) {

	var token *oauth2.Token
	var err error
	log.Println("hi")


	// 세션 가져오기
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Println(err)
	}

	// 로그인 타입 가져오기
	vars := mux.Vars(r)
	loginType := vars["my_type"]

	// 엑세스 토큰 가져오기
	switch loginType {
	case "kakao":
		token, err = OAuthConfKaKao.Exchange(oauth2.NoContext, r.FormValue("code"))

	default:
		token, err = OAuthConfKaKao.Exchange(oauth2.NoContext, r.FormValue("code"))
	}
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	origin := r.Header.Get("Referer")
	//log.Println(origin)

	w.Header().Set("Location", origin)

	// 세션 저장하기전 옵션 설정하기
	session.Options = &sessions.Options{
		Path: "/",
		Domain: DomainName,
		MaxAge: 3600,

	}

	// 세션에다가 토큰 타입 정보 같은거 집어넣기
	session.Values["login_type"] = loginType
	session.Values["token"] = token
	log.Println(session)

	// 세션 저장하기
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, origin, http.StatusFound)
}

func LogoutSession(w http.ResponseWriter, r *http.Request) {

	var client *http.Client
	var logoutResp *http.Response
	var logout []byte

	w.Header().Set("Content-Type", "application/json")

	// 세션 가져오기
	session, err := Store.Get(r, "session")
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println(r)

	// 세션 정보 가져오기
	loginType := session.Values["login_type"]
	token := session.Values["token"]

	if token != nil {
		t, _ := reflect.ValueOf(token).Interface().(oauth2.Token)

		switch loginType {
		case "kakao":

			// 클라이언트 오브젝트에다가 토큰으로 설정 해놓기
			client = OAuthConfKaKao.Client(oauth2.NoContext, &t)

			// 카카오 url을 통해 로그아웃 요청 보내기
			logoutResp, err = client.Get(KaKaoLogoutAPIURL)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer logoutResp.Body.Close()

		}

		// json 바디에서 로그아웃 결과 읽어 들이기
		logout, err = ioutil.ReadAll(logoutResp.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNoContent)
			return
		}

	}

	//log.Println(string(logout))

	// 세션 저장하기전 옵션 설정하기. -1이면 세션 곧바로 삭제
	session.Options = &sessions.Options{
		Path: "/",
		Domain: DomainName,
		MaxAge: -1,
	}
	//origin := r.Header.Get("Referer")
	// 세션 저장하기
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	err = json.NewEncoder(w).Encode(logout)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusAccepted)
	return

}


func BaseResponseJSONStatus(status int, msg string) string {
	var rtnMsg string
	rtnMsg = fmt.Sprintf("{status:%d,msg:%s}", status, msg)
	return rtnMsg
}
