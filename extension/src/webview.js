// eslint-disable-next-line no-undef
const vscode = acquireVsCodeApi();


// Start to debunk a possibility:
if ("WebTransport" in window) {
    console.log("WebTransport support: ✅");
} else {
    console.error("WebTransport support: ❌");
}


// eslint-disable-next-line no-unused-vars
function switchToConvoView() {
    document.getElementById('homeView').style.display = 'none';
    document.getElementById('convoView').style.display = 'flex';
}

// eslint-disable-next-line no-unused-vars
function switchToHomeView() {
    document.getElementById('convoView').style.display = 'none';
    document.getElementById('homeView').style.display = 'flex';
}

document.addEventListener('DOMContentLoaded', function () {

    const sendButton = document.getElementById('sendBtn');
    const msgInput = document.getElementById('msgInput');

    sendButton.addEventListener('click', function () {
        // Grasp the user's message
        const userMsg = msgInput.value;

        // Clear input field
        msgInput.value = '';

        // Display the grasped message
        displayMessage(userMsg, 'user')

        vscode.postMessage({
            command: 'execGeminiMsg',
            message: userMsg
        });
    });
});

window.addEventListener('message', e => {
    const msg = e.data;
    switch (msg.command) {
        case 'displayGeminiRes':
            displayResponse(msg.data.message, 'bot');
            break;
    }
});

function displayMessage(text, sender) {
    const convoArea = document.getElementById('convoArea');
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.classList.add(sender);
    messageElement.textContent = text;
    convoArea.appendChild(messageElement);
    convoArea.scrollTop = convoArea.scrollHeight;
}

function displayResponse(text, sender) {
    const convoArea = document.getElementById('convoArea');
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.classList.add(sender);
    messageElement.textContent = text;
    convoArea.appendChild(messageElement);
    convoArea.scrollTop = convoArea.scrollHeight;
}
