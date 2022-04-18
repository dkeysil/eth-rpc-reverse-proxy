package controllers

import (
	"github.com/dgrr/websocket"
	"go.uber.org/zap"
)

func (s *Service) OnMessage(c *websocket.Conn, isBinary bool, data []byte) {
	zap.L().Info("got message from clientConn", zap.ByteString("message", data))
	err := s.Client.Send(c, data, s.BackendResolver.GetUpstreamHost("ws:*"))
	if err != nil {
		zap.L().Error(err.Error())
	}
}
