package handler

import (
	"errors"
	"log"
	nethttp "net/http"
	"sync"
	"time"

	//"time"

	"github.com/gin-gonic/gin"

	"github.com/andriykutsevol/DDDCasbinExample/internal/app/application"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/request"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/response"
)

type Demos interface {
	DemoAdd(c *gin.Context)
	DemoGet(c *gin.Context)
	DemoPut(c *gin.Context)
	DemoPub(c *gin.Context)
	DemoSub(c *gin.Context)
}

type demos struct {
	demosApp application.Demos
}

func NewDemos(demosApp application.Demos) Demos {
	return &demos{
		demosApp: demosApp,
	}
}

//var wg sync.WaitGroup

func (d *demos) DemoAdd(c *gin.Context) {

	log.Println("DemoAdd")

	//ctx := c.Request.Context()
	var item request.DemoParams
	if err := http.ParseJSON(c, &item); err != nil {
		log.Println("error")
		http.ResError(c, err)
		return
	}

	resp := new(response.DemoOk)
	http.ResSuccess(c, resp)

}

//==================================================================

var (
	clients   = make(map[chan string]bool)
	clientsMu sync.Mutex
)

func broadcast(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		client <- message
	}
}

func (d *demos) DemoPub(c *gin.Context) {

	//ctx := c.Request.Context()
	var item request.DemoParams
	if err := http.ParseJSON(c, &item); err != nil {
		log.Println("error")
		http.ResError(c, err)
		return
	}

	log.Println("DemoPub: item.RequestString:", item.RequestString)
	message := item.RequestString

	if message != "" {
		broadcast(message)
		//c.JSON(http.StatusOK, gin.H{"status": "ok"})
		resp := new(response.DemoOk)
		http.ResSuccess(c, resp)
	} else {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Message is empty"})
		http.ResError(c, errors.New("error: Message is empty"))
	}

}

// // This handler handles multiple requests at the same time already
// // You just need to emulate multiple requests properly (bash/curl does not works)
// func (d *demos) DemoSub(c *gin.Context){

// 	id := c.Param("id")
// 	log.Println("DemoSub: id ", id)

//     for i := 0; i < 50; i++ {
//         log.Println("longrunning task: ", i, id)
// 		time.Sleep(100 * time.Millisecond)
//     }

// 	resp := new(response.DemoOk)
// 	resp.Message = "DemoGet OK"
// 	http.ResSuccess(c, resp)

// }

func (d *demos) DemoSub(c *gin.Context) {

	id := c.Param("id")
	log.Println("DemoSub: id ", id)

	for i := 0; i < 50; i++ {
		log.Println("longrunning task: ", i, id)
		time.Sleep(100 * time.Millisecond)
	}

	messageChan := make(chan string)
	clientsMu.Lock()
	clients[messageChan] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, messageChan)
		close(messageChan)
		clientsMu.Unlock()
	}()

	select {
	case message := <-messageChan:
		c.JSON(nethttp.StatusOK, gin.H{"message": message})
	case <-time.After(30 * time.Second): // Timeout for long-polling
		c.JSON(nethttp.StatusOK, gin.H{"message": "Timeout"})
	}

}

//==================================================================

func (d *demos) DemoGet(c *gin.Context) {

	id := c.Param("id")

	log.Println("DemoGet Handler Start 1")
	log.Println("DemoGet id: ", c.Param("id"))

	handleget_dto := application.HandleGet_Dto{
		Id: id,
	}

	err := d.demosApp.HandleGet(c, handleget_dto)
	if err != nil {
		log.Println("ERROR: d.demosApp.HandleGet(): ", err)
	}

	resp := new(response.DemoOk)
	resp.Message = "DemoGet OK"

	chValue, exists := c.Get("channel")

	log.Println("DemoGet Handler Start 2")

	// Check if the key exists in the context
	if !exists {
		// Handle the case where the key does not exist
		log.Println("Key 'channel' does not exist in the context")
		// Handle the error appropriately, such as returning an error response
		return
	}

	log.Println("DemoGet Handler Start 3")

	// Check if the value has the expected type
	ch, ok := chValue.(chan response.DemoOk)
	if !ok {
		// Handle the case where the value has the wrong type
		log.Println("Value associated with 'channel' has the wrong type")
		// Handle the error appropriately, such as returning an error response
		return
	}

	log.Println("DemoGet Handler Start 4")

	ch <- *resp

	//http.ResSuccess(c, resp)

}

// func (d *demos) DemoGet(c *gin.Context){

// 	// d.demosApp.HandleGet(c, c.Param("id"))
// 	wg.Add(1)

// 	go func() {
// 		defer wg.Done()
// 		id := c.Param("id")

// 		log.Println("DemoGet Handler Start")
// 		log.Println("DemoGet id: ", c.Param("id"))

// 		d.demosApp.HandleGet(c, id)

// 		resp := new(response.DemoOk)
// 		resp.Message = "DemoGet OK"
// 		http.ResSuccess(c, resp)
// 	}()

// 	log.Println("DemoGet Handler Wait")
// 	wg.Wait()

// }

func (d *demos) DemoPut(c *gin.Context) {

	log.Println("DemoPut")

	//time.Sleep(5 * time.Second)

	var item request.DemoPUTParams
	if err := http.ParseJSON(c, &item); err != nil {
		log.Println("error")
		http.ResError(c, err)
		return
	}
	log.Println("item: ", item)

	handleput_dto := application.HandlePut_Dto{
		Id: item.Property1,
	}

	err := d.demosApp.HandlePut(c, handleput_dto)
	if err != nil {
		log.Println("ERROR: d.demosApp.HandlePut", err)
	}

	resp := new(response.DemoOk)
	resp.Message = "DemoPut OK"
	http.ResSuccess(c, resp)
}
