package logic

import (
	"fmt"
	"github.com/reptile/config"
	"github.com/reptile/dependency_pack/goquery"
	"github.com/reptile/mysql"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
)

var (
	httpUrl = "fudao.qq.com/"
)

type Education struct {
	Subject string
	Number int
	Url string
}

type Courses struct {
	CoursesName			string
	Courses_id		string
	Title   		string
	Price   		string
	TeacherName     string
	Link			string
}

func Start(url string) {
	fmt.Println("正在爬取腾讯企鹅辅导相关数据请稍等.....")
	fmt.Println("爬取完成后会自动打开浏览器并浏览数据网站！")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)
	var List  = make(map[string]string)
	List = GetCourseList(html, List)
	s := GetCourseNumberList(List)
	sl := GetCourseSumList(s)

	ress := []Education{}
	res := Education{}
	var i  int
	var wg sync.WaitGroup
	for k, v := range sl{
		for _, v := range v {
			i++
			wg.Add(1)
			go GetCourseData(k, v, &wg)
		}
		res.Subject = k
		res.Number = i
		ress = append(ress, res)
		res = Education{}
		i = 0
	}
	wg.Wait()

	tmpl := template.Must(template.ParseFiles(config.HTMLAddr + "/html/course.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, struct{ Education []Education }{ress})
	})

	tmpl1 := template.Must(template.ParseFiles(config.HTMLAddr + "/html/courseList.html"))
	http.HandleFunc("/arith", func(w http.ResponseWriter, r *http.Request) {
		queryData := QueryDb("数学")
		tmpl1.Execute(w, struct{ Courses []Courses }{queryData})
	})

	http.HandleFunc("/chinese", func(w http.ResponseWriter, r *http.Request) {
		queryData := QueryDb("语文")
		tmpl1.Execute(w, struct{ Courses []Courses }{queryData})
	})

	http.HandleFunc("/english", func(w http.ResponseWriter, r *http.Request) {
		queryData := QueryDb("英语")
		tmpl1.Execute(w, struct{ Courses []Courses }{queryData})
	})

	switch fmt.Sprint(runtime.GOOS) {
	case "windows":
		cmd := exec.Command(`cmd`, `/c`, `start`, `http://localhost:8080`)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Start()
	case "darwin":
		exec.Command(`open`, `https://www.jianshu.com`).Start()
	case "linux":
		exec.Command(`xdg-open`, `https://www.jianshu.com`).Start()
	}

	log.Println("Listen Serve Addr : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func GetCourseList(html *goquery.Document, courseList map[string]string) map[string]string {
	// '//a[@class="tt"]/@href'
	html.Find(".subject-item").Find("a[class]").Each(func(i int, selection *goquery.Selection) {
		url, _ := selection.Attr("href")
		course1 := selection.Text()
		if course1 != "全部" && course1 != "讲座" {
			courseList[course1] = "https://" + url[2:]
		}
	})
	return courseList
}

func GetCourseNumberList(courseList map[string]string) map[string][]string {
	var result = make(map[string][]string)

	for k, v := range courseList{
		courseNum := GetCourseSum(v)
		result[k] = courseNum
	}

	return result
}

func GetCourseSum(url string) (res []string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	html, err := goquery.NewDocumentFromReader(resp.Body)

	html.Find("div[class=grade-area]").Find("a[class=grade-item]").Each(func(i int, selection *goquery.Selection) {
		url, _ := selection.Attr("href")
		res = append(res, "https://" + url[2:])
	})


	return res
}

func GetCourseSumList(val map[string][]string) (url map[string][]string) {
	//var wg *sync.WaitGroup

	url  = make(map[string][]string)

	for s, v := range val{
		func() {
			for _, vv := range v {
				//fmt.Println(s, ":", string(vv))
				resp, err := http.Get(string(vv))
				if err != nil {
					log.Println(err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					log.Printf("Error: status code %d", resp.StatusCode)
					//wg.Done()
					return
				}

				//bytes, err := ioutil.ReadAll(resp.Body)
				//
				//if err != nil {
				//	fmt.Println("ioutil.ReadAll err=",err)
				//	return
				//}
				//fmt.Println(string(bytes))

				html, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					log.Println(err)
					return
				}

				html.Find("li[class=courseGroup-card--wrapper]").Find("a[data-modid=sys_course_collection]").Each(func(i int, selection *goquery.Selection) {
					urlsss, _ := selection.Attr("href")
					//fmt.Println("https://" + urlsss)
					func() {

						resp, err := http.Get("https://" + urlsss[2:])
						if err != nil {
							log.Println(err)
							return
						}
						defer resp.Body.Close()
						if resp.StatusCode != http.StatusOK {
							log.Printf("Error: status code %d", resp.StatusCode)
							return
						}
						html, err := goquery.NewDocumentFromReader(resp.Body)
						html.Find("li[class=course-card]").Find("a[target=_blank]").Each(func(i int, selection *goquery.Selection) {
							//course, _ := selection.Attr("data-tdw")
							urls, _ := selection.Attr("href")
							url[s] = append(url[s], "https://" + httpUrl + urls[1:])
						})
					}()
				})

				html.Find("li[class=course-card]").Find("a[target=_blank]").Each(func(i int, selection *goquery.Selection) {
					//course, _ := selection.Attr("data-tdw")
					urls, _ := selection.Attr("href")
					url[s] = append(url[s], "https://" + httpUrl + urls[1:])
				})
			}
		}()
	}
	return url
}

func GetCourseData(coursename, url string, wg *sync.WaitGroup) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		wg.Done()
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: status code %d", resp.StatusCode)
		wg.Done()
		return
	}

	//bytes, err := ioutil.ReadAll(resp.Body)
	//
	//if err != nil {
	//	fmt.Println("ioutil.ReadAll err=",err)
	//	return
	//}
	//fmt.Println(string(bytes))

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		wg.Done()
		return
	}

	course := Courses{}

	course.CoursesName = coursename
	course.Courses_id = url[len(url) - 6:]
	course.Link = url

	html.Find(".tt-word").Each(func(i int, selection *goquery.Selection) {
		tt := selection.Text()

		course.Title = tt
	})

	html.Find("div[class=caption]>p").Each(func(i int, selection *goquery.Selection) {
		name := selection.Text()
		course.TeacherName = name
	})

	html.Find(".tt-price-integer").Each(func(i int, selection *goquery.Selection) {
		price := selection.Text()
		course.Price = price
	})

	db := mysql.DBCon()
		stmt, err := db.Prepare(
			"insert into courses (`courses_id`, `courses_name`, `price`, `teacherName`, `title`, `link`) values (?,?,?,?,?,?)")
		if err != nil {
			log.Println(err)
			wg.Done()
		}
		defer stmt.Close()
		rs, err := stmt.Exec(course.Courses_id, course.CoursesName, course.Price, course.TeacherName, course.Title, course.Link)
		if err != nil {
			log.Println(err)
			wg.Done()
		}
		if id, _ := rs.LastInsertId(); id > 0 {
			//log.Println("插入成功")
		}
	wg.Done()

}


func QueryDb(course string) (res []Courses) {
	courses := Courses{}
	db := mysql.DBCon()
	rows, err := db.Prepare("select title,teacherName,courses_id,price,link from courses where courses_name='" + course + "'")
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return
	}
	re, _ := rows.Query()
	for re.Next() {
		err = re.Scan(&courses.Title, &courses.TeacherName, &courses.Courses_id, &courses.Price,&courses.Link) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return
		}
		res = append(res, courses)
		courses = Courses{}
	}
	return res
}
