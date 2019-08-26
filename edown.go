package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"net/url"
	"github.com/zserge/lorca"
)

var (fname string)

func main() {
	ui, err := lorca.New("", "", 520, 320)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Bind("download", func(URL string) {
		ui.Eval(`document.querySelector('.done').innerText = ''`)
		UI(URL)
		ui.Eval(`document.querySelector('.done').innerText = '` + URL + ` done'`)
	})
	// Load HTML after Go functions are bound to JS
	ui.Load("data:text/html," + url.PathEscape(`
	<html>
	<head>
	<title>Edownloader</title>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	</head>
		<body>			
			<div class="field half">
			<label for="name" style="font-size:16px;">請輸入網址</label>
				<input id="URL" type="text" value=""  SIZE=40  height="35" style="font-size:16px;">
			</div>
			<input type="button" onclick="download(document.querySelector('#URL').value)" style="width:100px;height:30px;font-size:16px;" value="Download">
			<div class="done"></div>		
		</body>
	</html>
	`))
	<-ui.Done()
}
func UI(filename string) {

	html := gethtml(filename)
	fname = getname(html)
	pages := getpages(html) //找出頁數

	for i := 0; i < pages; i++ {
		html = gethtml(filename + "/?p=" + strconv.Itoa(i))
		picarr := getpics(html)
		for j := 1; j < len(picarr); j++ {
			picweb(picarr[j])
		}
	}
}

//call

//使用網址找出html
func gethtml(website string) string {
	res, err := http.Get(website)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	sitemap, err := ioutil.ReadAll(res.Body) //使用ioutil讀取body
	if err != nil {
		log.Fatal(err)
	}
	return string(sitemap)
}
func getURL() string {
	fmt.Println("Please enter URL: ")
	var URL string
	fmt.Scanln(&URL)
	return URL
}

//找出名稱
func getname(html string) string {
	name := catch(html, "<div id=\"gd2\">", "<div id=\"gright\">")
	name1 := catch(name, "<h1 id=\"gn\">", "</h1>")
	name2 := catch(name, "<h1 id=\"gj\">", "</h1></div>")
	var pagename string

	if name2 == "" {
		pagename = name1
	} else {
		pagename = name2
	}
	os.Mkdir(pagename, os.ModePerm) // 當前目錄建立資料夾
	return pagename
}

//抓出總共幾頁
func getpages(html string) int {
	html2 := catch(html, "<div class=\"gtb\">", "<div id=\"gdo\">")
	html3 := strings.Split(html2, "onclick=\"return false\"")
	var pages int

	if len(html3) > 2 {
		html4 := html3[len(html3)-2]
		pages, _ = strconv.Atoi(catch(html4, ">", "</a>"))
	} else {
		pages = 1
	}
	return pages
}

//將圖片的網址丟進陣列
func getpics(html string) []string {
	html2 := catch(html, "<div id=\"gdt\">", "</a></div></div><div class=\"c\"></div></div>")
	num := strings.Count(html, "<div class=\"gdtm\"")
	picarr := strings.Split(html2, "<div class=\"gdtm\"")

	for i := 1; i <= num; i++ {
		picarr[i] = catch(picarr[i], "<a href=\"", "\"><img alt=")
	}
	return picarr
}

//把單張圖片的網址抓出來
func picweb(website string) {
	html := gethtml(website)
	pic := catch(html, "<div id=\"i3\">", "<div id=\"i4\">")
	pic = catch(pic, "<img id=\"img\" src=\"", "\" style=\"")
	//fmt.Println(pic)
	getImg(pic)
}

//tool
//下載圖片
func getImg(url string) (n int64, err error) {
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		//name = fname + "/" + path[len(path)-1]
		name = fname + "/" + path[len(path)-1]
	}
	out, err := os.Create(name)
	defer out.Close()
	resp, err := http.Get(url)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	n, err = io.Copy(out, bytes.NewReader(pix))
	return
}

//抓出中間文字
func catch(input, str, end string) string {
	strn := strings.Index(input, str) + strings.Count(str, "") - 1
	endn := strings.Index(input, end)
	catchtml := string(input[strn:endn])
	return catchtml
}

/*
//將每一頁的網址找出來
func pageweb(html string, pages int) {
	block := catch(html, "<div class=\"gtb\">", "<div id=\"gdo\">")
	pagearr = strings.Split(block, "<a href=")
	//pages := getpages(html)
	for index := 1; index <= pages; index++ {
		pagearr[index] = catch(pagearr[index], "\"", "\" onclick=")
		//fmt.Println(pagearr[index])
	}
}
*/
