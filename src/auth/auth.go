//登录认证工具包
package auth

import (
	_ "bytes"
	_ "crypto"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "html"
	_ "io"
	"net/http"
	_ "net/url"
	"strconv"
	"strings"
	_ "strings"
	"time"
)

//type Token struct {
//	Raw String   //原始令牌。解析令牌时填充的
//	Method SigningMethid  //请求方法
//	Header map[string]interface{}  //在header头部里面的token
//	Claims Claims  //token的第二个参数
//	Signature string //签名
//	Valid bool  //验证令牌是否有效
//}

const (
	k = "96e3ba7d0560cb2d48a483d525aae4ca"
)

//生成token
func NewToken(key string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(key))
	return tokenString
}

func ParseToken(tokenString string, key string) (interface{}, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println(err)
		return "", false
	}
}


func main() {
	type UserInfo map[string]interface{}
	t := time.Now()
	key := "welcome to XXY's code world"
	userInfo := make(UserInfo)
	var expTime int64 = 1000 * 60 * 10
	var tokenState string
	userInfo["exp"] = "1515482650719371100"
	userInfo["unique"] = "0"

	tokenString := NewToken(key)
	claims, ok := ParseToken(tokenString, key)
	if ok {
		oldT, _ := strconv.ParseInt(claims.(jwt.MapClaims)["exp"].(string), 10, 64)
		ct := t.UTC().UnixNano()
		c := ct - oldT
		if c > expTime {
			ok = false
			tokenState = "Token 已过期"
		} else {
			tokenState = "Token 正常"
		}

	} else {
		tokenState = "token无效"
	}
	fmt.Println(tokenState)
	fmt.Println(claims)
}

//注册聊天账号
func RunServe() {
	//http.Handle("/register",registerHandle)
	http.HandleFunc("/register", registerHandle)
	http.ListenAndServe(":8000", nil)
}

func registerHandle(write http.ResponseWriter, request *http.Request) {
	username := request.PostFormValue("username")
	password := request.PostFormValue("password")

	//decoder := json.NewDecoder(request.Body)
	//解析参数
	//var params map[string]string
	//decoder.Decode(&params)
	fmt.Printf("POST json:username=%s,password=%s\n", username, password)
	fmt.Fprintf(write, `{"code":0}`)

}

//验证登录
func loginHandle(id int) {

	//models.Conn.DbSelect(id)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	const (
		SecretKey = "welcome to lujianjin's blog"
	)
	type Token struct {
		Token string `json:"token"`
	}

	type UserCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var user UserCredentials
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if strings.ToLower(user.Username) != "someone" {
		if user.Password != "p@ssword" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error extracting the key")
		fatal(err)
	}

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}
	response := Token{tokenString}
	JsonResponse(response, w)
}

func JsonResponse(responde interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(responde)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//
//
//
//func auth(){
//	res,err := http.PostForm(fmt.Sprintf("http://127.0.0.1:%v/authenticate",serverPort),url.Values{
//		"user" :{"test"},
//		"pass" : {"known"},
//	})
//
//	if err != nil{
//		fatal(err)
//	}
//
//	if res.StatusCode != 200 {
//		fmt.Println("Unexpected status code",res.StatusCode)
//	}
//
//	buf := new(bytes.Buffer)
//	io.Copy(buf,res.Body)
//	res.Body.Close()
//	tokenString := strings.TrimSpace(buf.String())
//
//	token,err := jwt.ParseWithClaims(tokenString,&CustomClaimsExample{},func(token *jwt.Token)(interface{},error){
//		return verifyKey,nil
//	})
//	fatal(err)
//
//	claims := token.Claims.(*CustomClaimsExample)
//	fmt.Println(claims.CustomerInfo.Name)
//
//}

func fatal(err error) {
	fmt.Println(err)
}

//func useTokenViaHttp(){
//	token,err = createToken("foo")
//	fatal(err)
//	req,err := http.NewRequest("GET",fmt.Sprintf("http://localhost:%v/restricted",serverPort),nil)
//	fatal(err)
//	req.Header.Set("Authorization",fmt.Sprintf("Bearer %v",token))
//	res,err := http.DefaultClient.Do(req)
//	fatal(err)
//	//读取响应体
//	buf := new(bytes.Buffer)
//	io.Copy(buf,res.Body)
//	res.Body.Close()
//	fmt.Println(buf.String())
//}
