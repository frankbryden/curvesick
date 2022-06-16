var key = {};
const PLAYER_RADIUS = 8;
const HISTORY_SHIFT = 5;
const GAP_TIME = 3.5;
const GAP_LENGTH = 0.3;
const TURN_ANGLE = Math.PI*1.1;
const BASE_SPEED = 200;
const KEYS_OF_INTEREST = new Set(["ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown"])

class Point {
    constructor(x, y) {
        this.x = x;
        this.y = y;
        this.radius = 7;
    }

    dist(point){
        return Math.sqrt(Math.pow(this.x - point.x, 2) + Math.pow(this.y - point.y, 2));
    }

    isTouching(point){
        return this.dist(point) <= (this.radius + point.radius);
    }
}

class Player {
    
    constructor(x, y){
        this.x = x;
        this.y = y;
        this.speed = BASE_SPEED;
        this.angle = Math.random()*2*Math.PI;
        this.turnAngle = TURN_ANGLE;
        this.history = [];
        this.alive = true;
        this.inGap = false;
        this.inGapTimer = 0;
        this.gapTimer = GAP_TIME;
        this.suicide = false;
    }

    kill() {
        this.alive = false;
    }

    update(delta){
        if (!this.alive) {
            return;
        }
        if (key["ArrowLeft"]) {
            this.angle -= this.turnAngle*delta;
        } else if (key["ArrowRight"]) {
            this.angle += this.turnAngle*delta;
        }

        this.x += this.speed*Math.cos(this.angle)*delta;
        this.y += this.speed*Math.sin(this.angle)*delta;

        if (this.gapTimer <= 0) {
            this.inGap = true;
            this.inGapTimer = GAP_LENGTH;
            this.gapTimer = GAP_TIME;
        }

        if (!this.inGap){
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
        ctx.arc(this.x, this.y, PLAYER_RADIUS-2, 0, 2*Math.PI, false);
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

        for (let point of this.history){
            if (lastDrawn.x != -1) {
                // if gap -> DON'T STROKE
                // no gap -> DO STROKE -> BUT we don't want stroke at each point.
                // This means, we STROKE the last point before drawing the first point AFTER a gap
                if (point.dist(lastDrawn) > PLAYER_RADIUS){
                    ctx.stroke();
                    ctx.beginPath();
                    ctx.moveTo(point.x, point.y);
                }
            }
            ctx.lineTo(point.x, point.y);
            lastDrawn.x = point.x*1.0;
            lastDrawn.y = point.y*1.0;
        }
        // * * * * * * __________ * * * *
        ctx.stroke();

        ctx.restore();
    }
}

class Game {
    constructor(canvas) {
        //In the future, most of this setup will be done server side, and the client will
        //simply be given this information
        this.player = new Player(canvas.width/2, canvas.height/2);
        this.serverMessage = this.serverMessage.bind(this);
        this.serverConn = new ServerConn(this.serverMessage);
        this.canvas = canvas;
        this.lastFrame = 0;
        this.delta = 0;
        this.gameState = {};
        this.serverPlayer = {x: 0, y: 0};
        for (let y = 0; y < this.toGameStateCoords(canvas.height); y++) {
            //Init a dict of dicts, where each coord [y][x] is true if a player has been to this pos
            this.gameState[y] = {};
        }
        console.log(`Initialising game state of size ${this.toGameStateCoords(canvas.width)} x ${this.toGameStateCoords(canvas.height)}`);
        this.lastSentKeystate = {};
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


    toGameStateCoords(coord) {
        return Math.round(coord/3);
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
        this.player.update(delta/1000);

        //Get game state version of coords
        const gameStateX = this.toGameStateCoords(this.player.x);
        const gameStateY = this.toGameStateCoords(this.player.y);
        if (this.player.x+PLAYER_RADIUS >= this.canvas.width || this.player.x - PLAYER_RADIUS <= 0 || this.player.y+PLAYER_RADIUS >= this.canvas.height || this.player.y - PLAYER_RADIUS <= 0) {
            this.player.kill();
        } else if (gameStateX in this.gameState[gameStateY]) {
            console.log(`Player died because at ${this.player.x}, ${this.player.y} and game state is ${gameStateX}, ${gameStateY}`);
            console.log(this.gameState);
            this.player.kill();
        }
        if (this.player.history.length > HISTORY_SHIFT) {
            const historyLen = this.player.history.length;
            const tailPoint = this.player.history[historyLen - HISTORY_SHIFT];
            this.gameState[this.toGameStateCoords(tailPoint.y)][this.toGameStateCoords(tailPoint.x)] = true;
        }
        if (key["d"]){
            this.player.kill();
        }
        // if (JSON.stringify(this.lastSentKeystate) !== JSON.stringify(key)) {
        if (this.stateHasChanged(key)) {
            console.log(key);
            this.serverConn.sendMessage(JSON.stringify(key));
            this.lastSentKeystate = Object.assign({}, key);
        }
        this.player.render(ctx);
        ctx.beginPath();
        ctx.arc(this.serverPlayer.x, this.serverPlayer.y, 15, 0, 2*Math.PI);
        ctx.fill();
        // for (let y of Object.keys(gameState)){
        //     for (let x of Object.keys(gameState[y])){
        //         ctx.beginPath();
        //         ctx.arc(parseInt(x)*3, parseInt(y)*3, 8, 0, 2*Math.PI);
        //         ctx.stroke()
        //     }
        // }
    }

    inGame() {
        return this.player.alive;
    }
}

(function () {
    const canvas = document.getElementById('canvas');
    const ctx = canvas.getContext('2d');

    // resize the canvas to fill browser window dynamically
    window.addEventListener('resize', resizeCanvas, false);

    //Keyboard listeners
    document.addEventListener("keydown", function(e) {
        if (KEYS_OF_INTEREST.has(e.key)){
            key[e.key] = true;
        }
    });
    document.addEventListener("keyup", function(e) {
        if (KEYS_OF_INTEREST.has(e.key)){
            key[e.key] = false;
        }
    });

    function resizeCanvas() {
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
    }


    resizeCanvas();

    const game = new Game(canvas);

    let lastFrame = 0;
    function gameLoop(clk){
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
 
        if (game.inGame()){
            requestAnimationFrame(gameLoop);
        }
    }

    requestAnimationFrame(gameLoop);
})();
