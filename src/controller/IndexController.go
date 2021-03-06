package controller

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	//"github.com/foolin/echo-template"
	echotemplate "github.com/foolin/echo-template"
	echosession "github.com/go-session/echo-session"

	"github.com/dgrijalva/jwt-go"

	"github.com/twinj/uuid"

	"github.com/labstack/echo"

	//db "mzc/src/databases/store"
	"github.com/cloud-barista/cb-webtool/src/model"
	"github.com/cloud-barista/cb-webtool/src/service"
	"github.com/cloud-barista/cb-webtool/src/util"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// type ReqInfo struct {
// 	Email    string `email`
// 	Password string `password`
// }

// func Index(c echo.Context) error {

// 	// fmt.Println("=========== DashBoard start ==============")
// 	// if loginInfo := CallLoginInfo(c); loginInfo.UserID != "" {

// 	// 	return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")

// 	// }
// 	fmt.Println("=========== Index Controller nothing ==============")
// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// }

func Index(c echo.Context) error {
	fmt.Println("============== index ===============")

	//TODO : virtual Machine 이미지 목록을 가져온다.: 스토어에 저장하면 되나?
	// virtualMachineImageInfoList, respStatus := service.LookupVirtualMachineImageList(paramConnectionName)
	//TODO : server spec 목록을 가져온다.  : 스토어에 저장하면 되나?
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func About(c echo.Context) error {
	// return c.Render(http.StatusOK, "About.html", map[string]interface{}{})
	return echotemplate.Render(c, http.StatusOK,
		"About", // 파일명
		map[string]interface{}{
			"message": "",
			"status":  200,
		})
}

func MainForm(c echo.Context) error {

	loginInfo := service.CallLoginInfo(c)
	if loginInfo.UserID == "" {
		// Login 정보가 없으므로 login화면으로
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	nameSpaceInfoList, nsStatus := service.GetNameSpaceList()
	if nsStatus.StatusCode == 500 {
		// if nsErr != 200 {
		// return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// cloudOsList , _ := service.GetCloudOSList()
	if len(nameSpaceInfoList) == 1 { // namespace가 1개이면 mcis 체크
		defaultNameSpace := nameSpaceInfoList[0]
		// mcis가 있으면 dashboard로 ( dashboard에서 mcis가 없으면 mcis 생성화면으로 : TODO 현재 미완성으로 MCIS관리화면으로 이동)
		mcisList, _ := service.GetMcisList(defaultNameSpace.ID)
		if len(mcisList) > 0 {
			log.Println(" mcisList  ", len(mcisList))
			return c.Redirect(http.StatusTemporaryRedirect, "/operation/manages/mcis/mngform")
		} else {
			log.Println(" mcisList is null ", mcisList)
			return c.Redirect(http.StatusTemporaryRedirect, "/operation/manages/mcis/regform")
		}
	} else {
		return echotemplate.Render(c, http.StatusOK,
			"auth/Main", // 파일명
			map[string]interface{}{
				"LoginInfo": loginInfo,
				// "CloudOSList":               cloudOsList,
				"NameSpaceList": nameSpaceInfoList,
				"message":       nsStatus.Message,
				"status":        nsStatus.StatusCode,
			})
	}
}

func ApiTestMngForm(c echo.Context) error {
	return echotemplate.Render(c, http.StatusOK,
		"ApiTest", // 파일명
		map[string]interface{}{})
}

// API 호출 Test
func ApiCall(c echo.Context) error {
	fmt.Println("============== ApiCall ===============")

	// params := make(map[string]string)
	params := echo.Map{}
	if err := c.Bind(&params); err != nil {
		fmt.Println("err = ", err) // bind Error는 나지만 크게 상관없는 듯.
	}
	fmt.Println(params)
	// apiInfo := util.AuthenticationHandler()
	//name := c.FormValue("name")
	// paramApiTarget := c.Param("ApiTarget") // SPIDER인지, Tumblebug인지
	// paramApiURL := c.Param("ApiURL")       // 호출되는 경로 : 변수가 있더라도 변수까지 반영 된 최종 호출 될 url
	// paramApiMethod := c.Param("ApiMethod") // GET인지 POST인지
	// paramApiObj := c.Param("ApiObj")       // 호출에 사용되는 parameter들 (json형태)

	// fmt.Println("paramApiTarget=", paramApiTarget)
	// fmt.Println("paramApiURL=", paramApiURL)
	// fmt.Println("paramApiMethod=", paramApiMethod)
	// fmt.Println("paramApiObj=", paramApiObj)
	// paramUserID := strings.TrimSpace(reqInfo.UserID)
	apiTarget := ""

	if params["ApiTarget"] == "SPIDER" {
		apiTarget = util.SPIDER
	} else if params["ApiTarget"] == "TUMBLEBUG" {
		apiTarget = util.TUMBLEBUG
	} else if params["ApiTarget"] == "DRAGONFLY" {
		apiTarget = util.DRAGONFLY
	} else if params["ApiTarget"] == "LADYBUG" {
		apiTarget = util.LADYBUG
	}

	apiMethod := ""
	if params["ApiMethod"] == "GET" {
		apiMethod = http.MethodGet
	} else if params["ApiMethod"] == "POST" {
		apiMethod = http.MethodPost
	} else if params["ApiMethod"] == "PUT" {
		apiMethod = http.MethodPut
	} else if params["ApiMethod"] == "DELETE" {
		apiMethod = http.MethodDelete
	}

	//url := util.TUMBLEBUG + "/ns"
	url := apiTarget + fmt.Sprintf("%v", params["ApiURL"])

	fmt.Println("url=", url)

	// if params["ApiObj"] != "" {// ApiObj유무에 따라 CommonHttp, CommonHttpWithoutParam으로 나눌까 하다가 하나로 호출.
	pbytes := []byte(fmt.Sprintf("%v", params["ApiObj"])) // 없으면 없는대로 CommonHttp호출.
	// pbytes, _ := json.Marshal(paramApiObj)
	fmt.Println("CommonHttp=")
	resp, err := util.CommonHttp(url, pbytes, apiMethod)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"err": err,
		})
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	resultStr := buf.String()

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(resultStr), &jsonMap)

	return c.JSON(http.StatusOK, map[string]interface{}{
		// "resp": resultStr,
		"resp": jsonMap,
		"err":  err,
	})
	// } else {
	// 	fmt.Println("CommonHttpWithoutParam=")
	// 	resp, err := util.CommonHttpWithoutParam(url, apiMethod)

	// 	if err != nil {
	// 		fmt.Println("result error=", err)
	// 		// fmt.Println("result error=", err)
	// 		return c.JSON(http.StatusOK, map[string]interface{}{
	// 			"err": err,
	// 		})
	// 	}
	// 	buf := new(bytes.Buffer)
	// 	buf.ReadFrom(resp.Body)
	// 	resultStr := buf.String()
	// 	log.Println(resultStr)
	// 	// data, err := ioutil.ReadAll(resp.Body)
	// 	// if err != nil {
	// 	// 	return c.JSON(http.StatusOK, map[string]interface{}{
	// 	// 		"err": err,
	// 	// 	})
	// 	// }
	// 	jsonMap := make(map[string]interface{})
	// 	err = json.Unmarshal([]byte(resultStr), &jsonMap)

	// 	return c.JSON(http.StatusOK, map[string]interface{}{
	// 		"resp": jsonMap,
	// 		"err":  err,
	// 	})
	// }
}

