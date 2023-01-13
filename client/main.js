window.addEventListener("load", function(evt) {
    let ws = null
    document.getElementById("form").onsubmit = function(e) {
        e.preventDefault()
        let playerId = document.querySelector('input[name="player"]:checked').value;
        // TODO: Get token after authentication (I'm gonna use token based authentication)
        let token = playerId == 1 ? "HIGHLY_SECURE_TOKEN" : "POORLY_SECURE_TOKEN"

        ws = new WebSocket("ws://localhost:3000/start?token=" + token)

        ws.onopen = function (e) {
            console.log("OPEN");
        }

        ws.onmessage = function (e) {
            // let cards = JSON.parse(e.data)
            console.log(e.data);
        }
    }

    let start = document.getElementById("start")
    start.onclick = function (e) {
        if (!ws) return

        ws.send("start")
    }
})