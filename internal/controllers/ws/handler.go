package controllers

import (
	"github.com/dgrr/websocket"
	"github.com/dkeysil/eth-rpc-reverse-proxy/internal/entities"
	"go.uber.org/zap"
)

func (s *Service) OnMessage(c *websocket.Conn, isBinary bool, data []byte) {
	zap.L().Info("got message from clientConn", zap.ByteString("message", data))

	request, err := entities.NewRequest(data)
	if err != nil {
		zap.L().Error("error while unmarshaling request", zap.Error(err))
		err = c.Close()
		if err != nil {
			zap.L().Error("error while closing client connection", zap.Error(err))
		}
		return
	}

	id := s.IDResolver.SetID(request.ID, c.ID())

	err = s.Client.Send(c, data, s.BackendResolver.GetUpstreamHost("ws:*"), id)
	if err != nil {
		zap.L().Error(err.Error())
	}

	if request.Method == "eth_call" {
		err = s.Client.Send(c, data, s.BackendResolver.GetUpstreamHost("ws:eth_call"), id)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
}
