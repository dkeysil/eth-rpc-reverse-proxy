package client

import (
	"encoding/json"
	"errors"

	"github.com/dgrr/websocket"
	resolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/id_resolver"
	ws "github.com/fasthttp/websocket"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

/*
1. read message and recover connection if it's down
2. getting id from resolver
3. updating message (set original id)
4. getting client connection
5. writing to client connection
*/
func (wsc *wsReverseProxyClient) listener(backendConn *ws.Conn, host string) {
	for {
		message, err := wsc.readMessage(backendConn, host)
		if err != nil {
			// TODO: remove host which cannot be retried from list of upstreams
			panic(err)
		}

		zap.L().Debug("got message from backendConn", zap.ByteString("message", message), zap.String("host", host))

		id, err := wsc.getID(message)
		if err != nil {
			zap.L().Debug(
				"problems in listener while getting id",
				zap.Error(err),
				zap.ByteString("message", message),
				zap.String("host", host),
			)
			continue
		}

		err = wsc.sendMessage(message, id)
		if err != nil {
			wsc.clientPool.Delete(id.ClientID)
			zap.L().Error(
				"error while sending message",
				zap.Error(err),
				zap.Uint64("client_id", id.ClientID),
				zap.ByteString("original_id", id.OriginalID),
				zap.Uint64("request_id", id.RequestID),
			)
		}
	}
}

func (wsc *wsReverseProxyClient) readMessage(backendConn *ws.Conn, host string) ([]byte, error) {
	_, message, err := backendConn.ReadMessage()
	if err != nil {
		zap.L().Debug("error in listener", zap.Error(err), zap.String("host", host))

		backendConn, _, err := ws.DefaultDialer.Dial(host, nil)
		if err != nil {
			zap.L().Error("can't revive backend connection", zap.Error(err), zap.String("host", host))
			return nil, err
		}

		zap.L().Debug("revive backend connection", zap.String("host", host))
		wsc.backendPool.Store(host, backendConn)
	}

	return message, err
}

func (wsc *wsReverseProxyClient) getID(message []byte) (resolver.ID, error) {
	var requestID ID
	err := json.Unmarshal(message, &requestID)
	if err != nil {
		return resolver.ID{}, err
	}

	id, ok := wsc.idResolver.PopID(requestID.ID)
	if !ok {
		return resolver.ID{}, errors.New("duplicated reauest")
	}

	return id, nil
}

func (wsc *wsReverseProxyClient) sendMessage(message []byte, id resolver.ID) error {
	message, err := sjson.SetBytes(message, "id", id.OriginalID)
	if err != nil {
		return err
	}

	conn, ok := wsc.clientPool.Load(id.ClientID)
	if !ok {
		return nil
	}

	clientConn := conn.(*websocket.Conn)

	_, err = clientConn.Write(message)
	if err != nil {
		return err
	}

	return nil
}
