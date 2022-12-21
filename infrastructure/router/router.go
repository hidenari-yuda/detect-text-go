package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/domain/utility"
	"github.com/hidenari-yuda/paychan-server/infrastructure/database"
	"github.com/hidenari-yuda/paychan-server/infrastructure/driver"
	"github.com/hidenari-yuda/paychan-server/infrastructure/router/routes"
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
		db        = database.NewDB(r.cfg.DB, true)
		firebase  = driver.NewFirebaseImpl(r.cfg.Firebase)
		basicAuth = utility.NewBasicAuth(r.cfg)
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
		"http://localhost:8080",
		"https://paychan.jp",
		"https://app.paychan.jp",
		"https://api.paychan.jp",
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

	var (
		userRoutes     = routes.UserRoutes{}
		presentRoutes  = routes.PresentRoutes{}
		richMenuRoutes = routes.RichMenuRoutes{}
	)

	noAuthAPI := api.Group("api")
	{
		noAuthAPI.GET("/healthz", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		// ユーザーの新規登録
		noAuthAPI.POST("/signup", userRoutes.SignUp(db, firebase))

		// // ユーザーのログイン
		// noAuthAPI.PUT("/signin", userRoutes.SignIn(db, firebase))

		noAuthAPI.POST("/line", userRoutes.GetLineWebHook(db, firebase))

		// noAuthAPI.POST("/present", presentRoutes.Create(db, firebase))

		// noAuthAPI.DELETE("/expired", presentRoutes.DeleteByExpired(db, firebase))

	}

	/****************************************************************************************/
	/// UserAPI
	//
	// userAPI := noAuthAPI.Group("/user")
	{

		// firebaseTokenからユーザーを取得
		// userAPI.GET("/firebaseToken", userRoutes.GetByFirebaseToken(db, firebase))

		// ユーザーのログイン
		// userAPI.GET("/lineUserId", userRoutes.GetByLineUserId(db, firebase))

	}

	/****************************************************************************************/
	/// PresentAPI
	//
	// presentAPI := noAuthAPI.Group("/present")
	// {

	// // create
	// presentAPI.POST("", presentRoutes.Create(db, firebase))

	// // update
	// presentAPI.PUT("", presentRoutes.Update(db, firebase))

	// // getbyId
	// presentAPI.GET("/:id", presentRoutes.GetById(db, firebase))

	// // getbyLineUserId
	// presentAPI.GET("/lineUserId/:lineUserId", presentRoutes.GetByLineUserId(db, firebase))

	// }

	/****************************************************************************************/
	/// AdminAPI
	//
	adminAPI := r.Engine.Group("/admin")
	adminNoAuthAPI := adminAPI.Group("")
	{
		adminNoAuthAPI.PUT("/authorize", routes.AdminAuthorize(db, r.cfg.App))
	}

	// Admin basic auth
	adminBasicAuthAPI := adminAPI.Group("/basic")
	adminBasicAuthAPI.Use(middleware.BasicAuth(basicAuth.BasicAuthValidator))
	{
		adminBasicAuthAPI.GET("/healthz", func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

	}
	/****************************************************************************************/
	/// UserAPI
	//
	adminForUserAPI := adminAPI.Group("/user")
	adminForUserAPI.Use(middleware.BasicAuth(basicAuth.BasicAuthValidator))
	{
		// getbyLineUserId
		adminForUserAPI.GET("/lineUserId/:lineUserId", userRoutes.GetByLineUserId(db, firebase))

		adminForUserAPI.GET("/all", userRoutes.GetAll(db, firebase))
	}
	/****************************************************************************************/
	/// PresentAPI
	//

	adminForPresentAPI := adminAPI.Group("/present")
	adminForPresentAPI.Use(middleware.BasicAuth(basicAuth.BasicAuthValidator))

	{
		// create
		adminForPresentAPI.POST("", presentRoutes.Create(db, firebase))

		// update
		adminForPresentAPI.PUT("", presentRoutes.Update(db, firebase))

		// getbyId
		adminForPresentAPI.GET("/:id", presentRoutes.GetById(db, firebase))

		// getbyLineUserId
		adminForPresentAPI.GET("/lineUserId/:lineUserId", presentRoutes.GetByLineUserId(db, firebase))

		// getall
		adminForPresentAPI.GET("/all", presentRoutes.GetAll(db, firebase))

		//delete
		adminForPresentAPI.DELETE("/expired", presentRoutes.DeleteByExpired(db, firebase))
	}
	/****************************************************************************************/
	/// RichMenuAPI
	//
	adminForRichMenuAPI := adminAPI.Group("/richMenu")
	adminForRichMenuAPI.Use(middleware.BasicAuth(basicAuth.BasicAuthValidator))
	{
		// create
		adminForRichMenuAPI.POST("", richMenuRoutes.Create(db, firebase))

		// uploadImage
		adminForRichMenuAPI.POST("/uploadImage/:richMenuId/:imagePath", richMenuRoutes.UploadImage(db, firebase))

		// createAlias
		adminForRichMenuAPI.POST("/alias/:richMenuId/:aliasId", richMenuRoutes.CreateAlias(db, firebase))

		//updateAlias
		adminForRichMenuAPI.PUT("/alias/:richMenuId/:aliasId", richMenuRoutes.UpdateAlias(db, firebase))

		//setAlias
		adminForRichMenuAPI.PUT("/setAlias/:richMenuId/:aliasId", richMenuRoutes.SetAlias(db, firebase))

		// getAll
		adminForRichMenuAPI.GET("/all", richMenuRoutes.GetAll(db, firebase))

		//deleteRichMenu
		adminForRichMenuAPI.DELETE("/:richMenuId", richMenuRoutes.DeleteRichMenu(db, firebase))

		//deleteAlias
		adminForRichMenuAPI.DELETE("/alias/:aliasId", richMenuRoutes.DeleteAlias(db, firebase))
	}

	/****************************************************************************************/

	return r
}

func (r *Router) Start() {
	r.Engine.Start(fmt.Sprintf(":%d", 8080))
}
