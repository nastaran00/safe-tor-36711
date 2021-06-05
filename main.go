package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dbms/rec/handler"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

/////////////////////////////////////////////////////
func New() (http.Handler, func(h http.Handler) gin.HandlerFunc) {
	nextHandler := new(connectHandler)
	makeGinHandler := func(h http.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			state := &middlewareCtx{ctx: c}
			ctx := context.WithValue(c.Request.Context(), nextHandler, state)
			h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
			if !state.childCalled {
				c.Abort()
			}
		}
	}
	return nextHandler, makeGinHandler
}

// Wrap takes the common HTTP middleware function signature, calls it to generate
// a handler, and wraps it into a Gin middleware handler.
//
// This is just a convenience wrapper around New.
func Wrap(f func(h http.Handler) http.Handler) gin.HandlerFunc {
	next, adapter := New()
	return adapter(f(next))
}

type connectHandler struct{}

// pull Gin's context from the request context and call the next item
// in the chain.
func (h *connectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.Context().Value(h).(*middlewareCtx)
	defer func(r *http.Request) { state.ctx.Request = r }(state.ctx.Request)
	state.ctx.Request = r
	state.childCalled = true
	state.ctx.Next()
}

type middlewareCtx struct {
	ctx         *gin.Context
	childCalled bool
}

// our main function
func main() {
	fmt.Println("Starting...")
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Println(e)
	}

	if compose := os.Getenv("indocker"); compose == "dockercompose" {
		fmt.Println("Sleeping for 5 to wait for DB")
		time.Sleep(time.Second * 5)
	}

	if needsAuth := os.Getenv("NEEDS_AUTH"); needsAuth == "yes" {
		fmt.Println("Readying DB")
		handler.ReadyDB()
	}

	router.LoadHTMLGlob("views/*.html")
	router.Static("/css", "views/css")
	router.Static("/fonts", "views/fonts")
	router.Static("/img", "views/img")
	router.Static("/js", "views/js")
	router.GET("/", handler.RenderHome)
	router.GET("/signuplogin", handler.Signuplogin)
	router.GET("/signup", handler.Signup)
	router.Use(Wrap(handler.JwtAuthentication))

	router.GET("/home", handler.Homelogged)
	router.POST("/user/new", handler.Addpeople)
	router.POST("/user/login", handler.Authenticate)
	router.Run()
	handlers.LoggingHandler(os.Stdout, router)

	// port := os.Getenv("PORT")

	// router := mux.NewRouter()
	// router.Use(handler.JwtAuthentication)

	// router.HandleFunc("/api/user/new", handler.CreateAccount).Methods("POST")
	// router.HandleFunc("/api/user/login", handler.Authenticate).Methods("POST")

	// fmt.Sprintf(":%v", port)

	// handlers.LoggingHandler(os.Stdout, router)

	//log.Fatal(http.ListenAndServe(p, loggedRouter))

}
