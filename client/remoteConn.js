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
            // console.log('Message from server ', event.data);
            msgCallback(JSON.parse(event.data));
        });
    }

    sendMessage(eventType, data) {
        let obj = {
            type: eventType,
            data: data,
        };
        this.socket.send(JSON.stringify(obj));
    }

    register(name) {
        this.sendMessage("register", {
            name: name,
        });
    }

    _sendLobbyEvent(subtype) {
        this.sendMessage("lobby", {
            subtype: subtype,
        });
    }

    sendLobbyEventReady(){
        this._sendLobbyEvent("ready");
    }
    
    sendLobbyEventUnready(){
        this._sendLobbyEvent("unready");
    }
    
    sendLobbyEventUnregister(){
        this._sendLobbyEvent("unregister");
    }

    sendKeyboardStateUpdate(keyboard_state) {
        this.sendMessage("keyboard_state", keyboard_state)
    }

    _sendGameEvent(subtype) {
        this.sendMessage("gameEvent", {
            subtype: subtype,
        });
    }

    sendState(state) {
        this.sendMessage("state", {
            state: state,
        })
    }
}