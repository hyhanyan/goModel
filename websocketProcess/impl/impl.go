package impl

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	maxConnId int64
	wsConnAll = make(map[int64]*Connection)
)

type Connection struct {
	wsConnect *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte

	mutex    sync.Mutex // 对closeChan关闭上锁
	isClosed bool       // 防止closeChan被关闭多次
	Id       int64      // 记录当前所有连接的下标
}

func GetLengthConnPool() int {
	return len(wsConnAll)
}

func GetConnPool() map[int64]*Connection {
	return wsConnAll
}

func InitConnection(wsConn *websocket.Conn) (*Connection, error) {
	maxConnId++
	conn := &Connection{
		wsConnect: wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
		Id:        maxConnId,
	}
	wsConnAll[conn.Id] = conn
	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()
	return conn, nil
}

func (conn *Connection) ReadMessage() ([]byte, error) {
	var data []byte
	var err error = nil
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		data = nil
		err = errors.New("connection is closeed")
	}
	return data, err
}

func (conn *Connection) WriteMessage(data []byte) error {
	var err error = nil
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closeed")
	}

	return err
}

func (conn *Connection) Close() {
	// 线程安全，可多次调用
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
		delete(wsConnAll, conn.Id)
	}
}

func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}
		//阻塞在这里，等待inChan有空闲位置
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto ERR
		}
	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
