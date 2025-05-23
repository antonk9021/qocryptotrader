package alphapoint

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/antonk9021/qocryptotrader/encoding/json"
	exchange "github.com/antonk9021/qocryptotrader/exchanges"
	"github.com/antonk9021/qocryptotrader/log"
)

const (
	alphapointDefaultWebsocketURL = "wss://sim3.alphapoint.com:8401/v1/GetTicker/"
)

// WebsocketClient starts a new webstocket connection
func (a *Alphapoint) WebsocketClient() {
	for a.Enabled {
		var dialer websocket.Dialer
		var err error
		var httpResp *http.Response
		endpoint, err := a.API.Endpoints.GetURL(exchange.WebsocketSpot)
		if err != nil {
			log.Errorln(log.WebsocketMgr, err)
		}
		a.WebsocketConn, httpResp, err = dialer.Dial(endpoint, http.Header{})
		httpResp.Body.Close() // not used, so safely free the body

		if err != nil {
			log.Errorf(log.ExchangeSys, "%s Unable to connect to Websocket. Error: %s\n", a.Name, err)
			continue
		}

		if a.Verbose {
			log.Debugf(log.ExchangeSys, "%s Connected to Websocket.\n", a.Name)
		}

		err = a.WebsocketConn.WriteMessage(websocket.TextMessage, []byte(`{"messageType": "logon"}`))

		if err != nil {
			log.Errorln(log.ExchangeSys, err)
			return
		}

		for a.Enabled {
			msgType, resp, err := a.WebsocketConn.ReadMessage()
			if err != nil {
				log.Errorln(log.ExchangeSys, err)
				break
			}

			if msgType == websocket.TextMessage {
				type MsgType struct {
					MessageType string `json:"messageType"`
				}

				msgType := MsgType{}
				err := json.Unmarshal(resp, &msgType)
				if err != nil {
					log.Errorln(log.ExchangeSys, err)
					continue
				}

				if msgType.MessageType == "Ticker" {
					ticker := WebsocketTicker{}
					err = json.Unmarshal(resp, &ticker)
					if err != nil {
						log.Errorln(log.ExchangeSys, err)
						continue
					}
				}
			}
		}
		a.WebsocketConn.Close()
		log.Debugf(log.ExchangeSys, "%s Websocket client disconnected.", a.Name)
	}
}
