package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"learn_go/webook/internal/service"
	"log"
	"math/rand"
	"net/http"
)

type SMSHandler struct {
	codeService service.CodeService
}

func NewSMSHandler(codeService service.CodeService) *SMSHandler {
	return &SMSHandler{codeService}
}

func (h *SMSHandler) Send(ctx *gin.Context) {
	type SMSReq struct {
		Phone string `json:"phone"`
	}
	req := &SMSReq{}
	if err := ctx.Bind(req); err != nil {
		return
	}

	// 这里可以验证手机号的格式
	if req.Phone == "" {
		ctx.String(http.StatusOK, "手机号不能为空")
		return
	}

	code := h.generateCode()
	log.Printf("send code: %v\n", code)
	err := h.codeService.Send(ctx, "login", req.Phone, code)
	switch err {
	case nil:
		ctx.String(http.StatusOK, "success")
	case service.ErrTooManySend:
		ctx.String(http.StatusOK, "验证码发送太多次了")
	default:
		log.Printf("sms send error: %v\n", err)
		ctx.String(http.StatusOK, "系统错误")
	}

}

func (h *SMSHandler) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func (h *SMSHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/sms/send", h.Send)
}
