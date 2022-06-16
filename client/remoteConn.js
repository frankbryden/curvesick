class ServerConn {
    constructor(msgCallback) {
        this.socket = new WebSocket("ws://localhost:8090/ws");
        let self = this;
        // Connection opened
        this.socket.addEventListener('open', function (event) {
            // self.socket.send("Hey there Theo");
        });

        // Listen for messages
        this.socket.addEventListener('message', function (event) {
            console.log('Message from server ', event.data);
            msgCallback(JSON.parse(event.data));
        });
    }

    sendMessage(message) {
        this.socket.send(message);
    }
}