package gofm

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Server struct {
	Manager
}

func NewServer() *Server {
	return &Server{
		Manager: NewManager(),
	}
}
func (s *Server) Run() {
	e := gin.Default()
	pprof.Register(e, "/debug/pprof")
	liveRoomGroup := e.Group("/api")
	{
		liveRoomGroup.PUT("/incadcs", s.handleIncAdcs)
		liveRoomGroup.PUT("/decadcs", s.handleDecAdcs)
	}

	e.Run(":3016")
}

// handleIncAdcs increase audiences 增加在线人数
func (s *Server) handleIncAdcs(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	nums, _ := strconv.Atoi(ctx.Query("nums"))
	adcs, err := s.Manager.IncreaseAudience(roomID, nums)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &LRAdcsResponse{CurAdcs: adcs, Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &LRAdcsResponse{CurAdcs: adcs, Message: ""})
}

// handleDecAdcs decrease audiences 减少在线人数
func (s *Server) handleDecAdcs(ctx *gin.Context) {
	roomID, _ := strconv.Atoi(ctx.Query("room_id"))
	nums, _ := strconv.Atoi(ctx.Query("nums"))
	adcs, err := s.Manager.DecreaseAudience(roomID, nums)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &LRAdcsResponse{CurAdcs: adcs, Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &LRAdcsResponse{CurAdcs: adcs, Message: ""})
}

type LRAdcsResponse struct {
	Message string `json:"message"`
	CurAdcs int    `json:"cur_adcs"`
}
