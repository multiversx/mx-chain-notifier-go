package ws

type ArgsWSDispatcher struct {
	argsWebSocketDispatcher
}

func NewTestWSDispatcher(args ArgsWSDispatcher) (*websocketDispatcher, error) {
	wsArgs := argsWebSocketDispatcher{
		Hub:  args.Hub,
		Conn: args.Conn,
	}

	return newWebSocketDispatcher(wsArgs)
}

func (wd *websocketDispatcher) WritePump() {
	wd.writePump()
}

func (wd *websocketDispatcher) ReadPump() {
	wd.readPump()
}
