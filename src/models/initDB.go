package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	_"github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var config = map[string]string {
	"host" : "127.0.0.1",
	"port" : "3306",
	"dbName" : "go_test",
	"userName" : "root",
	"password" : "",
}

type Model struct {
	link *sql.DB   //储存连接对象
	tableName string	//储存表明
	fields string	//字段名
	allFields []string  //储存查询出来的字段
	where string	//条件
	order string	//排序
	limit string	//偏移量
}



func NewModel(table string) Model{
	fmt.Println("test newModel")
	var this Model
	this.fields = "*"
	this.tableName = table
	fmt.Println("test newModel2")
	this.getConnect()
	fmt.Println("test newModel3")

	this.getFields()
	fmt.Println(this.allFields)

	return this
	//sql := "select field from table where map order by orderField limit limitOffset "
}


func (this *Model)getConnect(){
	fmt.Println("test connect")
	//var config config
	fmt.Println(config["userName"]+":"+config["password"]+"@tcp("+config["host"]+":"+config["port"]+")/go_test?charset=utf8")

	db,err := sql.Open("mysql",config["userName"]+":"+config["password"]+"@tcp("+config["host"]+":"+config["port"]+")/go_test?charset=utf8")
	if err != nil{
		fmt.Println("connect fail")
		log.Fatal(err)
	}
	fmt.Println("test connect1")
	this.link = db
}


func (this *Model)getFields(){
	fmt.Println("test getFields")
	query := "desc "+this.tableName
	result,err := this.link.Query(query)
	fmt.Println("test query")
	Fatal(err)
	for result.Next(){
		var Field string
		var Type string
		var Null string
		var Key string
		var Default sql.NullString
		var Extra string
		err := result.Scan(&Field,&Type,&Null,&Key,&Default,&Extra)
		Fatal(err)
		this.allFields = append(this.allFields, Field)
	}

}


func (this *Model)Query(sql string) interface{}{
	rows,err := this.link.Query(sql)
	if err != nil{
		fmt.Println("sql执行错误："+sql)
		return returnRes(0,``,err)
	}

	cols,err := rows.Columns()
	if err != nil{

		return returnRes(0,``,err)
	}
	//表示一行所有列的值，用byte表示
	vals := make([][]byte,len(cols))
	//这里表示一行填充数据
	scans := make([]interface{},len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k,_ := range vals{
		scans[k] = &vals[k]
	}

	i := 0
	result := make(map[int]map[string]string)

	for rows.Next(){
		//填充数据
		rows.Scan(scans...) //将slic地址传入
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k,v := range vals{
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return returnRes(1,result,"success")
}

func (this *Model) Get() interface{}{
	sql := `select `+this.fields +` from ` + this.tableName + `` + this.where + `` + this.order + `` + this.limit
	//执行并发送SQL
	result := this.Query(sql)
	return result
}

func (this *Model) find(id int) interface{}{
	where := `where id=` + strconv.Itoa(id)
	sql := `select `+ this.fields + ` from ` + this.tableName + where +` limit 1 `
	//执行并发送sql
	result := this.Query(sql)
	return result
}

/**
	设置要查询的字段信息
 */
func (this *Model) Field(field string) *Model{
	this.fields = field
	return this
}


/**
	order排序条件
	string $order 以此为基准进行排序
	返回自己，保证连贯操作
 */
func (this *Model) Order(order string) *Model{
	this.order = `order by` + order
	return this
}


func (this *Model) Limit(limit int) *Model{
	this.limit = `limit` + strconv.Itoa(limit)
	return this
}

func (this *Model) Where(where string) *Model{
	this.where = ` where ` + where
	return this
}

/**
统计总条数
 */
func (this *Model) count() interface{}{
	sql := `select count(*) as total from ` + this.tableName + `limit 1 order by id`
	result := this.Query(sql)
	return result
}

/**
	执行并发送SQL语句（赠删改）
 */
func (this *Model) MyExec(sql string) interface{}{
	fmt.Println(sql)
	res,err := this.link.Exec(sql)
	fmt.Println(err)
	fmt.Println(res.LastInsertId())
	if err != nil{
		return returnRes(0,``,err)
	}
	return returnRes(1,res,"success")
}

/**
添加操作
 */
func (this *Model) add(data map[string]interface{}) interface{}{
	key := ""
	value := ""
	//过滤非法字段
	for k,v :=range data{
		if result := in_array(v,this.allFields);result != true {
			delete(data,k)
		}else{
			key += `,`+k
			value += `,`+`'` + v.(string) + `'`
		}
	}

	//将map中取出的键转为字符串拼接
	key = strings.TrimLeft(key,",")
	//将map中的值转化为字符串拼接
	value = strings.TrimLeft(value,",")
	//准备SQL语句
	sql := `insert into `+this.tableName + ` (`+ key + `) values (` +value + `)`
	//执行并发送SQL
	result := this.MyExec(sql)
	return result

}


func (this *Model)update(data map[string]interface{})interface{}{
	str := ""
	for k,v:=range data{
		if res:=in_array(v,this.allFields);res != true{
			delete(data,k)
		}else {
			str += k+` = '`+v.(string)+`',`
		}
	}
	//去掉最右侧的逗号
	str = strings.TrimRight(str,",")

	//判断是否有条件
	if this.where == ""{
		fmt.Println("没有条件")
	}

	sql := `update ` + this.tableName + ` set ` + str + ` ` +this.where

	result := this.MyExec(sql)
	return result
}



//是否存在数组中
func in_array(need interface{},hayStack []string) bool{
	for _,v:=range hayStack{
		if v == need{
			return true
		}
	}
	return false
}


func (this *Model) delete(id int) interface{}{
	//判断id是否存在
	where := `select * from ` + this.tableName + `where id = `+ strconv.Itoa(id)
	sql := `delete From` + this.tableName + `` + where
	res := this.MyExec(sql)
	return res
}





func Fatal(err error){
	if err != nil{
		log.Fatal(err)
	}
}

func returnRes(errCode int,res interface{},msg interface{}) string{
	result := make(map[string]interface{})
	result["errCode"] = errCode
	result["result"] = res
	result["msg"] = msg
	data,_ := json.Marshal(result)
	return string(data)
}

//func initDb(){
//
//}





//func main(){
//	M := models.NewModel()
//}