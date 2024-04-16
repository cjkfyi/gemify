// eslint-disable-next-line no-undef
const vscode = acquireVsCodeApi();

// ℹ️: Flags for res
let resInProg = false;
let resChunkQueue = [];
let test = '';

// ℹ️: Communicate early
vscode.postMessage({
    command: 'execProjList',
});

// ℹ️: Listen for new messages, act upon
window.addEventListener('message', e => {
    const msg = e.data;

    switch (msg.command) {

        case 'returnProjList':
            var arr = msg.data.projects;
            listRecentProjs(arr);
            break;

        case 'returnChatList':
            var arr = msg.data.chats;
            var proj = msg.data.proj;
            listRecentChats(arr, proj);
            break;

        case 'returnMsgList':
            var arr = msg.data.msgs;
            var chat = msg.data.chat;
            listRecentMsgs(arr, chat);
            break;

        case 'returnMsg':
            resChunkQueue.push(msg.data);
            if (!resInProg) {
                resInProg = true;
                streamGeminiRes();
            }
            break;
    };
});

function listRecentProjs(list) {
    const projListEl = document.getElementById('proj-list');

    // Breaks if doesn't exist

    const displayCount = 4;
    const displayConversations = list.slice(0, displayCount);

    displayConversations.forEach((proj, index) => {

        // Create Tile Elements
        const listItem = document.createElement('li');
        const tileLink = document.createElement('a');
        const tileName = document.createElement('h3');
        const tileDesc = document.createElement('p');

        // Add Styling Classes
        listItem.classList.add('proj-tile');
        tileLink.classList.add('proj-tile-link');
        tileName.classList.add('proj-tile-name');
        tileDesc.classList.add('proj-tile-desc');

        // Populate Tile Content
        tileName.textContent = proj.name;
        tileDesc.textContent = proj.desc;

        // Assemble Tile 
        tileLink.href = '#';
        tileLink.appendChild(tileName);
        tileLink.appendChild(tileDesc);
        listItem.appendChild(tileLink);
        projListEl.appendChild(listItem);


        // cleanup after click so they don't linger?
        listItem.addEventListener('click', () => {

            selectProj();
            
            vscode.postMessage({
                command: 'execChatList',
                data: {
                    projID: proj.projID,
                    proj: proj,
                },
            });
        });
    });
};

function listRecentChats(list, proj) {
    const chatListEl = document.getElementById('chat-list-area');
    //
    const projectTitleEl = document.getElementById('proj-title');
    projectTitleEl.textContent = proj.name; 
    //
    const backBtnEl = document.getElementById('back-button');
    backBtnEl.addEventListener('click', () => {
        chatListEl.innerHTML = ''; 
        returnHome()
    });

    const displayCount = 4; 
    const displayChats = list.slice(0, displayCount);

    displayChats.forEach((chat, index) => {

        // Create Tile Elements
        const listItem = document.createElement('li');
        const tileLink = document.createElement('a');
        const tileName = document.createElement('h3');
        const tileDesc = document.createElement('p');

        // Add Styling Classes
        listItem.classList.add('chat-tile');
        tileLink.classList.add('chat-tile-link');
        tileName.classList.add('chat-tile-name');
        tileDesc.classList.add('chat-tile-desc');

        // Populate Tile Content
        tileName.textContent = chat.name;
        tileDesc.textContent = chat.desc; 

        // Assemble Tile 
        tileLink.href = '#';
        tileLink.appendChild(tileName);
        tileLink.appendChild(tileDesc);
        listItem.appendChild(tileLink);
        chatListEl.appendChild(listItem);


        listItem.addEventListener('click', () => {

            selectChat();
            
            vscode.postMessage({
                command: 'execMsgList',
                data: {
                    chat: chat,
                    proj: proj,
                },
            });
        });
    });
};

