package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"encoding/json"
    "io/ioutil"
	"log"
	"os"
	"math/rand"
	"strings"
	"time"
	"net/http"
	"html/template"
	_ "github.com/go-sql-driver/mysql"
)

type Dagashi_data struct {
    Data     [][]string `json:"data"`
}


type Result struct {
	id           int
	name         string
	taste        string
	texture      string
	information  string
	price        int
}

type Values struct {
	Id           int
	Name         string
	Taste        string
	Texture      string
	Information  string
	Price        int
}



// データベース読み込み
func ReadAll(db *sql.DB, sizzle_strings string) []Result {
	var results []Result //資料(https://golangstart.com/struct_slice/)

	rows, err := db.Query("select * from sizzle where taste in" + sizzle_strings + "or texture in" + sizzle_strings + "or information in"+ sizzle_strings)

	if err != nil {
		panic(err)
	}

	// データベースへのスキャン
	for rows.Next() {
		result := Result{}
		err = rows.Scan(&result.id, &result.name, &result.taste, &result.texture, &result.information, &result.price)
		if err != nil {
			panic(err)
		}
		results = append(results, result)
	}
	rows.Close()

	return results
}

// データベースログイン
func open(path string, count uint) *sql.DB {
	db, err := sql.Open("mysql", path)
	if err != nil {
		log.Fatal("open error:", err)
	}

	if err = db.Ping(); err != nil {
		time.Sleep(time.Second * 2)
		count--
		fmt.Printf("retry... count:%v\n", count)
		return open(path, count)
	}

	fmt.Println("db connected!!")
	return db
}

// データベース設定
func connectDB() *sql.DB {
	var path string = fmt.Sprintf("%s:%s@tcp(db:3306)/%s?charset=utf8&parseTime=true",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_DATABASE"))

	return open(path, 100)
}

// データベースに挿入
func InsertDB(db *sql.DB, name,taste,texture,information string, price int) {

	// insert処理
    ins, err := db.Prepare("INSERT INTO sizzle(name,taste,texture,information,price) VALUES(?,?,?,?,?)")
    if err != nil {
        log.Fatal(err)
    }
    ins.Exec(name, taste, texture, information, price)
}

// テーブルデータ削除
func TruncateDB(db *sql.DB) {
	db.Exec("TRUNCATE sizzle")
}

// リストシャッフル
func Shuf(rand_slice []Result) {

	//fmt.Println(rand_slice)
	rand.Seed(time.Now().UnixNano()) //乱数のシード設定

	rand.Shuffle(len(rand_slice), func(i, j int) {
		rand_slice[i], rand_slice[j] = rand_slice[j], rand_slice[i]
	})
	
	//fmt.Println(rand_slice)
	//return rand_slice
}

//フォーム画面を出力
func HandlerInput(w http.ResponseWriter, r *http.Request) {
	
	// HTMLファイルとの接続
	template := template.Must(template.ParseFiles("./index.html"))
	
	// 画面表示
	err := template.ExecuteTemplate(w, "index.html", nil)
	
	if err != nil {
		log.Fatalln(err)
	}
	
}


// 入力された名前を画面に表示
func HandlerResult(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("画面遷移完了")

	// フォーム値を格納
	setting_price, _ := strconv.Atoi(r.FormValue("budget"))                 // 設定金額
	setting_taste := strings.Split(r.FormValue("taste"), ",")               // 味覚系シズル
	setting_texture := strings.Split(r.FormValue("texture"), ",")           // 食感系シズル
	setting_information := strings.Split(r.FormValue("information"), ",")   // 情報系シズル

	// シズル設定
	setting_sizzle := append(append(setting_taste, setting_texture...), setting_information...)
	
	// シズル情報送信用
    sizzle_strings := "('"
	for i := 0; i < len(setting_sizzle); i++ {
		sizzle_strings = sizzle_strings + setting_sizzle[i]
		if i == (len(setting_sizzle)-1) {
            sizzle_strings = sizzle_strings +"')"
			break
		} else {
			sizzle_strings = sizzle_strings + "', '"
		}
	}

	// データベース接続
	db := connectDB()
	defer db.Close()

	// 駄菓子格納リスト
	var Read_results []Result

    // お菓子データ抽出
	Read_results = ReadAll(db, sizzle_strings)

	// 駄菓子格納リスト_2
	Read_results_2 := make([]Result, len(Read_results))
	copy(Read_results_2, Read_results)

	// 駄菓子格納リスト_3
	Read_results_3 := make([]Result, len(Read_results))
	copy(Read_results_3, Read_results)

	// シャッフル
	Shuf(Read_results)
	// シャッフル
	Shuf(Read_results_2)
	// シャッフル
	Shuf(Read_results_3)
    
	// 合計金額変数
	var price_sum int = 0 

    // 表示させる駄菓子リスト
	var select_results []Result

	all_Read_results := append(append(Read_results, Read_results_2...), Read_results_3...)
	
	for _, v := range all_Read_results {

		// 設定金額 > 現在の合計金額+商品金額 の場合
		if setting_price > (price_sum + v.price) {
			
			// リストに追加
			select_results = append(select_results, v)
			
			// 金額合算
			price_sum += v.price
			
		} else { // それ以外
			continue
		}
		
	}
	// 表示用データ
	var result_data = []Values{}
	
	// html反映用のデータ作成
    for i := 0; i < len(select_results); i++ {

		values := Values{
			Id: select_results[i].id,
			Name: select_results[i].name,
			Taste: select_results[i].taste,
			Texture: select_results[i].texture,
			Information: select_results[i].information,
			Price: select_results[i].price,
		}

		result_data = append(result_data, values)
	}

	// HTMLファイルとの接続
	template := template.Must(template.ParseFiles("result.html"))

	// 画面表示
	err := template.ExecuteTemplate(w, "result.html", result_data[:])

	if err != nil {
		log.Fatalln(err)
	}
}



func main() {

	// 変数定義
	var dagashi Dagashi_data

	// データベース接続
	db := connectDB()

    // テーブルクリア
	TruncateDB(db)	

	// jsonファイル読み込み
	raw, err := ioutil.ReadFile("./sizzle_budget_table.json")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

	// jsonから構造体に
    json.Unmarshal(raw, &dagashi)
	

	// データをテーブルに挿入
	for _, daga_list := range dagashi.Data[:] {
		temp_name := daga_list[0]
		temp_taste := daga_list[1]
		temp_texture := daga_list[2]
		temp_information := daga_list[3]
		temp_price, _ := strconv.Atoi(daga_list[4]) 

		// insert処理
		InsertDB(db, temp_name, temp_taste, temp_texture, temp_information, temp_price)
    }

	// データベース閉じる
	defer db.Close()

	// 確認
	fmt.Println("実行OK")

	// css,jsファイルの読み込み
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	
	// 初期画面実行
	http.HandleFunc("/", HandlerInput)
	
	// 画面遷移実行
	http.HandleFunc("/result", HandlerResult)

	// サーバー起動
	http.ListenAndServe(":8080", nil)

}