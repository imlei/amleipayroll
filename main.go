// I am learning the golang, and this one is help me to improve my work profermanace.
// Idlookup v0.1 6/16/2020
package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

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

type idlookup struct {
	id      int
	sysid   int
	mrpid   string
	nsid    string
	Custype string
	Cusname string
}

var db *sql.DB
var view *template.Template

func main() {

	// 连接数据库
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	var err error
	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 准备模板
	err = LoadTemplate()
	if err != nil {
		panic(err)
	}

	// 注册处理函数
	http.HandleFunc("/load", loadHandler)
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/idlookup", idlookupHandler)

	// 启动服务器
	err = http.ListenAndServe(":12345", nil)
	if err != nil {
		panic(err)
	}
}

// 加载模板
func LoadTemplate() error {
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
