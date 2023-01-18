// Commands
const COMMAND_DEAL_CARDS = 'DEAL_CARDS'


window.addEventListener("load", function (evt) {

    let roomId = "xxx-yyy-zzz"
    let ws = null
    document.getElementById("form").onsubmit = function (e) {
        e.preventDefault()
        let playerId = document.querySelector('input[name="player"]:checked').value;
        // TODO: Get token after authentication (I'm gonna use token based authentication)

        ws = new WebSocket("ws://localhost:3000/ws?roomId=" + roomId + "&userId=" + playerId)
        ws.onopen = function (e) {
            console.log("OPEN");
        }

        ws.onmessage = function (e) {
            let message = JSON.parse(e.data)
            let content = message['content']
            let command = message['command']
            if (message['contentType'] === 'json') {
                content = JSON.parse(message['content'])
            }

            console.log(message)

            messageHandler(command, content);
        }
    }

    document.getElementById("create-room").onclick = function (e) {
        console.log("create-room")
        e.preventDefault()
        fetch('http://localhost:3000/rooms', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({"id": roomId, "name": "market"})
        })
            .then(response => response.json())
            .then(response => console.log(JSON.stringify(response)))
    }

    let start = document.getElementById("start")
    start.onclick = function (e) {
        if (!ws) return

        ws.send("start")
    }

    let bet = document.getElementById("bet")
    bet.onclick = function (e) {
        console.log("Bet")
        if (!ws) return

        ws.send("bet")
    }
})

function addCard(color, number, container) {
    let card = document.createElement("div")
    card.style.backgroundColor = color
    card.classList.add("card")
    let cardNumber = document.createElement("span")
    cardNumber.innerText = number
    card.appendChild(cardNumber)
    container.appendChild(card)
}


function messageHandler(command, content) {
    let cardsContainer = document.getElementById("cards-container")

    if (command === COMMAND_DEAL_CARDS) {
        content.forEach((card) => addCard(card['color'], card['number'], cardsContainer))
    }
}