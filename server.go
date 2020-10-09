package gofm

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Manager
}

// handleIncAdcs increase audiences 增加在线人数
func (s *Server) handleIncAdcs() func(ctx gin.Context) {
	return func(ctx gin.Context) {
		roomID := ctx.GetInt("room_id")
		nums := ctx.GetInt("nums")
		adcs, err := s.Manager.IncreaseAudience(roomID, nums)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &LRAdcsResponse{CurAdcs: adcs, Message: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, &LRAdcsResponse{CurAdcs: adcs, Message: ""})
	}
}

// handleDecAdcs decrease audiences 减少在线人数
func (s *Server) handleDecAdcs() func(ctx gin.Context) {
	return func(ctx gin.Context) {
		roomID := ctx.GetInt("room_id")
		nums := ctx.GetInt("nums")
		adcs, err := s.Manager.DecreaseAudience(roomID, nums)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &LRAdcsResponse{CurAdcs: adcs, Message: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, &LRAdcsResponse{CurAdcs: adcs, Message: ""})
	}
}

func (s *Server) Run() {
	e := gin.Default()
	liveRoomGroup := e.Group("/api")
	liveRoomGroup.Handle(http.MethodPut, "/incauds", func(context *gin.Context) {})
	liveRoomGroup.Handle(http.MethodPut, "/decadcs")

	e.Run(":3016")
}

func NewServer() *Server {
	return &Server{
		Manager: &manager{
			mu:    nil,
			rooms: nil,
		},
	}
}

type LRAdcsResponse struct {
	Message string `json:"message"`
	CurAdcs int    `json:"cur_adcs"`
}
