package auth
//
//import (
//	"fmt"
//	"github.com/dgrijalva/jwt-go"
//	"github.com/dgrijalva/jwt-go/request"
//	"net/http"
//	"time"
//)
//const(
//	k = "96e3ba7d0560cb2d48a483d525aae4ca"
//)
//
//func new(){
//	timeUnix := time.Now().Unix()
//	//生成token
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
//		ExpiresAt:time.Now().AddDate(0,0,1).Unix(),
//		Id : "1",
//	})
//	fmt.Println(token)
//	t,_ := token.SignedString([]byte(k))
//	//验证token
//	req,_ := http.NewRequest("GET","TEST",nil)
//	req.Header.Add("token",t)
//	fmt.Println(req.Header)
//	token2, err := request.ParseFromRequest(req,request.HeaderExtractor{"token"},func(token *jwt.Token) (interface{},error){
//		return []byte(k),nil
//	})
//
//	if err != nil {
//		fmt.Println(err.Error())
//		panic(err)
//	}
//
//	sc := token2.Claims.(jwt.MapClaims)
//	fmt.Println(sc)
//
//}
//
//func check(){
//
//
//}
//
//func main(){
//
//}