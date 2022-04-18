package client

import (
	"encoding/json"
	"sync"

	"github.com/dgrr/websocket"
	resolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/id_resolver"
	ws "github.com/fasthttp/websocket"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

/*
TODO:
1. duplicate call on eth_call
2. keep alive backend upstreams
3. refactor
*/

type WSReverseProxyClient interface {
	Send(c *websocket.Conn, data []byte, host string, id resolver.ID) (err error)
}

type wsReverseProxyClient struct {
	clientPool  sync.Map
	backendPool sync.Map
	idResolver  resolver.IDResolver
}

func NewWSReverseProxyClient(upstreams []string, idResolver resolver.IDResolver) WSReverseProxyClient {
	client := &wsReverseProxyClient{idResolver: idResolver}
	for _, host := range upstreams {
		backendConn, _, err := ws.DefaultDialer.Dial(host, nil)
		if err != nil {
			panic(err)
		}
		client.backendPool.Store(host, backendConn)

		go client.listener(backendConn)
	}

	return client
}

func (wsc *wsReverseProxyClient) Send(clientConn *websocket.Conn, data []byte, host string, id resolver.ID) (err error) {
	conn, ok := wsc.backendPool.Load(host)
	if !ok {
		panic("can't get backend conn by host")
	}

	backendConn := conn.(*ws.Conn)

	wsc.clientPool.LoadOrStore(id.ClientID, clientConn)
	data, err = sjson.SetBytes(data, "id", id.RequestID)
	if err != nil {
		return err
	}

	err = backendConn.WriteMessage(ws.TextMessage, data)
	if err != nil {
		return err
	}

	return nil
}

type ID struct {
	ID uint64 `json:"id"`
}

func (wsc *wsReverseProxyClient) listener(backendConn *ws.Conn) {
	for {
		_, message, err := backendConn.ReadMessage()
		if err != nil {
			zap.L().Error("error in listener", zap.Error(err))
			return
		}
		zap.L().Info("got message from backendConn", zap.ByteString("message", message))

		var requestID ID
		err = json.Unmarshal(message, &requestID)
		if err != nil {
			zap.L().Error("error while unmarshaling", zap.Error(err))
			continue
		}

		id, ok := wsc.idResolver.PopID(requestID.ID)
		if !ok {
			zap.L().Debug("got duplicated request")
			continue
		}

		conn, ok := wsc.clientPool.Load(id.ClientID)
		if !ok {
			zap.L().Error("can't get client conn", zap.Uint64("client_id", id.ClientID), zap.Uint64("original_id", id.OriginalID), zap.Uint64("request_id", requestID.ID))
			continue
		}

		clientConn := conn.(*websocket.Conn)

		message, err = sjson.SetBytes(message, "id", id.OriginalID)
		if err != nil {
			zap.L().Error("error while setting original id to message", zap.Error(err), zap.Uint64("client_id", id.ClientID), zap.Uint64("original_id", id.OriginalID), zap.Uint64("request_id", requestID.ID))
			wsc.clientPool.Delete(id.ClientID)
			return
		}

		_, err = clientConn.Write(message)
		if err != nil {
			zap.L().Error("error while writing to clientConn", zap.Error(err), zap.Uint64("client_id", id.ClientID), zap.Uint64("original_id", id.OriginalID), zap.Uint64("request_id", requestID.ID))
			wsc.clientPool.Delete(id.ClientID)
			return
		}
	}
}