// Server API 호출 Test : 너무 복잡함....
func ServerCall(c echo.Context) error {
	fmt.Println("============== ServerCall ===============")

	params := make(map[string]string)
	if err := c.Bind(&params); err != nil {
		fmt.Println("err = ", err) // bind Error는 나지만 크게 상관없는 듯.
	}
	fmt.Println(params)
	// apiInfo := util.AuthenticationHandler()
	//name := c.FormValue("name")
	// paramApiTarget := c.Param("ApiTarget") // SPIDER인지, Tumblebug인지
	// paramApiURL := c.Param("ApiURL")       // 호출되는 경로 : 변수가 있더라도 변수까지 반영 된 최종 호출 될 url
	// paramApiMethod := c.Param("ApiMethod") // GET인지 POST인지
	// paramApiObj := c.Param("ApiObj")       // 호출에 사용되는 parameter들 (json형태)

	// fmt.Println("paramApiTarget=", paramApiTarget)
	// fmt.Println("paramApiURL=", paramApiURL)
	// fmt.Println("paramApiMethod=", paramApiMethod)
	// fmt.Println("paramApiObj=", paramApiObj)
	// paramUserID := strings.TrimSpace(reqInfo.UserID)
	apiTarget := ""

	if params["ApiTarget"] == "SPIDER" {
		apiTarget = util.SPIDER
	} else if params["ApiTarget"] == "TUMBLEBUG" {
		apiTarget = util.TUMBLEBUG
	} else if params["ApiTarget"] == "DRAGONFLY" {
		apiTarget = util.DRAGONFLY
	}

	apiMethod := ""
	if params["ApiMethod"] == "GET" {
		apiMethod = http.MethodGet
	} else if params["ApiMethod"] == "POST" {
		apiMethod = http.MethodPost
	} else if params["ApiMethod"] == "PUT" {
		apiMethod = http.MethodPut
	} else if params["ApiMethod"] == "DELETE" {
		apiMethod = http.MethodDelete
	}

	//url := util.TUMBLEBUG + "/ns"
	url := apiTarget + params["ApiURL"]

	fmt.Println("url=", url)

	// if params["ApiObj"] != "" {// ApiObj유무에 따라 CommonHttp, CommonHttpWithoutParam으로 나눌까 하다가 하나로 호출.
	pbytes := []byte(params["ApiObj"]) // 없으면 없는대로 CommonHttp호출.
	// pbytes, _ := json.Marshal(paramApiObj)
	fmt.Println("CommonHttp=")
	resp, err := util.CommonHttp(url, pbytes, apiMethod)

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"err": err,
		})
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	resultStr := buf.String()

	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(resultStr), &jsonMap)

	return c.JSON(http.StatusOK, map[string]interface{}{
		// "resp": resultStr,
		"resp": jsonMap,
		"err":  err,
	})
}

