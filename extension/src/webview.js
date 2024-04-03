// eslint-disable-next-line no-undef
const vscode = acquireVsCodeApi();

let msgInProg = false;
let chunkQueue = [];

// Wait for DOM, before attempting anything element-wise
document.addEventListener('DOMContentLoaded', function () {
    const sendButton = document.getElementById('sendBtn');
    const msgInput = document.getElementById('msgInput');

    sendButton.addEventListener('click', function () {
        const userMsg = msgInput.value;
        msgInput.value = '';

        // Display the grasped message
        displayMessage(userMsg, 'user');

        vscode.postMessage({
            command: 'execNewMsg',
            message: userMsg,
        });
    });
});

// Listen for new messages, act upon them
window.addEventListener('message', e => {
    const msg = e.data;

    switch (msg.command) {
        case 'updateDisplay':
            chunkQueue.push(msg.data);

            if (!msgInProg) {
                msgInProg = true;
                processChunkQueue();
            }
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

function processChunkQueue() {
    if (chunkQueue.length === 0) {
        msgInProg = false;
        return;
    }

    const area = document.getElementById('convoArea');
    let msgEl = document.querySelector('.message.bot.streaming');
    if (!msgEl) {
        msgEl = document.createElement('div');
        msgEl.classList.add('message', 'bot', 'streaming');
        area.appendChild(msgEl);
    }

    const chunk = chunkQueue.shift();

    // Check for EOF and reset flags
    if (chunk === 'EOF') {
        msgEl.classList.remove('streaming');
        storedResponse = [];
        msgInProg = false;
        return;
    }

    let messageContent;
    try {
        const parsedChunk = JSON.parse(chunk);
        messageContent = parsedChunk.message;
    } catch (error) {
        messageContent = chunk;
    }

    let currentIndex = 0;
    const animationInterval = 10;

    function updateMessage() {
        if (currentIndex < messageContent.length) {
            msgEl.textContent += messageContent[currentIndex];
            currentIndex++;
            setTimeout(updateMessage, animationInterval);
        } else {
            setTimeout(processChunkQueue, animationInterval);
        }
    }

    updateMessage();

    area.scrollTop = area.scrollHeight;
}

// eslint-disable-next-line no-unused-vars
function switchToConvoView() {
    vscode.postMessage({
        command: 'execNewConvo',
        message: ''
    });

    document.getElementById('homeView').style.display = 'none';
    document.getElementById('convoView').style.display = 'flex';
}

// eslint-disable-next-line no-unused-vars
function switchToHomeView() {
    vscode.postMessage({
        command: 'execReturnHome',
        message: ''
    });

    document.getElementById('convoView').style.display = 'none';
    document.getElementById('homeView').style.display = 'flex';
}

