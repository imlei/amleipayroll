// I am learning the golang, and this one is help me to improve my work profermanace.
// Idlookup v0.0.1 6/16/2020
package main

import (
	"database/sql"
	"fmt"
	"time"

	//	"io"
	//	"net/http"

	//	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

//数据库连接信息
const (
	USERNAME = "idlookup"
	PASSWORD = "idlookup"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "idlookup"
)

type idaccount struct {
	ID      int    `json:"id" from:"id"`
	Sysid   int    `json:"sysid" from:"sysid"`
	Nsid    string `json:"nsid" from:"nsid"`
	Mrpid   string `json:"mrpid" from:"mrpid"`
	Cusname string `json:"Cusname" from:"Cusname"`
	Custype string `json:"Custype" from:"Custype"`
}

//var db *sql.DB
//var view *template.Template

func main() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	DB, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	}

	DB.SetConnMaxLifetime(200 * time.Second) //最大连接周期，超时的连接就close
	DB.SetMaxOpenConns(100)                  //设置最大连接数
	//CreateTable(DB)
	InsertData(DB)
	QueryOne(DB)
	QueryMulti(DB)
	//UpdateData(DB)
	//DeleteData(DB)

	// 准备模板
	//err = LoadTemplate()
	//if err != nil {
	//	panic(err)
	//}

	// 注册处理函数
	//http.HandleFunc("/load", loadHandler)
	//http.HandleFunc("/", listHandler)
	//http.HandleFunc("/idlookup", idlookupHandler)

	// 启动服务器
	//err = http.ListenAndServe(":12345", nil)
	//if err != nil {
	//	panic(err)
	//}

}

//插入数据
func InsertData(DB *sql.DB) {
	result, err := DB.Exec("insert INTO idaccount(Sysid,Nsid,Mrpid,Cusname,Custype) values(?,?,?,?,?)", "9829", "VC-EC9937", "5254", "2025 SOLUTIONS LTD", "CAD")
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}
	lastInsertID, err := result.LastInsertId() //获取插入数据的自增ID
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return
	}
	fmt.Println("Insert data id:", lastInsertID)

	rowsaffected, err := result.RowsAffected() //通过RowsAffected获取受影响的行数
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}

//查询单行

func QueryOne(DB *sql.DB) {

	idlook := new(idaccount) //用new()函数初始化一个结构体对象
	row := DB.QueryRow("select ID, Sysid, Mrpid, Cusname, Custype from idaccount where id=?", 8)
	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	if err := row.Scan(&idlook.ID, &idlook.Sysid, &idlook.Mrpid, &idlook.Cusname, &idlook.Custype); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Println("Single row data:", *idlook)

}

//查询多行
func QueryMulti(DB *sql.DB) {
	idlook := new(idaccount)
	rows, err := DB.Query("select ID, Sysid, Mrpid, Cusname, Custype from idaccount where id < ?", 10)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		fmt.Printf("queryMulti Query failed,err:%v\n", err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&idlook.ID, &idlook.Sysid, &idlook.Mrpid, &idlook.Cusname, &idlook.Custype)
		if err != nil {
			fmt.Printf("queryMulti Scan failed,err:%v\n", err)
			return
		}
		fmt.Print(*idlook)
	}
	fmt.Printf("\n")
}

// 加载模板
/* func LoadTemplate() error {
	// 准备模板函数
	funcs := make(template.FuncMap)
	//funcs["showtime"] = ShowTime

	// 准备模板
	v := template.New("view")
	v.Funcs(funcs)

	_, err := v.ParseGlob("view/*.htm")
	if err != nil {
		return err
	}

	view = v
	return nil
}

// 动态加载模板 /load
func loadHandler(w http.ResponseWriter, req *http.Request) {
	err := LoadTemplate()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	io.WriteString(w, `模板加载完成`)
}

// 显示留言页面 /
func listHandler(w http.ResponseWriter, req *http.Request) {
	// 查询数据
	rows, err := db.Query("SELECT * FROM id_account")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	// 获取数据
	lys := []idlookup{}
	for rows.Next() {
		ly := idlookup{}
		err := rows.Scan(&ly.id, &ly.sysid, &ly.mrpid, &ly.nsid, &ly.Cusname, &ly.Custype)
		if nil != err {
			http.Error(w, err.Error(), 500)
			return
		}
		lys = append(lys, ly)
	}

	// 显示数据
	err = view.ExecuteTemplate(w, "index.htm", lys)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 留言页面
func idlookupHandler(w http.ResponseWriter, req *http.Request) {
	if "POST" == req.Method {
		// 获取参数
		Cusname := strings.TrimSpace(req.FormValue("name"))
		Custype := strings.TrimSpace(req.FormValue("type"))
		mrpid := strings.TrimSpace(req.FormValue("id"))

		// 检查参数
		if name == "" || content == "" {
			io.WriteString(w, "参数错误!\n")
			return
		}

		// sql语句
		sql, err := db.Prepare("insert into idlookup(Cusname, Custype, mrpid) values(?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer sql.Close()

		// sql参数,并执行
		_, err = sql.Exec(Cusname, Custype, mrpid)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// 跳转
		w.Header().Add("Location", "/")
		w.WriteHeader(302)

		// 提示信息
		io.WriteString(w, "提交成功!\n")

		return
	}

	// 显示表单
	err := view.ExecuteTemplate(w, "idlist.htm", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
*/