func LoginForm(c echo.Context) error {
	fmt.Println("============== Login Form ===============")

	user := os.Getenv("LoginUser")
	email := os.Getenv("LoginEmail")
	pass := os.Getenv("LoginPassword")

	store := echosession.FromContext(c)
	obj := map[string]string{
		"userid":           user,
		"username":         user,
		"email":            email,
		"password":         pass,
		"defaultnamespage": "",
		"accesstoken":      "",
		"refreshtoken":     "",
	}
	store.Set(user, obj)
	store.Save() // 사용자정보를 따로 저장하지 않으므로 설정파일에 유저를 set.
	fmt.Println("user : ", user)

	return echotemplate.Render(c, http.StatusOK, "auth/Login", nil)
	//return c.Render(http.StatusOK, "Login.html", map[string]interface{}{})
}

func LoginProc(c echo.Context) error {
	fmt.Println("============== Login proc ===============")
	store := echosession.FromContext(c) // store내 param은 모두 소문자.

	// reqInfo := new(model.ReqInfo)
	// if err := c.Bind(reqInfo); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 		"message": "fail",
	// 		"status":  "fail",
	// 	})
	// }

	// paramUserID := strings.TrimSpace(reqInfo.UserID)
	// // paramEmail := strings.TrimSpace(reqInfo.Email)
	// paramPass := strings.TrimSpace(reqInfo.Password)

	paramUserID := c.FormValue("userID")
	paramPass := c.FormValue("password")
	fmt.Println("paramUser & getPass : ", paramUserID, paramPass)

	// echoSession에서 가져오기
	result, ok := store.Get(paramUserID)

	if !ok {
		log.Println(" login proc err  ", ok)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{ //401
			"message": " 정보가 없으니 다시 등록바랍니다.",
			"status":  "fail",
		})
	}
	storedUser := result.(map[string]string)
	fmt.Println("Stored USER:", storedUser)
	if paramUserID != storedUser["userid"] && paramUserID != storedUser["password"] {
		log.Println(" invalid id or pass  ")
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{ //401
			"message": "invalid user or password",
			"status":  "fail",
		})
	}

	newToken, createTokenErr := createToken(paramUserID)
	if createTokenErr != nil {
		log.Println(" login proc err  ", createTokenErr)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{ //401
			"message": " 로그인 처리 요류",
			"status":  "fail",
		})
	}
	log.Println("newToken  ", newToken)
	// "accesstoken" : "",
	// 	"refreshtoken" : "",
	// td.RefreshToken
	storedUser["accesstoken"] = newToken.AccessToken
	storedUser["refreshtoken"] = newToken.RefreshToken
	// store.Set(paramUser, storedUser)
	// store.Save()

	//////// 현재구조에서는 nsList 부분을 포함해야 함. TODO : 이부분 호출되는 화면에서 필요할 듯 한데.. 공통으로 뺄까?
	nsList, nsStatus := service.GetNameSpaceList()
	log.Println(nsStatus)
	if nsStatus.StatusCode == 500 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": nsStatus.Message,
			"status":  nsStatus.StatusCode,
		})
	}
	if len(nsList) == 0 {
		nameSpaceInfo, createNameSpaceStatus := service.CreateDefaultNamespace()
		if createNameSpaceStatus.StatusCode == 500 {
			log.Println(" default namespace create failed  ", createNameSpaceStatus)
			// login은 되었으므 main 화면까지는 보내야되지 않을까?
		} else {
			nsList = append(nsList)
			storedUser["defaultnamespacename"] = nameSpaceInfo.Name
			storedUser["defaultnamespaceid"] = nameSpaceInfo.Name
			// storedUser["defaultnamespaceid"] = nameSpaceInfo.ID
		}
	} else if len(nsList) == 1 {
		for i, nameSpaceInfo := range nsList {
			log.Println(i, nameSpaceInfo)
			storedUser["defaultnameSpacename"] = nameSpaceInfo.Name
			storedUser["defaultnamespaceid"] = nameSpaceInfo.Name
			// storedUser["defaultnamespaceid"] = nameSpaceInfo.ID
			// defaultNameSpace = nameSpaceInfo.Name // ID로 handling 하려면 ID로
		}
	} else {
		storedUser["defaultnameSpacename"] = ""
		storedUser["defaultnamespaceid"] = ""
	}

	store.Set("namespacelist", nsList)
	///////

	/////// connectionconfig 목록 조회 ////////
	cloudConnectionConfigInfoList, _ := service.GetCloudConnectionConfigList()
	store.Set("connectionconfig", cloudConnectionConfigInfoList)
	/////// connectionconfig 목록 조회 끝 ////////

	// // result := map[string]string{}
	// result := get.(map[string]string)
	// fmt.Println("result mapping : ", result)
	// for k, v := range get.(map[string]string) {
	// 	fmt.Println(k, v)
	// 	result[k] = v
	// }

	// Username:             storedUser["username"],
	// 	AccessToken:          storedUser["accesstoken"],
	// 	DefaultNameSpaceID:   storedUser["defaultnamespaceid"],
	// 	DefaultNameSpaceName: storedUser["defaultnameSpacename"],

	store.Set(paramUserID, storedUser)
	store.Save()

	loginInfo := model.LoginInfo{
		UserID:      paramUserID,
		Username:    paramUserID,
		AccessToken: storedUser["accesstoken"],
		//Username:  result["username"],
		DefaultNameSpaceID:   storedUser["defaultnamespaceid"],
		DefaultNameSpaceName: storedUser["defaultnamespacename"],
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":                       "success",
		"status":                        "200",
		"LoginInfo":                     loginInfo,
		"NameSpaceList":                 nsList,
		"CloudConnectionConfigInfoList": cloudConnectionConfigInfoList,
	})

	// return c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 	"message": "invalid user",
	// 	"status":  "403",
	// })

}

