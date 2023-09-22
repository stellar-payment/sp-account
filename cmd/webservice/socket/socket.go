package socket

import (
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	tagLoggerHandleLiveAPI = "[HandleLiveAPI]"
)

type SocketParams struct {
	Logger *logrus.Entry
}

// Wrapper to allow CORS on Socket.IO
type NextSocketIO struct {
	Server *socketio.Server
	Logger *logrus.Entry
}

func (nsio NextSocketIO) ServeHTTP(c echo.Context) error {
	nsio.Server.ServeHTTP(c.Response(), c.Request())
	return nil
}

func HandleLiveAPI(params *SocketParams) (nsio *NextSocketIO) {
	allowAllOrigins := func(r *http.Request) bool { return true }

	nsio = &NextSocketIO{
		Server: socketio.NewServer(&engineio.Options{
			Transports: []transport.Transport{
				&polling.Transport{CheckOrigin: allowAllOrigins},
				&websocket.Transport{CheckOrigin: allowAllOrigins},
			},
		}),
	}

	// onConnect Flat
	nsio.Server.OnConnect(defaultNamespace, func(c socketio.Conn) error {
		c.SetContext("") // TODO: Decided what to put
		params.Logger.Println("--> connected:", c.ID())
		return nil
	})

	// onConnect Chat
	return
}
