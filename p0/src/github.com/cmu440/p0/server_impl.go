// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0
import (
	"net"
	// "fmt"
	"strconv"
	"bufio"
)
type client_t struct {
	quitCh chan bool
	dataOutCh chan []byte
	fd net.Conn
}

type multiEchoServer struct {
	maxClientNum int 
	ConnectNum int
	Maxfd int
	
	quitCh chan bool
	newClientCh chan net.Conn
	dataToSendCh chan []byte
	clientToDelCh chan int
	countResquestCH chan bool
	countResponseCh chan int

	Clients [100]client_t
	Listenfd net.Listener
	e error
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	var v multiEchoServer

	v.ConnectNum=0
	v.Maxfd=0
	v.maxClientNum=100

	v.quitCh=make(chan bool)
	v.newClientCh=make(chan net.Conn)
	v.dataToSendCh=make(chan []byte,100)
	v.clientToDelCh=make(chan int)
	v.countResquestCH = make(chan bool)
	v.countResponseCh = make(chan int)

	return &v
}

func (mes *multiEchoServer) Start(port int) error {
	mes.Listenfd,mes.e=net.Listen("tcp",":"+strconv.Itoa(port))
	if mes.e!=nil {
		// fmt.Println("Error on listen:",mes.e)
		return mes.e
	}
	go runServer(mes)
	return nil
}

func (mes *multiEchoServer) Close() {
	mes.quitCh<-true
}

func (mes *multiEchoServer) Count() int {
	mes.countResquestCH<-true
	return <-mes.countResponseCh
}

func runServer(mes *multiEchoServer) {
	go acceptNewClient(mes)
	for{
		select{
			case <-mes.quitCh:
				mes.Listenfd.Close()
				for i:=0;i<mes.Maxfd;i++ {
					if mes.Clients[i].fd!=nil{
						mes.Clients[i].quitCh<-true
					}
				}
				return

			case num:=<-mes.clientToDelCh:
				cleanClient(mes,num)

			case conn:=<-mes.newClientCh:
				num:=addClient(mes,conn)
				if num>=0 {
					go handleConn(mes,mes.Clients[num],num)
				}

			case buf:=<-mes.dataToSendCh:
				sendMessageToAll(mes,buf)

			case <-mes.countResquestCH:
				mes.countResponseCh<-mes.ConnectNum
		}
	}
}
func acceptNewClient(mes *multiEchoServer) {
	for{
		// fmt.Println("waiting for connection")	
		conn,e:=mes.Listenfd.Accept()
		if e!=nil {
			// fmt.Println("Error on accept:",e)
			return
		}
		mes.newClientCh<-conn		
	}
}

func addClient(mes *multiEchoServer,conn net.Conn) int{
	var i int
	for i=0;i<mes.maxClientNum;i++{
		if mes.Clients[i].fd==nil{
			mes.Clients[i].fd=conn
			mes.Clients[i].quitCh=make(chan bool)
			mes.Clients[i].dataOutCh=make(chan []byte,100)
			break
		}
	}
	if i==mes.maxClientNum{
		conn.Write([]byte("we cannot handle anymore!"))
		conn.Close()
		return -1
	}
	if i>mes.Maxfd {
		mes.Maxfd=i
	}
	mes.ConnectNum++
	return i
}
func handleConn(mes *multiEchoServer,c client_t,num int){
	go readFromConn(mes,c.fd,num)
	for{
		select{
		case <-c.quitCh:
			c.fd.Close()
			return
		case buf:=<-c.dataOutCh:
			c.fd.Write(buf)
		}
	}
}
func readFromConn(mes *multiEchoServer,conn net.Conn,num int){
	reader:=bufio.NewReader(conn)
	for{		
		msg,err := reader.ReadBytes('\n')
		if err != nil {
	        // fmt.Println("Error on read: ", err,"--closing")
	        conn.Close()
	        mes.clientToDelCh<-num
	        return 
	    }
	    // fmt.Println("Client ",num," sent:", string(msg))
	    mes.dataToSendCh<-msg
	}
}

func cleanClient(mes *multiEchoServer,num int){
	if mes.Clients[num].fd!=nil{
		mes.Clients[num].fd=nil
		mes.Clients[num].quitCh=nil
		mes.Clients[num].dataOutCh=nil
		mes.ConnectNum--

	}
}
func sendMessageToAll(mes *multiEchoServer, message []byte){
	for i:=0;i<=mes.Maxfd;i++{
		if mes.Clients[i].fd!=nil{
			if(len(mes.Clients[i].dataOutCh)<cap(mes.Clients[i].dataOutCh)){
				mes.Clients[i].dataOutCh<-message
			}
		}
	}
}





