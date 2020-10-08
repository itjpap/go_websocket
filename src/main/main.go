package main

import (
	"auth"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"gitee.com/qq582826210/go-expect"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/websocket"
	"log"
	"models"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type User struct {

	conn *websocket.Conn
	msg  chan []byte
}



type Hub struct {
	//用户列表，保存所有用户
	userList map[*User]bool
	//注册chan，用户注册时添加到chan中
	register chan *User
	//注销chan，用户退出时添加到chan中，再从map中删除
	unregister chan *User
	//广播消息，将消息广播给所有连接
	broadcast chan []byte
}

//处理中心处理获取到的信息
func (h *Hub) run() {
	for {
		fmt.Println("阻塞通道读取消息")
		select {
		//从注册chan中取数据
		case user := <-h.register:
			//取到数据后将数据添加到用户列表中
			h.userList[user] = true
		case user := <-h.unregister:
			if _, ok := h.userList[user]; ok {
				delete(h.userList, user)
			}
		case data := <-h.broadcast:
			for u := range h.userList {
				select {
				case u.msg <- data:
				default:
					delete(h.userList, u)
					close(u.msg)
				}
			}
		}
	}
}

//定义一个websocket升级器
var up = &websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		if r.Method != "GET" {
			fmt.Println("请求方式错误")
			return false
		}
		if r.URL.Path != "/chat" {
			fmt.Println("请求路径错误")
			return false
		}
		//还可以根据其他需求定制校验规则
		return true
	},
}

var hub = &Hub{
	userList:   make(map[*User]bool),
	register:   make(chan *User),
	unregister: make(chan *User),
	broadcast:  make(chan []byte),
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	//通过升级后的升级器得到链接
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("获取连接失败:", err)
		return
	}
	fmt.Println("成功获取连接")

	user := &User{
		conn: conn,
		msg:  make(chan []byte),
	}

	fmt.Println("把用户放进注册列表中")
	fmt.Println(user)
	go hub.run()
	hub.register <- user

	//用户断开链接时添加用户到取消注册通道中
	defer func() {
		fmt.Println("用户断开链接")
		hub.unregister <- user
	}()
	fmt.Println("读取用户信息")
	go read(user)
	write(user)

}

func read(user *User) {
	fmt.Println("开始读取用户输入")
	//从连接中循环读取信息
	for {
		_, msg, err := user.conn.ReadMessage()
		fmt.Println(msg)
		if err != nil {
			fmt.Println("用户退出：", user.conn.RemoteAddr().String())
			hub.unregister <- user
			break
		}

		hub.broadcast <- msg
	}
}



func write(user *User) {
	for data := range user.msg {
		fmt.Println("开始写入用户流")
		err := user.conn.WriteMessage(1, data)
		if err != nil {
			fmt.Println("写入错误")
			break
		}
	}
}


func main() {
	fmt.Println("test log")
	//return  unsupported
	c,err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil{
		log.Fatal(err)
	}
	defer c.Close()

	cmd := exec.Command("./tt")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	go func(){
		c.ExpectEOF()
	}()

	err = cmd.Start()
	if err != nil{
		log.Fatal(err)
	}

	time.Sleep(time.Second)
	c.Send("Hello\r")
	time.Sleep(time.Second)
	c.Send("world\r")

	err = cmd.Wait()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("test log")
	//生成token
	//runtest()
	//runRegister()
	//model := models.NewModel("user")
	//res := model.MyExec(`INSERT INTO user (username,password,salt,access_token,sex,brithday)
	//VALUES ('baobao','123456','1','_OPAXS','男','1997-05-01')`)
	//fmt.Println(res)
	//auth.
	//fmt.Println("test run newModel")
	//model := models.NewModel("user")
	//fmt.Println(model.Fields)
	////fmt.Println(timeUnix)
	//go hub.run()
	//http.HandleFunc("/login",httpHandle)         //将chat请求交给wshandle处理  //wsHandle
	//http.HandleFunc("/chat",chatHandle)         //将chat请求交给wshandle处理  //wsHandle
	//http.ListenAndServe("127.0.0.1:8888", nil) //开始监听
}



