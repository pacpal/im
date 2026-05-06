package api

import (
	"IM/server/gateway/auth"
	msgService "IM/server/msgservice"
	"IM/server/msgservice/hub"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowed := []string{"http://localhost:8080"}
		for _, allow := range allowed {
			if origin == allow {
				return true
			}
		}
		log.Printf("wrong origin:%s", &origin)
		return false
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsHandler struct {
	Hub        *hub.Hub
	msgService *msgService.MessageService
}

func NewWsHandler(s *msgService.MessageService) *wsHandler {
	return &wsHandler{msgService: s}
}

func (ws *wsHandler) HandleWs(c *gin.Context) {
	// 从cookie查token鉴权
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "'missing token'"})
		return
	}
	token := cookie.Value

	claims, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization"})
		return
	}
	c.Set("uid", claims["sub"])
	c.Set("userName", claims["name"])

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrader wrong", err)
		return
	}
	defer conn.Close()
	userID := claims["sub"].(string)

	client := ws.Hub.Register(conn, userID)

	client.Start()

}
