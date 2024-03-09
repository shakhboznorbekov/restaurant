package socket

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/websocket"
//	"github.com/restaurant/internal/auth"
//
//	//"gitlab.com/m9693/oxo/oxo-backend/websocket_service/api/models"
//	"log"
//	"time"
//)
//
//const (
//	// Time allowed to write a message to the peer.
//	writeWait = 10 * time.Second
//
//	// Time allowed to read the next pong message from the peer.
//	pongWait = 60 * time.Second
//
//	// Send pings to peer with this period. Must be less than pongWait.
//	pingPeriod = (pongWait * 9) / 10
//
//	// Maximum message size allowed from peer.
//	maxMessageSize = 512
//)
//
//// Client is a middleman between the _websocket connection and the hub.
//type Client struct {
//	hub *Hub
//
//	// The _websocket connection.
//	conn *websocket.Conn
//
//	Send chan ResMessage
//
//	// close read
//	closeRead chan bool
//
//	Search string
//
//	// close read
//	closeWrite chan bool
//
//	//waiter_new_order
//	NewOrder int
//
//	//user map status
//	mapStatus bool
//
//	auth.ClaimsAuth
//}
//
//// readPump pumps messages from the _websocket connection to the hub.
////
//// The application runs readPump in a per-connection goroutine. The application
//// ensures that there is at most one reader on a connection by executing all
//// reads from this goroutine.
//func (c *Client) readPump() {
//	closeR := true
//	defer func() {
//		if closeR {
//			c.hub.unregister <- c
//			c.conn.Close()
//		}
//	}()
//	c.conn.SetReadLimit(maxMessageSize)
//	c.conn.SetReadDeadline(time.Now().Add(pongWait))
//	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
//	for {
//		var action = Action{}
//		_, m, err := c.conn.ReadMessage()
//		json.Unmarshal(m, &action)
//
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("error: %v", err)
//			}
//			break
//		}
//		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
//		//data := map[string][]byte{
//		//	"message": message,
//		//	"id":      []byte(c.ID),
//		//}
//		//userMessage, _ := json.Marshal(data)
//		res := Message{
//			UserID:  c.ID,
//			Action:  action.Action,
//			Message: m,
//		}
//		c.hub.Broadcast <- res
//	}
//	select {
//	case <-c.closeRead:
//		closeR = false
//		return
//	default:
//		return
//	}
//}
//
//// writePump pumps messages from the hub to the _websocket connection.
////
//// A goroutine running writePump is started for each connection. The
//// application ensures that there is at most one writer to a connection by
//// executing all writes from this goroutine.
//func (c *Client) writePump() {
//	ticker := time.NewTicker(pingPeriod)
//	closeW := true
//	defer func() {
//		if closeW {
//			ticker.Stop()
//			c.conn.Close()
//		}
//	}()
//	for {
//		select {
//		case message, ok := <-c.Send:
//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if !ok {
//				// The hub closed the channel.
//				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
//				return
//			}
//
//			w, err := c.conn.NextWriter(websocket.TextMessage)
//			if err != nil {
//				return
//			}
//			w.Write(message.Message)
//			if err := w.Close(); err != nil {
//				return
//			}
//			//Add queued chat messages to the current _websocket message.
//			n := len(c.Send)
//			for i := 0; i < n; i++ {
//				w, err := c.conn.NextWriter(websocket.TextMessage)
//				if err != nil {
//					return
//				}
//				message := <-c.Send
//				w.Write(message.Message)
//				if err := w.Close(); err != nil {
//					return
//				}
//			}
//		case <-ticker.C:
//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
//				return
//			}
//		}
//	}
//	select {
//	case <-c.closeWrite:
//		closeW = false
//		return
//	default:
//		return
//	}
//}
//
//// ServeWs handles _websocket requests from the peer.
//func ServeWs(hub *Hub, conn *websocket.Conn, user auth.ClaimsAuth) {
//	client := &Client{hub: hub, conn: conn, Send: make(chan ResMessage, 256)}
//	client.ClaimsAuth = user
//
//	client.closeWrite = make(chan bool, 2)
//	client.closeRead = make(chan bool, 2)
//	key := ""
//	if client.BranchID != nil {
//		key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//	} else if client.RestaurantID != nil {
//		key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//	} else {
//		key = client.Roles
//	}
//	if _, ok := hub.Clients[key]; ok {
//		if v, ok := hub.Clients[key][client.ID]; ok {
//			v.closeRead <- true
//			v.closeWrite <- true
//			v.conn.Close()
//			hub.Clients[key][client.ID] = client
//		} else {
//			client.hub.register <- client
//		}
//	} else {
//		client.hub.register <- client
//	}
//
//	// Allow collection of memory referenced by the caller by doing all work in
//	// new goroutines.
//	go client.writePump()
//	go client.readPump()
//}
//
////
////func (c *Client) writeMap(searchChan *chan string) {
////	defer func() {
////		c.conn.Close()
////	}()
////	search := ""
////	for {
////		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
////		w, err := c.conn.NextWriter(websocket.TextMessage)
////		if err != nil {
////			return
////		}
////		select {
////		case s := <-*searchChan:
////			search = s
////		default:
////		}
////		message, err := c.hub.userHub.OnLineUserLocation(c.OrganizationID, search)
////		if err != nil {
////			er, _ := json.Marshal(map[string]string{
////				"error": err.Error(),
////			})
////			w.Write(er)
////		} else {
////		}
////		if err := w.Close(); err != nil {
////			return
////		}
////		time.Sleep(5 * time.Second)
////	}
////}
////
////func (c *Client) readMap(searchChan *chan string) {
////	defer func() {
////		c.conn.Close()
////	}()
////	c.conn.SetReadLimit(maxMessageSize)
////	c.conn.SetReadDeadline(time.Now().Add(pongWait))
////	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
////	for {
////		//var message = ReqMessage{}
////		_, m, err := c.conn.ReadMessage()
////		json.Unmarshal(m, &message)
////
////		if err != nil {
////			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
////				log.Printf("error: %v", err)
////			}
////			break
////		}
////
////		*searchChan <- "message.Search"
////	}
////}