//http处理中间件
func authMiddleware(next http.Handler) http.Handler{
	TestApiKey := "test_api_key"
	return http.HandlerFunc(func(rw http.ResponseWriter,req *http.Request){
		var apiKey string
		if apiKey = req.Header.Get("X-Api-Key"); apiKey != TestApiKey{
			log.Printf("bad auth api key :%s",apiKey)
			rw.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(rw,req)
	})
}

func MiddleWare(h http.Handler,middleware ...func(http.Handler) http.Handler) http.Handler{
	for _,mw := range middleware{
		h = mw(h)
	}
	return h

}


func chatHandle(w http.ResponseWriter,r *http.Request){
	query := r.URL.Query()
	jwtQuery,ok := query["jwt"]
	signQuery,ok2 := query["sign"]
	token := ""
	sign := ""
	if ok && ok2{
		token = jwtQuery[0]
		sign = signQuery[0]
	}else{
		fmt.Println("缺少验证参数")
		fmt.Println("缺少验证参数")
		return
	}

	fmt.Println(token)
	_,okSign := auth.ParseToken(token,sign)
	fmt.Println(okSign)

	if okSign {
		fmt.Println("用户:"+sign+"通过验证！")
	}
	wsHandle(w,r)

}

/**
	通过header方法验证jwt
 */
func headerJwt(w http.ResponseWriter,r *http.Request){
	header := r.Header
	fmt.Println(header)
	jwtHeader,ok := header["Jwt"]
	signHeader,ok2 := header["Sign"]
	token := ""
	sign := ""
	if ok && ok2 {
		token = jwtHeader[0]
		sign = signHeader[0]
	}else {
		log.Fatal("缺少header参数")
	}
	fmt.Println(token)
	fmt.Println(sign)
}



func Fatal(err error){
	fmt.Println(err)
}

func httpHandle(w http.ResponseWriter,r *http.Request){
	type sqlParse struct {
		errCode int
		msg string
		result interface{}
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	model := models.NewModel("user")
 	model.Where("username=" + "'" + username + "'")
	resQuery := model.Get()
	resMap := make(map[string] interface{})
	err := json.Unmarshal([]byte(resQuery.(string)),&resMap)
	if err != nil{
		fmt.Println(err)
	}
	//resQuery := model.Query("select * from where username= '" + username + "' limit 1")
	fmt.Println("json to mao",resMap)
	fmt.Println("the value of key1 is ",resMap["1"])
	fmt.Println("the value of key1 is ",resMap["errCode"])
	fmt.Println("the value of key1 is ",resMap["result"])
	ok := resMap["result"].(map[string]interface{})
	fmt.Println(ok)
	if len(ok) == 0{
		println(ok)
		fmt.Fprintln(w,"用户不存在")
	}

	dataArr := ok["0"].(map[string]interface{})
	detail := dataArr["password"]
	if detail != password {
		fmt.Fprintln(w,"账号密码不正确")
		return
	}

	token := auth.NewToken(dataArr["id"].(string))
	fmt.Fprintln(w,token)
	//fmt.Println(password)
	//type UserData struct{
	//	username string `json:"username"`
	//	password string `json:"password"`
	//}
	//var user UserData
	//err := json.NewDecoder(r.Body).Decode(&user)
	//fmt.Fprintln(w, "这是请求中的路径：", r.URL.Path)
	//fmt.Fprintln(w, "这是请求中的路径?后面的参数：", r.URL.RawQuery)
	//fmt.Fprintln(w, "这是请求中的User-Agent信息：", r.Header["User-Agent"])
	//fmt.Fprintln(w, "这是请求中的User-Agent信息：", r.Header.Get("User-Agent"))
	//
	//// 获取请求体内容的长度
	//// len := r.ContentLength
	//// body := make([]byte, len)
	//// r.Body.Read(body)
	//// fmt.Fprintln(w, "请求体中的内容是：", string(body))
	//
	//// 解析表单，在调用r.Form r.PostForm之前执行
	//r.ParseForm()
	//// fmt.Fprintln(w, "表单信息：", r.Form)
	//fmt.Fprintln(w, "表单信息：", r.PostForm)
	//
	//// fmt.Fprintln(w, "用户名：", r.FormValue("username"))
	//// fmt.Fprintln(w, "密码：", r.FormValue("password"))
	//fmt.Fprintln(w, "密码：", r.PostFormValue("password"))
	//panic("123")
	//if err != nil{
	//	w.WriteHeader(http.StatusForbidden)
	//}



}




//func runRegister(){
//	auth.RunServe()
//
//}


func runtest(){

	k := "96e3ba7d0560cb2d48a483d525aae4ca"
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.StandardClaims{
		ExpiresAt:time.Now().AddDate(0,0,1).Unix(),
		Id : "1",
	})
	fmt.Println(token)
	t,_ := token.SignedString([]byte(k))
	//验证token
	req,_ := http.NewRequest("GET","TEST",nil)
	req.Header.Add("token",t)
	fmt.Println(req.Header)

	token2, err := request.ParseFromRequest(req,request.HeaderExtractor{"token"},func(token *jwt.Token) (interface{},error){
		return []byte(k),nil
	})

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	sc := token2.Claims.(jwt.MapClaims)
	fmt.Println(sc)
}
