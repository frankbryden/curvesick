var key = {};
const PLAYER_RADIUS = 8;
const HISTORY_SHIFT = 5;
const GAP_TIME = 3.5;
const GAP_LENGTH = 0.3;
const TURN_ANGLE = Math.PI * 1.1;
const BASE_SPEED = 200;
const KEYS_OF_INTEREST = new Set(["ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown"])

class Point {
    constructor(x, y) {
        this.x = x;
        this.y = y;
        this.radius = 7;
    }

    dist(point) {
        return Math.sqrt(Math.pow(this.x - point.x, 2) + Math.pow(this.y - point.y, 2));
    }

    isTouching(point) {
        return this.dist(point) <= (this.radius + point.radius);
    }
}

class Player {

    constructor(x, y, speed, angle, turnAngle) {
        this.x = x;
        this.y = y;
        this.speed = speed;
        this.angle = angle;
        this.turnAngle = turnAngle;
        this.history = [];
        this.alive = true;
        this.inGap = false;
        this.inGapTimer = 0;
        this.gapTimer = GAP_TIME;
        this.suicide = false;
    }
    // constructor(x, y) {
    //     this.x = x;
    //     this.y = y;
    //     this.speed = BASE_SPEED;
    //     this.angle = Math.random() * 2 * Math.PI;
    //     this.turnAngle = TURN_ANGLE;
    //     this.history = [];
    //     this.alive = true;
    //     this.inGap = false;
    //     this.inGapTimer = 0;
    //     this.gapTimer = GAP_TIME;
    //     this.suicide = false;
    // }

    kill() {
        this.alive = false;
    }

    update(delta) {
        if (!this.alive) {
            return;
        }
        if (key["ArrowLeft"]) {
            this.angle -= this.turnAngle * delta;
        } else if (key["ArrowRight"]) {
            this.angle += this.turnAngle * delta;
        }

        this.x += this.speed * Math.cos(this.angle) * delta;
        this.y += this.speed * Math.sin(this.angle) * delta;

        if (this.gapTimer <= 0) {
            this.inGap = true;
            this.inGapTimer = GAP_LENGTH;
            this.gapTimer = GAP_TIME;
        }

        if (!this.inGap) {
            this.gapTimer -= delta;
            this.history.push(new Point(this.x, this.y));
            if (this.suicide) {
                this.kill();
            }
        } else {
            this.inGapTimer -= delta;
            if (this.inGapTimer <= 0) {
                this.inGap = false;
                // this.suicide = true;
            }
        }
    }

    render(ctx) {
        //draw head
        ctx.beginPath();
        ctx.arc(this.x, this.y, PLAYER_RADIUS - 2, 0, 2 * Math.PI, false);
        ctx.moveTo(this.x, this.y);
        ctx.fill();

        ctx.save();

        //draw head
        ctx.lineWidth = PLAYER_RADIUS;
        if (this.alive) {
            ctx.strokeStyle = 'purple';
        } else {
            console.log("Printing in grey");
            ctx.strokeStyle = 'grey';
        }
        ctx.beginPath();
        let lastDrawn = new Point(-1, -1);
        // const toDraw = this.inGap ? 
        //     this.history :
        //     this.history.slice(0, this.history.length - HISTORY_SHIFT);

        for (let point of this.history) {
            if (lastDrawn.x != -1) {
                // if gap -> DON'T STROKE
                // no gap -> DO STROKE -> BUT we don't want stroke at each point.
                // This means, we STROKE the last point before drawing the first point AFTER a gap
                if (point.dist(lastDrawn) > PLAYER_RADIUS) {
                    ctx.stroke();
                    ctx.beginPath();
                    ctx.moveTo(point.x, point.y);
                }
            }
            ctx.lineTo(point.x, point.y);
            lastDrawn.x = point.x * 1.0;
            lastDrawn.y = point.y * 1.0;
        }
        // * * * * * * __________ * * * *
        ctx.stroke();

        ctx.restore();
    }
}

class Game {
    constructor(canvas, serverConn) {
        //In the future, most of this setup will be done server side, and the client will
        //simply be given this information
        this.player = new Player(canvas.width / 2, canvas.height / 2);
        this.serverMessage = this.serverMessage.bind(this);
        this.serverConn = serverConn;
        this.canvas = canvas;
        this.lastFrame = 0;
        this.delta = 0;
        this.gameState = {};
        this.ready = false;
        this.serverDims = {
            width: -1,
            height: -1,
        };
        this.serverPlayers = {};
        this.lastSentKeystate = {};
        this.running = false;
    }