// ----------- 로그인이 성공하면 Namespace가 없으면 생성 ----------/
// ----------- Name Space가 1개 있으면 Dashboard로 이동 ----------/
// ----------- Name Space가 1개 이상 있으면 Dashboard로 이동 및 Namespace선택 Modal 띄우기 ----------/
// ----------- MCIS가 등록되어있지 않으면 등록화면으로 ----------/
// form 전송으로 쓰려 했으나 사용안함.
// func LoginProcess(c echo.Context) error {
// 	store := echosession.FromContext(c)

// 	paramUser := c.FormValue("username")
// 	paramPass := c.FormValue("password")

// 	// reqInfo := new(model.ReqInfo)
// 	// if err := c.Bind(reqInfo); err != nil {
// 	// 	return c.Redirect(http.StatusTemporaryRedirect, "/login")
// 	// }

// 	// paramUser := strings.TrimSpace(reqInfo.UserName)
// 	// // paramEmail := strings.TrimSpace(reqInfo.Email)
// 	// paramPass := strings.TrimSpace(reqInfo.Password)
// 	fmt.Println("paramUser & paramPass : ", paramUser, paramPass)

// 	// echoSession에서 가져오기
// 	result, ok := store.Get(paramUser)

// 	if !ok {
// 		log.Println(" login proc err  ", ok)
// 		return c.Redirect(http.StatusTemporaryRedirect, "/login")
// 	}
// 	storedUser := result.(map[string]string)
// 	// fmt.Println("Stored USER:", storedUser)
// 	if paramUser != storedUser["username"] || paramUser != storedUser["password"] {
// 		log.Println(" invalid id or pass  ")
// 		return c.Redirect(http.StatusTemporaryRedirect, "/login")
// 	}

