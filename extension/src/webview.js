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

        case 'returnMsg':
            resChunkQueue.push(msg.data);
            if (!resInProg) {
                resInProg = true;
                streamGeminiRes();
            }
            break;
    }
});

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
    const convoArea = document.getElementById('convo-area');
    const messageElement = document.createElement('div');
    messageElement.classList.add('message');
    messageElement.classList.add('user');
    messageElement.textContent = text;
    convoArea.appendChild(messageElement);
    convoArea.scrollTop = convoArea.scrollHeight;
};

//

function listRecentProjs(list) {
    const convoListEl = document.getElementById('convo-list');

    // Breaks if doesn't exist
    if (list.length === 0) {
        // ℹ️: Problematic alongside of the `show more...` el
        convoListEl.innerHTML = '<li>Create your first conversation?</li>';
        return;
    }

    const displayCount = 4;
    const displayConversations = list.slice(0, displayCount);

    displayConversations.forEach((proj, index) => {
        // Create Tile Elements
        const listItem = document.createElement('li');
        const tileLink = document.createElement('a');
        const tileName = document.createElement('h3');
        const tileDesc = document.createElement('p');

        // Add Styling Classes
        listItem.classList.add('convo-tile');
        tileLink.classList.add('convo-tile-link');
        tileName.classList.add('convo-tile-name');
        tileDesc.classList.add('convo-tile-desc');

        // Populate Tile Content
        tileName.textContent = proj.name;
        tileDesc.textContent = proj.desc;

        // Assemble Tile 
        tileLink.href = '#';
        tileLink.appendChild(tileName);
        tileLink.appendChild(tileDesc);
        listItem.appendChild(tileLink);
        convoListEl.appendChild(listItem);

        listItem.addEventListener('click', () => {
            vscode.postMessage({
                command: 'execConvoView',
                data: proj.projID,
            });
        });
    });
};

//

function newProjectAction() {
    vscode.postMessage({
        command: 'execNewProj',
    });
    document.getElementById('home-view').style.display = 'none';
    document.getElementById('new-proj-view').style.display = 'flex';
};

function returnHome() {
    vscode.postMessage({
        command: 'execReturnHome',
    });
    // document.getElementById('project-view').style.display = 'none';
    document.getElementById('new-proj-view').style.display = 'none';
    //
    document.getElementById('home-view').style.display = 'flex';
};

// function reopenConvoView() {
//     document.getElementById('home-view').style.display = 'none';
//     //
//     // document.getElementById('convo-view').style.display = 'flex';
// };
