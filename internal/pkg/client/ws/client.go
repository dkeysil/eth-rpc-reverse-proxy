package client

import (
	"sync"

	"github.com/dgrr/websocket"
	resolver "github.com/dkeysil/eth-rpc-reverse-proxy/internal/id_resolver"
	ws "github.com/fasthttp/websocket"
	"github.com/tidwall/sjson"
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

		go client.listener(backendConn, host)
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