// 	newToken, createTokenErr := createToken(paramUser)
// 	if createTokenErr != nil {
// 		log.Println(" createTokenErr  ", createTokenErr)
// 		return c.Redirect(http.StatusTemporaryRedirect, "/login")
// 	}

// 	storedUser["accesstoken"] = newToken.AccessToken
// 	storedUser["refreshtoken"] = newToken.RefreshToken
// 	store.Set(paramUser, storedUser)
// 	store.Save()
// 	// return c.Render(http.StatusBadRequest, "/setting/connections/CloudConnection", map[string]interface{}{
// 	// return c.Redirect(http.StatusTemporaryRedirect, "/setting/connections/CloudConnection")

// 	loginInfo := model.LoginInfo{
// 		Username:    storedUser["username"],
// 		AccessToken: storedUser["accesstoken"],
// 	}
// 	// log.Println("LoginProcess loginInfo  ", loginInfo)

// 	c.Response().Header().Set("username", loginInfo.Username)
// 	c.Response().Header().Set("logininfo", "username="+loginInfo.Username+", accesstoken="+loginInfo.AccessToken)
// 	//c.Response().WriteHeader(http.StatusOK)
// 	return c.Redirect(http.StatusMovedPermanently, "/setting/connections/CloudConnection")
// 	// return c.Redirect(http.StatusPermanentRedirect, "../setting/connections/CloudConnection") // POST 그대로 보낼 떄
// 	// return echotemplate.Render(c, http.StatusOK,
// 	// 	"setting/connections/CloudConnection",
// 	// 	map[string]interface{}{
// 	// 		"LoginInfo": loginInfo,
// 	// })
// 	// return echotemplate.Render(c, http.StatusOK,
// 	// 	"setting/connections/CloudConnection", nil,
// 	// )
// 	// return echotemplate.Render(c, http.StatusOK,
// 	// 	"Map.html",
// 	// 	map[string]interface{}{
// 	// 			"LoginInfo": loginInfo,
// 	// 		},
// 	// )
// 	// return echotemplate.Render(c, http.StatusOK,
// 	// 	"setting/connections/CloudConnection.html",
// 	// 	map[string]interface{}{
// 	// 			"LoginInfo": loginInfo,
// 	// 		},
// 	// )
// 	// return c.Render(http.StatusOK, "/setting/connections/CloudConnection.html", loginInfo)

// 	// 	storedUser["defaultnamespage"] = nameSpaceInfo.ID

