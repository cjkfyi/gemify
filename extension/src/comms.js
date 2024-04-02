Object.assign(global, { WebSocket: require('ws') });

async function sendGemifyMsg(msg, id) {
    const url = 'ws://localhost:8000/ws/chat/' + id;
    console.log(url);

    // Establish WebSocket Connection
    const ws = new WebSocket(url);

    // Promise to resolve when the WebSocket connection is open
    const openPromise = new Promise((resolve, reject) => {
        ws.onopen = () => resolve();
        ws.onerror = (err) => reject(err);
    });

    await openPromise; // Wait for the connection to open

    // Send the Message
    ws.send(JSON.stringify({ message: msg }));

    // Function to handle incoming messages
    const receiveMessagePromise = new Promise((resolve, reject) => {
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data); 
            resolve(data);
        };
        ws.onerror = (err) => reject(err);
    });

    const data = await receiveMessagePromise; // Wait for the response
    return data;
}


async function newConvo() {
    const res = await fetch('http://localhost:8000/chat', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
    });

    if (!res.ok) {
        throw new Error(`HTTP Error: ${res.status}`);
    }

    const data = await res.json();
    return data;
}

export { sendGemifyMsg, newConvo }