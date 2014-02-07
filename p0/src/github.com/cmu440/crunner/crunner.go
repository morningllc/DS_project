package main

import (
	"fmt"
	"net"
	"strconv"
	// "bufios"
)

// const (
// 	defaultHost = "localhost"
// 	defaultPort = 9999
// )

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's echoed response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	fmt.Println("Not implemented.")
	// var ttt [1024]byte
	// ttt:=[]byte("aaa\nbbb\nccc\n\nddd\n")
	// var e error
	// numm:=0
	// var num int
	// var xzz [1024]byte
	// x:=xzz[:]
	// for {
	// 	num,x,e=bufio.ScanLines(ttt[numm:],false)
	// 	numm=numm+num
	// 	if e!=nil{
	// 		fmt.Println("Error slice: ", e)
	// 		break
	// 	}

	// 	fmt.Println(numm,string(x))
	// 	if numm==len(ttt) {break}
	// }
	// // fmt.Println(bufio.ScanLines([]byte("aaa\nbbb\nccc\n\nddd\n"),true))


	// fmt.Println("Not implemented.\n\n\n")
	var buf [1024]byte
	hostport:=net.JoinHostPort("localhost", strconv.Itoa(11223))
	if addr, err := net.ResolveTCPAddr("tcp", hostport); err != nil {
			return
		} else if conn, err := net.DialTCP("tcp", nil, addr); err != nil {
			return
		} else {

			for j:=0;j<100;j++{
				writeMsg:=fmt.Sprintf("%d\n",j)
				conn.Write([]byte(writeMsg))
			}
			for i:=0;;i++{
				n,err := conn.Read(buf[:])
				if err != nil {
					fmt.Println("Error on read: ", err,"--closing")
					conn.Close()
					return
				}
			
				fmt.Println("receive",i," : ", string(buf[0:n]))
			}
		}
}
