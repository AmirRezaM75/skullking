// Server Commands
const COMMAND_DEAL_CARDS = 'DEAL_CARDS'
const COMMAND_BETTING_STARTED = 'BETTING_STARTED'
const COMMAND_BETTING_ENDED = 'BETTING_ENDED'


// Client Commands
const COMMAND_BET = 'BET'
const COMMAND_START = 'START'

let ws = null

window.addEventListener("load", function (evt) {

    let roomId = "xxx-yyy-zzz"
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

        ws.send(makeMessage(COMMAND_START))
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

function addBet(number, container) {
    let betCard = document.createElement("div")
    betCard.classList.add("bet")
    betCard.dataset.value = number
    betCard.addEventListener('click', betClickHandler)
    let betNumber = document.createElement("span")
    betNumber.innerText = number
    betCard.appendChild(betNumber)
    container.appendChild(betCard)
}

function betClickHandler(e) {
    let message = makeMessage(COMMAND_BET, e.target.getAttribute('data-value'))
    ws.send(message)
}


function messageHandler(command, content) {
    let cardsContainer = document.getElementById("cards-container")
    let betsContainer = document.getElementById("bets-container")
    let timer = document.getElementById("timer")
    if (command === COMMAND_DEAL_CARDS) {
        content.forEach((card) => addCard(card['color'], card['number'], cardsContainer))
    }

    if (command === COMMAND_BETTING_STARTED) {
        for (let i = 0; i <= content['round']; i++) {
            addBet(i, betsContainer)
        }
        let endsAt = content['endsAt']
        let timesRemaining
        let x = setInterval(() => {
            let now = new Date().getTime() / 1000 // In seconds
            timesRemaining = Math.floor(endsAt - now)
            if (timesRemaining < 0) {
                clearInterval(x)
            } else {
                timer.innerText = timesRemaining.toString()
            }
        }, 1000)
    }
}

function makeMessage(command, content) {
    return JSON.stringify({
        command: command,
        content: content
    })
}