// 	// 	// 저장 성공하면 namespace 목록 조회
// 	// 	nsList2, nsErr2 := service.GetNameSpaceList()
// 	// 	if nsErr2 != nil {
// 	// 		log.Println(" nsErr2  ", nsErr2)
// 	// 		return c.Redirect(http.StatusTemporaryRedirect, "/setting/connections/CloudConnection")
// 	// 	}
// 	// 	log.Println("nsList2  ", nsList2)
// 	// 	nsList = nsList2
// 	// }
// 	// log.Println("nsList  ", nsList)
// 	// store.Set("namespacelist", nsList)// 이게 유효한가?? 쓸모없을 듯
// 	// store.Save()

// 	// // mcis가 있으면 dashboard로

// 	// // mcis가 없으면 mcis 등록화면으로

// 	// // return c.Render(http.StatusBadRequest, "/setting/connections/CloudConnection", map[string]interface{}{
// 	// return c.Redirect(http.StatusTemporaryRedirect, "/setting/connections/CloudConnection")
// }

func createToken(userID string) (*TokenDetails, error) {

	// var err error

	// atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	// atClaims["username"] = username
	// atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	// at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// token, err := at.SignedString([]byte(os.Getenv("LoginAccessSecret")))

	// if err != nil {
	//    return "", err
	// }
	// return token, nil

	// 액세스 토큰(access token)이 만료된 경우 리프레시 토큰(refresh token)을 사용하여
	// 새 액세스 토큰을 생성하여 액세스 토큰이 만료가 되더라도 사용자가 다시 로그인을 하지 않게
	// refresh token은 30분에 token만료.
	// action이 있을 때마다 at, rt(refresh token)은 갱신
	// 5분 넘어 action이 발생했을 때(at가 expired) rt가 유효하면 로그인 된 것으로
	// 30분동안 action이 없으면 refresh token이 expire되므로 이후에는 로그인 필요
	// 페이지 호출할 때마다 유효성 검증 후 expired 시간 재할당.
	var err error

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 5).Unix()
	td.AccessUuid = uuid.NewV4().String()
	// td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	//atClaims["user_id"] = userid
	atClaims["userid"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("LoginAccessSecret")))

	if err != nil {
		log.Println("create accessToken  ", err)
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	// rtClaims["user_id"] = userid
	rtClaims["userid"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("LoginRefreshSecret")))

	if err != nil {
		log.Println("create RefreshToken  ", err)
		return nil, err
	}

	return td, nil
}

// Login 없이 접근가능
// Login이 필요없는 화면에서 호출하는게 의미 있나? 없이 써도 되는 듯.
func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

// Token이 있어야 접근가능
// login 이 필요한 page에서 호출하여 값이 true일 때만 접근가능
func restricted(c echo.Context) error {
	user := c.Get("UserID").(*jwt.Token)
	// user := c.Get("email").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["userid"].(string)
	return c.String(http.StatusOK, "Welcome "+userID+"!")
}

func RegUser(c echo.Context) error {
	//comURL := GetCommonURL()

	user := os.Getenv("LoginEmail")
	pass := os.Getenv("LoginPassword")

	store := echosession.FromContext(c)
	obj := map[string]string{
		"userid":   user,
		"password": pass,
	}
	store.Set(user, obj)
	err := store.Save()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"message": "Fail",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "SUCCESS",
		"user":    user,
	})

}

func LogoutForm(c echo.Context) error {
	fmt.Println("============== Logout Form ===============")
	//comURL := GetCommonURL()
	return c.Render(http.StatusOK, "logout.html", nil)
}

// 세션을 초기화 하고 login 화면으로 보낸다.
func LogoutProc(c echo.Context) error {
	fmt.Println("============== Logout proc ===============")
	store := echosession.FromContext(c)

	reqInfo := new(model.ReqInfo)

	getUser := strings.TrimSpace(reqInfo.UserID)

	store.Set(getUser, nil)
	store.Save()
	log.Println(" auth expired ")

	// return c.Render(http.StatusOK, "login.html", nil)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")

}
