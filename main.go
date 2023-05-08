package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Blog struct {
	ID       int
	Title    string
	Content  string
	Image    string
	PostDate time.Time
	// StartDate string
	// EndDate   string
	Author string
	// NodeJs     bool
	// NextJs     bool
	// ReactJs    bool
	// TypeScript bool
}
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name    string
}

var userData = SessionData{}

func main() {

	connection.DatabaseConnect()

	// create new echo instance
	e := echo.New()

	// serve static files from directory
	e.Static("/public", "public")

	// initialitation to use session
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	//routing
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/blog", blog)
	e.GET("/detail/:id", detail)
	e.GET("/delete-blog/:id", deleteBlog)
	e.GET("/form-register", formRegister)
	e.GET("/form-login", formLogin)
	e.GET("/logout", logout)
	e.POST("/add-blog", addBlog)
	e.POST("/register", register)
	e.POST("/login", login)

	e.Logger.Fatal(e.Start("localhost:8001"))
}

func home(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	// map(tipe data) => key and value
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, content, image, post_date FROM tb_projects") // exc , start_date, end_date
	fmt.Println(data)

	var result []Blog

	for data.Next() {
		var each = Blog{}

		err := data.Scan(&each.ID, &each.Title, &each.Content, &each.Image, &each.PostDate) // exc , &each.StartDate, &each.EndDate
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		each.Author = "Asep"
		result = append(result, each)

	}

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	blogs := map[string]interface{}{
		"Blog":         result,
		"DataSession":  userData,
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
	}
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())
	return tmpl.Execute(c.Response(), blogs)
}

func contact(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/contact-me.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}

	blogs := map[string]interface{}{
		"DataSession": userData,
	}
	return tmpl.Execute(c.Response(), blogs)
}

func blog(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/blog.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func detail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	tmpl, err := template.ParseFiles("views/blog-detail.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	var BlogDetail = Blog{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, content, image, post_date FROM tb_projects WHERE id = $1", id).Scan(&BlogDetail.ID, &BlogDetail.Title, &BlogDetail.Content, &BlogDetail.Image, &BlogDetail.PostDate)
	// exc startDate, endDate, , &BlogDetail.StartDate, &BlogDetail.EndDate

	BlogDetail.Author = "ASEP"

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	data := map[string]interface{}{
		"Blog": BlogDetail,
	}

	return tmpl.Execute(c.Response(), data)
}

func addBlog(c echo.Context) error {
	title := c.FormValue("title") // Batch 46
	// startDate := c.FormValue("startDate")
	// endDate := c.FormValue("endDate")
	content := c.FormValue("content") // Finishing CRUD
	image := "tech1.jpg"

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(title, content, image, post_date) VALUES ($1, $2, $3, $4)", title, content, image, time.Now()) // exc start_date, end_date, startDate, endDate,

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteBlog(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) // id = 1 string => 1 int

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)

	// DELETE FROM tb_projects WHERE id=1
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}
func formRegister(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func formLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)
	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["alertStatus"], // true / false
		"FlashMessage": sess.Values["message"],     // "Register success"
	}

	delete(sess.Values, "message")
	delete(sess.Values, "alertStatus")

	tmpl, err := template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return redirectWithMessage(c, "Email Salah !", false, "/form-login")
	}

	fmt.Println(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return redirectWithMessage(c, "Password Salah !", false, "/form-login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 // 3 jam
	sess.Values["message"] = "Login Success"
	sess.Values["status"] = true // show alert
	sess.Values["name"] = user.Name
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true // access login
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func register(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// generate password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		redirectWithMessage(c, "Register failed, please try again :)", false, "/form-register")
	}

	return redirectWithMessage(c, "Register success", true, "/form-login")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Values["isLogin"] = false
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, path)
}
