package infrastructure

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hidenari-yuda/umerun-resume/domain/config"
	"github.com/hidenari-yuda/umerun-resume/infrastructure/database"
	"github.com/hidenari-yuda/umerun-resume/infrastructure/driver"
	"github.com/hidenari-yuda/umerun-resume/infrastructure/router/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct {
	cfg    config.Config
	Engine *echo.Echo
}

func NewRouter(cfg config.Config) *Router {
	return &Router{
		cfg:    cfg,
		Engine: echo.New(),
	}
}

func (r *Router) SetUp() *Router {
	var (
		db       = database.NewDB(r.cfg.DB, true)
		firebase = driver.NewFirebaseImpl(r.cfg.Firebase)
	// basicAuth = utility.NewBasicAuth(r.cfg)
	)

	// r.Engine.HidePort = true
	// r.Engine.HideBanner = true
	// r.Engine.Use(middleware.Recover())
	// // TODO: Webクライアントのドメインが決まったら設定する 👆の`r.Engine.Use(middleware.CORS())`は消す
	// // r.Engine.Use(middleware.CORSWithConfig((middleware.CORSConfig{
	// // AllowOrigins: r.cfg.App.CorsDomains,
	// // 	AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderAccessControlAllowOrigin},
	// // 	AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	// // })))
	// r.Engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	// }))
	// r.Engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Skipper: func(c echo.Context) bool {
	// 		if strings.Contains(c.Request().URL.Path, "healthz") {
	// 			return true
	// 		} else {
	// 			return false
	// 		}
	// 	},
	// }))
	r.Engine.HidePort = true
	r.Engine.HideBanner = true
	r.Engine.Use(middleware.Recover())

	var origins = []string{
		"http://localhost:9090",
		"http://localhost:3000",
	}

	// if r.cfg.App.Env == "local" {
	// 	origins = []string{
	// 		"http://localhost:9090",
	// 		"http://localhost:3000",
	// 	}
	// } else if r.cfg.App.Env == "dev" {
	// 	origins = []string{
	// 	}
	// } else if r.cfg.App.Env == "prd" {
	// 	origins = []string{
	// 	}
	// }

	fmt.Println("------------")
	fmt.Println(r.cfg.App.Env)
	fmt.Println(origins)
	fmt.Println("------------")

	r.Engine.Use(middleware.CORSWithConfig((middleware.CORSConfig{
		AllowOrigins: origins,
		AllowHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderContentType,
			echo.HeaderOrigin,
			echo.HeaderAccessControlAllowOrigin,
			"FirebaseAuthorization",
		},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	})))

	r.Engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "healthz") {
				return true
			} else {
				return false
			}
		},
	}))

	api := r.Engine.Group("")
	{
		api.GET("/healthz", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		api.GET("/*", func(c echo.Context) error {
			return c.NoContent(http.StatusNotFound)
		})

		api.POST("/*", func(c echo.Context) error {
			return c.NoContent(http.StatusNotFound)
		})

		api.PUT("/*", func(c echo.Context) error {
			return c.NoContent(http.StatusNotFound)
		})
	}

	/****************************************************************************************/
	/// No Auth API
	//

	userRoutes := routes.UserRoutes{}

	noAuthAPI := api.Group("api")
	{
		noAuthAPI.GET("/healthz", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		// ユーザーの新規登録
		noAuthAPI.POST("/signup", userRoutes.SignUp(db, firebase))

		// ユーザーのログイン
		noAuthAPI.PUT("/signin", userRoutes.SignIn(db, firebase))

		noAuthAPI.POST("/detect/resume/job_seeker", userRoutes.DetectTextFromJobSeekerResume(db, firebase))

	}

	/****************************************************************************************/
	/// UserAPI
	//
	// userAPI := noAuthAPI.Group("/user")
	{
		// ユーザーのログイン

	}

	return r
}

func (r *Router) Start() {
	r.Engine.Start(fmt.Sprintf(":%d", 9090))
}
