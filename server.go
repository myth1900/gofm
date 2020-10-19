package gofm

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	Manager
	e *gin.Engine
}

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Next()
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}
func NewServer() *Server {
	e := gin.Default()
	e.Use(Cors())
	pprof.Register(e, "/debug/pprof")

	s := &Server{
		Manager: NewManager(),
		e:       e,
	}
	s.handleStaticFiles()
	s.handleRoutes()
	return s
}

func (s *Server) Run() {
	s.e.Run(":3016")
}

type Route struct {
	Method string
	Path   string
	Handle func(ctx *gin.Context)
}

func (s *Server) handleRoutes() {
	routes := []Route{
		{
			http.MethodPut, "/api/room/:room_id/audience", s.handleAudience,
		},
		{
			http.MethodGet, "/api/rooms/status", s.handleStatus,
		},
	}

	for _, r := range routes {
		s.e.Handle(r.Method, r.Path, r.Handle)
	}
}

func (s *Server) handleStaticFiles() {
	infoFunc := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	fsCss := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: infoFunc,
		Prefix:    "web/dist/css",
		Fallback:  "index.html",
	}
	fsJS := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: infoFunc,
		Prefix:    "web/dist/js",
		Fallback:  "index.html",
	}
	//fs := &assetfs.AssetFS{
	//	Asset:     Asset,
	//	AssetDir:  AssetDir,
	//	AssetInfo: infoFunc,
	//	Prefix:    "web/dist/",
	//	Fallback:  "index.html",
	//}
	s.e.StaticFS("/css", fsCss)
	s.e.StaticFS("/js", fsJS)
	//s.e.StaticFS("/", fs)

	s.e.GET("/", func(ctx *gin.Context) {
		ctx.Writer.WriteHeader(http.StatusOK)
		idxHtml, _ := Asset("web/dist/index.html")
		ctx.Writer.Write(idxHtml)
		ctx.Writer.Header().Add("Accept", "text/html")
		ctx.Writer.Flush()
	})
	s.e.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.Writer.WriteHeader(http.StatusOK)
		fav, _ := Asset("web/dist/favicon.ico")
		ctx.Writer.Write(fav)
		ctx.Writer.Flush()

	})
}

func (s *Server) handleAudience(ctx *gin.Context) {
	if err := func() error {
		roomID, err := strconv.Atoi(ctx.Param("room_id"))
		if err != nil {
			return err
		}
		nums, err := strconv.Atoi(ctx.Query("nums"))
		if err != nil {
			return err
		}
		if err := s.UpdateAudienceWithRoomID(roomID, nums); err != nil {
			return err
		}
		return nil
	}(); err != nil {
		ctx.JSON(http.StatusInternalServerError, &Response{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &Response{Message: ""})
}
func (s *Server) handleStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, s.Status())
}

type Response struct {
	Message string `json:"message"`
}
