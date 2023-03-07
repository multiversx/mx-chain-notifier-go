package ws

// ArgsWSDispatcher -
type ArgsWSDispatcher struct {
	argsWebSocketDispatcher
}

// NewTestWSDispatcher -
func NewTestWSDispatcher(args ArgsWSDispatcher) (*websocketDispatcher, error) {
	wsArgs := argsWebSocketDispatcher{
		Hub:  args.Hub,
		Conn: args.Conn,
	}

	return newWebSocketDispatcher(wsArgs)
}

// WritePump -
func (wd *websocketDispatcher) WritePump() {
	wd.writePump()
}

// ReadPump -
func (wd *websocketDispatcher) ReadPump() {
	wd.readPump()
}

// ReadSendChannel -
func (wd *websocketDispatcher) ReadSendChannel() []byte {
	d := <-wd.send
	return d
}
