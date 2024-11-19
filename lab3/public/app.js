const peersContainer = document.getElementById("peers-container");


const socket = new WebSocket(`ws://${window.location.host}/ws`);

socket.onmessage = function(event) {
    const data = JSON.parse(event.data);
    renderPeers(data);
};

function renderPeers(peers) {
    peersContainer.innerHTML = ""; 

    peers.forEach(peer => {
        const peerDiv = document.createElement("div");
        peerDiv.className = "peer";

        const statusSpan = document.createElement("span");
        statusSpan.textContent = `Пир ${peer.port}: `;
        statusSpan.className = peer.status === "живой" ? "alive" : "dead";

        const statusText = document.createElement("span");
        statusText.textContent = peer.status;

        const killButton = document.createElement("button");
        killButton.textContent = "Убить";
        killButton.disabled = peer.status !== "живой";
        killButton.onclick = () => killPeer(peer.port);

        peerDiv.appendChild(statusSpan);
        peerDiv.appendChild(statusText);
        peerDiv.appendChild(document.createTextNode(" "));
        peerDiv.appendChild(killButton);

        peersContainer.appendChild(peerDiv);
    });
}

function killPeer(port) {
    fetch("/kill", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ port: port })
    })
    .then(response => {
        if (response.ok) {
            console.log(`Пир ${port} убит`);
        } else {
            console.error("Ошибка при убийстве пира");
        }
    })
    .catch(error => {
        console.error("Ошибка при отправке запроса:", error);
    });
}
