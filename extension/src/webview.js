// eslint-disable-next-line no-undef
const vscode = acquireVsCodeApi();

// Wait for DOM, before attempting anything element-wise
document.addEventListener('DOMContentLoaded', function () {

    const sendButton = document.getElementById('sendBtn');
    const msgInput = document.getElementById('msgInput');

    sendButton.addEventListener('click', function () {
        const userMsg = msgInput.value;
        msgInput.value = '';

        // Display the grasped message
        displayMessage(userMsg, 'user')

        vscode.postMessage({
            command: 'execNewMsg',
            message:  userMsg, 
        });
    });
});

// Listen for new messages, act upon them
window.addEventListener('message', e => {
    const msg = e.data;
    switch (msg.command) {            
        case 'returnNewMsg':
            displayResponse(msg.data.message, 'bot');
            break;
    }
});

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
