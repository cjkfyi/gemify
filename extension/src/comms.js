const WebSocket = require('ws');

async function getProjList() {
    
    const res = await fetch(
        'http://127.0.0.1:8080/projects', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        },
    );

    if (!res.ok) {
        throw new Error(`HTTP Error: ${res.status}`);
    };

    const data = await res.json();
    return data;
};

async function getChatList(projID) {
    
    const url = 'http://127.0.0.1:8080/p/' + 
        projID + '/chats';

    const res = await fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
    });

    if (!res.ok) {
        throw new Error(`HTTP Error: ${res.status}`);
    };

    const data = await res.json();
    return data;
};

async function getMsgList(projID, chatID) {

    const url = 'http://127.0.0.1:8080/p/' + 
        projID + '/c/' + chatID + '/history';

    const res = await fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        },
    });

    if (!res.ok) {
        throw new Error(`HTTP Error: ${res.status}`);
    };

    const data = await res.json();
    return data;
};

async function getNewMsg(msg, onChunkReceived) {

    const data = msg.data
    const projID = data.chat.projID
    const chatID = data.chat.chatID

    const url =
        'ws://localhost:8080/p/' + 
        projID + '/c/' + chatID + '/s';

    const ws = new WebSocket(url);

    await new Promise((resolve, reject) => {
        ws.onopen = () => resolve();
        ws.onerror = (err) => reject(err);
    });

    ws.send(JSON.stringify({ message: data.message }));

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        onChunkReceived(data.content);
    };
};

export {
    getProjList,
    getChatList,
    getMsgList,
    getNewMsg,
};