const WebSocket = require('ws');

async function sendGemifyMsg(msg, id, onChunkReceived) {

    const url = 'ws://localhost:8000/ws/chat/' + id;
    const ws = new WebSocket(url);

    await new Promise((resolve, reject) => {
        ws.onopen = () => resolve();
        ws.onerror = (err) => reject(err);
    });

    ws.send(JSON.stringify({ message: msg }));

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        onChunkReceived(data.content);
    };
};

async function sendNewConvo() {
    const res = await fetch('http://localhost:8000/chat', {
        method: 'POST',
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

async function getProjList() {
    const res = await fetch('http://localhost:8080/projects', {
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

export {
    sendGemifyMsg,
    sendNewConvo,
    getProjList,
};