package api

import (
	"IM/server/model"
	msgService "IM/server/msgservice"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MsgHandler struct {
	msgSv *msgService.MessageService
}

func NewMsgHandler(msgSv msgService.MessageService) *MsgHandler {
	return &MsgHandler{msgSv: &msgSv}
}
func (h *MsgHandler) SendMessage(c *gin.Context) {
	var req struct {
		RcID    string `json:"receive_id"`
		Content string `json:"content"`
		Type    string `json:"type"`
		Time    int64  `json:"time"`
	}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong register info"})
		return
	}
	userID, _ := c.Get("userID")
	err = h.msgSv.RouteMessage(context.Background(), model.Message{
		SdID:    userID.(string),
		RcID:    req.RcID,
		Content: req.Content,
		Type:    req.Type,
		Time:    req.Time,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "sent"})
}
func (h *MsgHandler) GetOfflineMsgs(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgs, err := h.msgSv.GetOfflineMsgs(context.Background(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

func (h *MsgHandler) GetOnlineStatus(c *gin.Context) {
	userID, _ := c.Get("userID")
	statuss, err := h.msgSv.GetOnlineStatus(context.Background(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"online": statuss})
}