function listRecentMsgs(list, chat) {

    const convoTitleEl = document.getElementById('convo-title');
    convoTitleEl.textContent = chat.name;

    const sendMsgBtn = document.getElementById('sendBtn');
    const msgInput = document.getElementById('msg-input');

    sendMsgBtn.addEventListener('click', function () {

        const msg = msgInput.value;
        renderUsrMsg(msg);
        
        msgInput.value = ''; // reset

        vscode.postMessage({
            command: 'execNewMsg',
            data: {
                chat: chat,
                message: msg,
            },
        });
    });

    //

    const convoArea = document.getElementById('convo-area');
    convoArea.innerHTML = ''; // reset
    list.forEach((msg) => {
        const messageElement = document.createElement('div');
        messageElement.classList.add('message');

        if (msg.isUser) {
            messageElement.classList.add('user');
        } else {
            messageElement.classList.add('bot'); 
        }

        messageElement.textContent = msg.message;
        convoArea.appendChild(messageElement);    
    });
    convoArea.scrollTop = convoArea.scrollHeight; 
}


//

function newProjectView() {
    document.getElementById('home-view').style.display = 'none';
    document.getElementById('new-proj-view').style.display = 'flex';
};

function newChatView() {
    // vscode.postMessage({
    //     command: 'execNewProj',
    // });
    document.getElementById('chat-view').style.display = 'none';
    document.getElementById('new-chat-view').style.display = 'flex';
}

//

function createProjectBtn() {

    // todo: take inputted data
    //       create a new project 
    //       route the user to it

    // shoot the req w/ passed data, shape?

    // vscode.postMessage({
    //     command: 'execNewProj',
    // });
    document.getElementById('new-proj-view').style.display = 'none';
    document.getElementById('chat-view').style.display = 'flex';
};


function createChatBtn() {

    // todo: take inputted data
    //       create a new chat  
    //       route the user to convo

    // shoot the req w/ passed data, shape?

    // vscode.postMessage({
    //     command: 'execNewProj',
    // });
    document.getElementById('new-chat-view').style.display = 'none';
    document.getElementById('convo-view').style.display = 'flex';
};

//

//

//

function returnHome() {
    // vscode.postMessage({
    //     command: 'execReturnHome',
    // });
    document.getElementById('settings-view').style.display = 'none';
    document.getElementById('new-chat-view').style.display = 'none';
    document.getElementById('chat-view').style.display = 'none';
    document.getElementById('convo-view').style.display = 'none';
    document.getElementById('new-proj-view').style.display = 'none';
    document.getElementById('home-view').style.display = 'flex';
};

function returnChat() {
    // vscode.postMessage({
    //     command: 'execReturnHome',
    // });
    document.getElementById('settings-view').style.display = 'none';
    document.getElementById('new-chat-view').style.display = 'none';
    document.getElementById('home-view').style.display = 'none';
    document.getElementById('convo-view').style.display = 'none';
    document.getElementById('new-proj-view').style.display = 'none';
    document.getElementById('chat-view').style.display = 'flex';
};

//

function selectProj() {
    document.getElementById('home-view').style.display = 'none';
    document.getElementById('chat-view').style.display = 'flex';
};

function selectChat() {
    document.getElementById('chat-view').style.display = 'none';
    document.getElementById('convo-view').style.display = 'flex';
};

//


function streamGeminiRes() {

    if (resChunkQueue.length === 0) {
        resInProg = false;
        return;
    };

    const area = document.getElementById('convo-area');
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

function renderUsrMsg(text) {
    const convoArea = document.getElementById('convo-area');
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.classList.add('user');
    messageElement.textContent = text;
    convoArea.appendChild(messageElement);
    convoArea.scrollTop = convoArea.scrollHeight;
};


// ℹ️: Wait for DOM, before attempting anything element-wise
// document.addEventListener('DOMContentLoaded', function () {
//     const sendMsgBtn = document.getElementById('send-btn');
//     const msgInput = document.getElementById('msg-input');

//     sendMsgBtn.addEventListener('click', function () {
//         const msg = msgInput.value;
//         msgInput.value = '';

//         renderUsrMsg(msg);

//         vscode.postMessage({
//             command: 'execNewMsg',
//             message: msg,
//         });
//     });
// });