    initGameState() {
        for (let y = 0; y < this.toGameStateCoords(this.serverDims.height); y++) {
            //Init a dict of dicts, where each coord [y][x] is true if a player has been to this pos
            this.gameState[y] = {};
        }
        // console.log(`Initialising game state of size ${this.toGameStateCoords(canvas.width)} x ${this.toGameStateCoords(canvas.height)}`);
    }

    serverMessage(data) {
        console.log(typeof data);
        if ("X" in data && "Y" in data) {
            this.serverPlayer = {
                x: data.X,
                y: data.Y,
            }
        }
    }

    buildWithServerDefinitions(state_def) {
        this.serverDims = {
            width: state_def.width,
            height: state_def.height,
        };
        for (let player of state_def.players) {
            this.serverPlayers[player.name] = new Player(player.x, player.y, player.speed, player.angle, player.turnAngle);
        }
        this.initGameState();
    }

    serverUpdate(state_def){
        for (let player of state_def.players) {
            this.serverPlayers[player.name].x = player.x_float;
            this.serverPlayers[player.name].y = player.y_float;
        }
    }

    toGameStateCoords(coord) {
        return {
            x: this.serverDims.width/this.canvas.width * coord.x,
            y: this.serverDims.height/this.canvas.height * coord.y,
        }
    }

    toDisplayCoords(coord) {
        return {
            x: this.canvas.width/this.serverDims.width * coord.x,
            y: this.canvas.height/this.serverDims.height * coord.y,
        }
    }

    keyStateToBinaryRepr(keyState) {
        //TODO implement later - 2 bits for key; one bit for state
    }

    stateHasChanged(newState) {
        for (let key of Object.keys(newState)) {
            if (!(key in this.lastSentKeystate) || this.lastSentKeystate[key] != newState[key]) {
                return true;
            }
        }
        return false;
    }

    gameLoop(ctx, delta) {
        this.player.update(delta / 1000);

        //Get game state version of coords
        const gameStateX = this.toGameStateCoords(this.player.x);
        const gameStateY = this.toGameStateCoords(this.player.y);
        // console.log(`${this.player.x} x became ${gameStateX}`);
        // if (this.player.x + PLAYER_RADIUS >= this.canvas.width || this.player.x - PLAYER_RADIUS <= 0 || this.player.y + PLAYER_RADIUS >= this.canvas.height || this.player.y - PLAYER_RADIUS <= 0) {
        //     this.player.kill();
        // } else if (gameStateX in this.gameState[gameStateY]) {
        //     console.log(`Player died because at ${this.player.x}, ${this.player.y} and game state is ${gameStateX}, ${gameStateY}`);
        //     console.log(this.gameState);
        //     this.player.kill();
        // }
        // if (this.player.history.length > HISTORY_SHIFT) {
        //     const historyLen = this.player.history.length;
        //     const tailPoint = this.player.history[historyLen - HISTORY_SHIFT];
        //     this.gameState[this.toGameStateCoords(tailPoint.y)][this.toGameStateCoords(tailPoint.x)] = true;
        // }
        // if (key["d"]) {
        //     this.player.kill();
        // }
        // if (JSON.stringify(this.lastSentKeystate) !== JSON.stringify(key)) {
        if (this.stateHasChanged(key)) {
            this.serverConn.sendKeyboardStateUpdate(key);
            this.lastSentKeystate = Object.assign({}, key);
        }
        // this.player.render(ctx);
        // ctx.beginPath();
        // ctx.arc(this.serverPlayer.x, this.serverPlayer.y, 15, 0, 2 * Math.PI);
        // ctx.fill();
        for (const [playerName, player] of Object.entries(this.serverPlayers)) {
            const translated = this.toDisplayCoords(player);
            ctx.beginPath();
            ctx.arc(translated.x, translated.y, 15, 0, 2 * Math.PI);
            ctx.fill();
            // console.log(`Filling arc at ${translated.x}, ${translated.y}`);
        }
        // for (let y of Object.keys(gameState)){
        //     for (let x of Object.keys(gameState[y])){
        //         ctx.beginPath();
        //         ctx.arc(parseInt(x)*3, parseInt(y)*3, 8, 0, 2*Math.PI);
        //         ctx.stroke()
        //     }
        // }
    }

