package controllers

import (
	"encoding/json"

	"github.com/dgrr/websocket"
	"go.uber.org/zap"
)

type Method struct {
	Method string `json:"method"`
}

type ID struct {
	ID uint64 `json:"id"`
}

func (s *Service) OnMessage(c *websocket.Conn, isBinary bool, data []byte) {
	zap.L().Info("got message from clientConn", zap.ByteString("message", data))

	var method Method
	err := json.Unmarshal(data, &method)
	if err != nil {
		zap.L().Error("error while getting method", zap.Error(err))
	}

	var originalID ID
	err = json.Unmarshal(data, &originalID)
	if err != nil {
		zap.L().Error("error while unmarshaling", zap.Error(err))
		return
	}

	id := s.IDResolver.SetID(originalID.ID, c.ID())

	err = s.Client.Send(c, data, s.BackendResolver.GetUpstreamHost("ws:*"), id)
	if err != nil {
		zap.L().Error(err.Error())
	}

	if method.Method == "eth_call" {
		err = s.Client.Send(c, data, s.BackendResolver.GetUpstreamHost("ws:eth_call"), id)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
}
