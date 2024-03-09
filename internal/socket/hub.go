package socket

//
//import (
//	"fmt"
//	"sync"
//)
//
//// Hub maintains the set of active clients and broadcasts messages to the
//type Hub struct {
//	// Registered clients.
//	Clients map[string]map[int64]*Client
//
//	// Registered clients.
//	printer map[int64]*Printer
//
//	// Inbound messages from the clients.
//	Broadcast chan Message
//
//	// Register requests from the clients.
//	register chan *Client
//
//	// Unregister requests from clients.
//	unregister chan *Client
//
//	// Register requests from the printer.
//	registerPrinter chan *Printer
//
//	// Unregister requests from printer.
//	unregisterPrinter chan *Printer
//}
//
//// NewHub ...
//func NewHub() *Hub {
//	return &Hub{
//		Broadcast:         make(chan Message, 256),
//		register:          make(chan *Client, 256),
//		unregister:        make(chan *Client, 256),
//		Clients:           make(map[string]map[int64]*Client),
//		registerPrinter:   make(chan *Printer, 256),
//		unregisterPrinter: make(chan *Printer, 256),
//		printer:           make(map[int64]*Printer),
//	}
//}
//
//// Run ...
//func (h *Hub) Run() {
//	var wg sync.WaitGroup
//	wg.Add(3)
//	go func() {
//		defer wg.Done()
//		for {
//			select {
//			case client := <-h.registerPrinter:
//				h.printer[client.BranchID] = client
//				for i := 0; i < len(h.register); i++ {
//					client := <-h.registerPrinter
//					h.printer[client.BranchID] = client
//				}
//			}
//		}
//	}()
//	go func() {
//		wg.Done()
//		for {
//			select {
//			case client := <-h.unregisterPrinter:
//				delete(h.printer, client.BranchID)
//				close(client.send)
//				for i := 0; i < len(h.register); i++ {
//					client := <-h.registerPrinter
//					delete(h.printer, client.BranchID)
//					close(client.send)
//				}
//			}
//		}
//	}()
//	go func() {
//		defer wg.Done()
//		for {
//			select {
//			case client := <-h.register:
//				key := ""
//				if client.BranchID != nil {
//					key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//				} else if client.RestaurantID != nil {
//					key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//				} else {
//					key = client.Roles
//				}
//				if _, ok := h.Clients[key]; ok {
//					h.Clients[key][client.ID] = client
//				} else {
//					h.Clients[key] = make(map[int64]*Client)
//					h.Clients[key][client.ID] = client
//				}
//				for i := 0; i < len(h.register); i++ {
//					client := <-h.register
//					key := ""
//					if client.BranchID != nil {
//						key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//					} else if client.RestaurantID != nil {
//						key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//					} else {
//						key = client.Roles
//					}
//					h.Clients[key][client.ID] = client
//				}
//			}
//		}
//	}()
//	go func() {
//		wg.Done()
//		for {
//			select {
//			case client := <-h.unregister:
//				key := ""
//				if client.BranchID != nil {
//					key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//				} else if client.RestaurantID != nil {
//					key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//				} else {
//					key = client.Roles
//				}
//				delete(h.Clients[key], client.ID)
//				close(client.Send)
//				for i := 0; i < len(h.register); i++ {
//					client := <-h.register
//					key := ""
//					if client.BranchID != nil {
//						key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//					} else if client.RestaurantID != nil {
//						key = fmt.Sprintf("%s_%d", client.Roles, *client.BranchID)
//					} else {
//						key = client.Roles
//					}
//					delete(h.Clients[key], client.ID)
//					close(client.Send)
//				}
//			}
//		}
//	}()
//	go func() {
//		wg.Done()
//		for {
//			select {
//			case Message := <-h.Broadcast:
//				h.controlMessage(Message)
//				for i := 0; i < len(h.register); i++ {
//					Message := <-h.Broadcast
//					h.controlMessage(Message)
//				}
//			}
//		}
//
//	}()
//	wg.Wait()
//}
//
//func (h *Hub) controlMessage(message Message) {
//	switch message.Action {
//	case "new_order":
//		if _, ok := h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)] {
//				h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//		if _, ok := h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)] {
//				h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//		if _, ok := h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)] {
//				h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//	case "new_food":
//		if _, ok := h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)] {
//				h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//		if _, ok := h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)] {
//				h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//		if _, ok := h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)] {
//				h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//	case "order_payment":
//		if _, ok := h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)] {
//				h.Clients[fmt.Sprintf("ADMIN_%d", message.RestaurantID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//		if _, ok := h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)] {
//				h.Clients[fmt.Sprintf("BRANCH_%d", message.BranchID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//		if _, ok := h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)] {
//				h.Clients[fmt.Sprintf("CASHIER_%d", message.BranchID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//		if _, ok := h.Clients["CLIENT"]; ok {
//			if v, ok := h.Clients["CLIENT"][message.UserID]; ok {
//				v.Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//	case "new_waiter_order":
//		fmt.Println(111, h.Clients, "222:", fmt.Sprintf("WAITER_%d", message.BranchID))
//		if _, ok := h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)] {
//				h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//	case "client_accepted_order":
//		if _, ok := h.Clients["CLIENT"]; ok {
//			if v, ok := h.Clients["CLIENT"][message.UserID]; ok {
//				v.Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//	case "printer_new_food":
//		if v, ok := h.printer[message.BranchID]; ok {
//			v.send <- ResMessage{
//				Message: message.Message,
//			}
//		}
//	case "waiter_accepted_order":
//		if _, ok := h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)]; ok {
//			for k, _ := range h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)] {
//				h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)][k].Send <- ResMessage{
//					Message: message.Message,
//				}
//			}
//		}
//	case "client_call":
//		if _, ok := h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)]; ok {
//			for k, v := range h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)] {
//				if v.RestaurantID != nil {
//					h.Clients[fmt.Sprintf("WAITER_%d", message.BranchID)][k].Send <- ResMessage{
//						Message: message.Message,
//					}
//				}
//			}
//		}
//	}
//
//}