    inGame() {
        return this.running && this.player.alive;
    }
}

(function () {
    const form = document.getElementById('form');
    const readyBtn = document.getElementById('readyBtn');
    const lobbyUl = document.getElementById('lobby');
    const nameInput = document.getElementById('name');
    const canvas = document.getElementById('canvas');
    const ctx = canvas.getContext('2d');

    canvas.style.display = "none";
    readyBtn.style.display = "none";

    const serverConn = new ServerConn(serverMessage)

    const game = new Game(canvas, serverConn);

    // resize the canvas to fill browser window dynamically
    window.addEventListener('resize', resizeCanvas, false);

    //Keyboard listeners
    document.addEventListener("keydown", function (e) {
        if (KEYS_OF_INTEREST.has(e.key)) {
            key[e.key] = true;
        }
    });
    document.addEventListener("keyup", function (e) {
        if (KEYS_OF_INTEREST.has(e.key)) {
            key[e.key] = false;
        }
    });

    function hideReadyBtn() {
        readyBtn.style.display = "none";
    }


    function showReadyBtn() {
        readyBtn.style.display = "block";
    }

    function hideForm() {
        form.style.display = "none";
    }

    function showForm() {
        form.style.display = "block";
    }

    function showCanvas() {
        canvas.style.display = "block";
    }

    function startGame() {
        hideForm();
        hideReadyBtn();
        showCanvas();
        game.running = true;
        requestAnimationFrame(gameLoop);
    }

    function toggleReady() {
        game.ready = !game.ready;
        if (game.ready) {
            serverConn.sendLobbyEventReady();
            readyBtn.innerHTML = "Unready";
        } else {
            serverConn.sendLobbyEventUnready();
            readyBtn.innerHTML = "Ready";
        }
    }
    function serverMessage(msg) {
        const payload = msg["data"];
        switch (msg["type"]) {
            case "lobby":
                displayLobby(payload["players"]);
                break;
            case "init_state_def":
                console.log("Got server defs!");
                console.log(payload);
                game.buildWithServerDefinitions(payload);
                break;
            case "game_event":
                switch (payload["sub_type"]) {
                    case "round_start":
                        game.running = true;
                        startGame();
                        break;
                    case "round_end":
                        break;
                }
                break;
            case "reg_state_def":
                game.serverUpdate(payload);
                break;
            default:
                console.log(`${msg["type"]} not found in switch!!`);
        }
    }

    function displayLobby(lobby) {
        lobby.sort(function (a, b) {
            const nameA = a.name.toUpperCase(); // ignore upper and lowercase
            const nameB = b.name.toUpperCase(); // ignore upper and lowercase
            if (nameA < nameB) {
                return -1;
            }
            if (nameA > nameB) {
                return 1;
            }

            // names must be equal
            return 0;
        });
        lobbyUl.innerHTML = "";
        for (const player of lobby) {
            const li = document.createElement("li");
            li.appendChild(document.createTextNode(`${player.name}: Ready ? ${player.isReady}`));
            lobbyUl.appendChild(li);
        }
    }

    function resizeCanvas() {
        canvas.width = 500;//window.innerWidth - 100;
        canvas.height = 500;//window.innerHeight - 100;
    }

    function submitRegistration(e) {
        e.preventDefault();
        console.log("Login from: " + nameInput.value);
        serverConn.register(nameInput.value);
        hideForm();
        showReadyBtn();
    }

    form.addEventListener('submit', submitRegistration);
    readyBtn.addEventListener('click', toggleReady);


    resizeCanvas();

    let lastFrame = 0;
    function gameLoop(clk) {
        delta = clk - lastFrame;
        lastFrame = clk;
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.beginPath();
        ctx.rect(0, 0, canvas.width, canvas.height);
        ctx.stroke();

        game.gameLoop(ctx, delta);

        ctx.strokeText("Delta: " + Math.round(delta), 10, 50);
        ctx.strokeText("In gap: " + game.player.inGap, 10, 100);
        ctx.strokeText("Key      : " + JSON.stringify(key), 10, 150);
        ctx.strokeText("Last sent: " + JSON.stringify(game.lastSentKeystate), 10, 200);

        if (game.inGame()) {
            requestAnimationFrame(gameLoop);
        }
    }
})();
