// eslint-disable-next-line no-undef
const vscode = acquireVsCodeApi();

// ℹ️: Flags for res
let resInProg = false;
let resChunkQueue = [];
let test = '';

// ℹ️: Communicate early
vscode.postMessage({
    command: 'execConvoList',
});

// ℹ️: Wait for DOM, before attempting anything element-wise
document.addEventListener('DOMContentLoaded', function () {
    const sendMsgBtn = document.getElementById('sendBtn');
    const msgInput = document.getElementById('msgInput');

    sendMsgBtn.addEventListener('click', function () {
        const msg = msgInput.value;
        msgInput.value = '';

        renderUsrMsg(msg);

        vscode.postMessage({
            command: 'execNewMsg',
            message: msg,
        });
    });
});

// ℹ️: Listen for new messages, act upon
window.addEventListener('message', e => {
    const msg = e.data;

    switch (msg.command) {
        case 'returnConvoList':
            var arr = msg.data.conversations;
            listRecentConvos(arr);
            break;
            
        case 'returnConvoView':
            reopenConvoView();
            break;

        case 'returnMsg':
            resChunkQueue.push(msg.data);
            if (!resInProg) {
                resInProg = true;
                streamGeminiRes();
            }
            break;
    }
});

function streamGeminiRes() {
    
    if (resChunkQueue.length === 0) {
        resInProg = false;
        return;
    };

    const area = document.getElementById('convoArea');
    let resEl = document.querySelector('.message.bot.streaming');
    if (!resEl) {
        resEl = document.createElement('div');
        resEl.classList.add('message', 'bot', 'streaming');
        area.appendChild(resEl);
    };

    const chunk = resChunkQueue.shift();

    // Check for EOF and reset flags
    if (chunk === 'EOF') {
        resEl.classList.remove('streaming');
        resInProg = false;

        const htmlContent = marked.parse(test);
        resEl.innerHTML = htmlContent;
        // const htmlContent = processMD(test)
        // resEl.innerHTML = htmlContent
        test = ''; 
        return;
    } else {
        test += chunk;  
    }
    
    let messageContent;
    messageContent = chunk; 
    
    let currentIndex = 0;
    const animationInterval = 10;

    function streamResponse() {
        if (currentIndex < messageContent.length) {
            resEl.textContent += messageContent[currentIndex];
            currentIndex++;
            setTimeout(streamResponse, animationInterval);
        } else {
            setTimeout(streamGeminiRes, animationInterval);
        };
    };

    streamResponse();

    area.scrollTop = area.scrollHeight;
}

//

function renderUsrMsg(text) {
    const convoArea = document.getElementById('convoArea');
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.classList.add('user');
    messageElement.textContent = text;
    convoArea.appendChild(messageElement);
    convoArea.scrollTop = convoArea.scrollHeight;
};

//

function listRecentConvos(list) {
    const convoListEl = document.getElementById('convo-list');

    if (list.length === 0) {
        // ℹ️: Problematic alongside of the `show more...` el
        convoListEl.innerHTML = '<li>Create your first conversation?</li>';
        return;
    }

    const displayCount = 4;
    const displayConversations = list.slice(0, displayCount);

    displayConversations.forEach((convo, index) => {
        const listItem = document.createElement('li');
        listItem.textContent = convo.title;
        listItem.id = `convo-item-${index}`;
        convoListEl.appendChild(listItem);

        listItem.addEventListener('click', () => {
            vscode.postMessage({
                command: 'execConvoView',
                data: convo.id,
            });
        });
    });
};

//

function newProjectAction() {
    // ℹ️: Will grow complex 

    vscode.postMessage({
        command: 'execNewConvo',
    });
    document.getElementById('homeView').style.display = 'none';
    // ℹ️: Hide all views, but the specific one needed rendered
    document.getElementById('convoView').style.display = 'flex';
};

function returnHome() {
    vscode.postMessage({
        command: 'execReturnHome',
    });
    document.getElementById('convoView').style.display = 'none';
    // ℹ️: Hide all views, but the specific one needed rendered
    document.getElementById('homeView').style.display = 'flex';
};

function reopenConvoView() {
    document.getElementById('homeView').style.display = 'none';
    // ℹ️: Hide all views, but the specific one needed rendered
    document.getElementById('convoView').style.display = 'flex';
};
