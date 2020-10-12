package gofm

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Server struct {
	Manager
	e *gin.Engine
}

func NewServer() *Server {
	e := gin.Default()
	pprof.Register(e, "/debug/pprof")

	s := &Server{
		Manager: NewManager(),
		e:       e,
	}
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
			http.MethodPut, "/api/:room_id/audience", s.handleAudience,
		},
		{
			http.MethodGet, "/debug/status", s.handleStatus,
		},
	}

	for _, r := range routes {
		s.e.Handle(r.Method, r.Path, r.Handle)
	}
}

func (s *Server) handleAudience(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Param("room_id"))
	nums, _ := strconv.Atoi(ctx.Query("nums"))
	err := s.UpdateAudienceWithRoomID(roomID, nums)
	if err != nil {
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
