export function newWebsocketConn(roomId: string) {
	return new WebSocket('ws://localhost:1437/ws/' + roomId);
}
