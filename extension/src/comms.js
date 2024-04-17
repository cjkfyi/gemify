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

    const chat = msg.data.chat;
    const projID = chat.projID;
    const chatID = chat.chatID;
    const message = msg.data.msg;


    const url =
        'ws://localhost:8080/p/' + 
        projID + '/c/' + chatID + '/s';

    const ws = new WebSocket(url);

    await new Promise((resolve, reject) => {
        ws.onopen = () => resolve();
        ws.onerror = (err) => reject(err);
    });

    ws.send(message);

    ws.onmessage = (e) => {
        const data = JSON.parse(e.data);
        onChunkReceived(data);
    };

    ws.onclose = () => {
        const data = { content: "EOF"};
        onChunkReceived(data);
    };
};

export {
    getProjList,
    getChatList,
    getMsgList,
    getNewMsg,
};