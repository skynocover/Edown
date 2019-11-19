package main

import (
	"github.com/zserge/lorca"
	"log"
	"net/url"

	"fmt"
	. "github.com/GoSpider"
	"os"
	"strconv"
	"strings"
)

var (
	viewtitle      = "Edownloader"
	inputlabel     = "請輸入網址"
	bottonname     = "Download"
	vieweight  int = 520
	viewheight int = 320

	ui, err        = lorca.New("", "", vieweight, viewheight)
	fname   string //資料夾名稱

)

func main() {
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	html := buildhtml()

	ui.Load("data:text/html," + url.PathEscape(html))
	ui.Bind("execute", func(input string) {
		ui.Eval(`document.querySelector('.done').innerText = ''`)
		download(input)
		if fname != "" {
			ui.Eval(`document.querySelector('.done').innerText = '` + fname + ` done'`)
		} else {
			ui.Eval(`document.querySelector('.done').innerText = 'LoadFail'`)
		}

	})
	<-ui.Done()
}

//out Call
func download(website string) {
	html := Gethtml(extoe(website)) //取得輸入網頁的html 修改EX的網址變成E

	//使用html得到本子名稱,頁數,跟圖片總張數
	fname = getname(html)

	if fname != "" {
		pages := getpages(html)
		images := getimages(html)
		//輪流至每一頁
		for i := 0; i < pages; i++ {
			//用每一頁的網址求出每一頁的html並取得當前頁面所有圖片的網址
			html = Gethtml(extoe(website) + "/?p=" + strconv.Itoa(i))
			picarr := getpics(html)
			//得到所有圖片網址後到每一個圖片內去載圖
			for j := 1; j < len(picarr); j++ {
				pic := picweb(picarr[j])
				go GetImg(fname, pic)
				//顯示當前載到第幾張圖
				ui.Eval(`document.querySelector('.done').innerText = '` + strconv.Itoa(i*40+j) + "/" + images + ` done'`)
			}
		}
	}
}

//修改網址,將EX改成E
func extoe(website string) (web string) {
	web = strings.Replace(website, "https://exhentai.org", "https://e-hentai.org", -1)
	return
}

//找出名稱並建立資料夾
func getname(html string) string {
	fmt.Println(html)
	name := catch(html, "<div id=\"gd2\">", "<div id=\"gright\">")
	if name != "" {
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
	} else {
		return ""
	}
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

//將找出圖片有幾張
func getimages(html string) string {
	html = catch(html, "<div id=\"asm\"><script", "<div id=\"gdo\">")
	html = catch(html, "<p class=\"gpc\">", "</p>")
	html = catch(html, "of", "images")
	html = strings.ReplaceAll(html, " ", "")
	//fmt.Println(html)
	return html
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
func picweb(website string) string {
	html := Gethtml(website)
	pic := catch(html, "<div id=\"i3\">", "<div id=\"i4\">")
	pic = catch(pic, "<img id=\"img\" src=\"", "\" style=\"")
	return pic
}

//tool

//抓出中間文字
func catch(input, str, end string) string {
	strn := strings.Index(input, str) + strings.Count(str, "") - 1
	endn := strings.Index(input, end)
	if strn < endn {
		catchtml := string(input[strn:endn])
		return catchtml
	} else {
		return ""
	}
}

//in call

func buildhtml() string {
	htmlstr := `<html>
	<head>
	<title>` + viewtitle + `</title>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	</head>
		<body>			
			<div class="field half"> `

	htmlend := `<div class="done"></div>		
	</body>
	</html> 
	`

	return htmlstr +
		`
	 	<label for="name" style="font-size:16px;">` + inputlabel + `</label>
		<input id="input" type="text" value=""  SIZE=40  height="35" style="font-size:16px;">
		</div>
		<input type="button" onclick="execute(document.querySelector('#input').value)" style="width:100px;height:30px;font-size:16px;" value="` + bottonname + `">
	` + htmlend
}
