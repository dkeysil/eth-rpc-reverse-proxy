package client

import (
	"encoding/json"
	"sync"

	"github.com/dgrr/websocket"
	ws "github.com/fasthttp/websocket"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

type WSReverseProxyClient interface {
	Send(c *websocket.Conn, data []byte, host string) (err error)
}

type wsReverseProxyClient struct {
	clientPool   sync.Map
	clientIDPool sync.Map
	backendPool  sync.Map
}

func NewWSReverseProxyClient(upstreams []string) WSReverseProxyClient {
	client := &wsReverseProxyClient{}
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

func (wsc *wsReverseProxyClient) Send(clientConn *websocket.Conn, data []byte, host string) (err error) {
	conn, ok := wsc.backendPool.Load(host)
	if !ok {
		panic("can't get backend conn by host")
	}

	backendConn := conn.(*ws.Conn)

	wsc.clientPool.LoadOrStore(clientConn.ID(), clientConn)

	var id ID
	err = json.Unmarshal(data, &id)
	if err != nil {
		zap.L().Error("error while unmarshaling", zap.Error(err))
		return err
	}

	wsc.clientIDPool.Store(clientConn.ID(), id.ID)

	data, err = sjson.SetBytes(data, "id", clientConn.ID())
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
			// if error have to close both connections
			return
		}
		zap.L().Info("got message from backendConn", zap.ByteString("message", message))

		var id ID
		err = json.Unmarshal(message, &id)
		if err != nil {
			zap.L().Error("error while unmarshaling", zap.Error(err))
			continue
		}

		conn, ok := wsc.clientPool.Load(id.ID)
		if !ok {
			zap.L().Error("can't get client conn", zap.Uint64("id", id.ID))
			continue
		}

		clientConn := conn.(*websocket.Conn)

		tempID, _ := wsc.clientIDPool.Load(clientConn.ID())
		originalID := tempID.(uint64)

		message, err = sjson.SetBytes(message, "id", originalID)
		if err != nil {
			zap.L().Error("error while setting original id to message", zap.Error(err), zap.Uint64("original_id", originalID), zap.Uint64("id", id.ID))
			wsc.clientPool.Delete(id.ID)
			wsc.clientIDPool.Delete(clientConn.ID())
			return
		}

		_, err = clientConn.Write(message)
		if err != nil {
			zap.L().Error("error while writing to clientConn", zap.Error(err), zap.Uint64("id", id.ID))
			wsc.clientPool.Delete(id.ID)
			wsc.clientIDPool.Delete(clientConn.ID())
			return
		}
	}
}
