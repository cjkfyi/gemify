async function sendGeminiReq(msg) {
    const res = await fetch('http://localhost:8080/api/gemini', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ message: msg })
    });

    if (!res.ok) {
        throw new Error(`HTTP Error: ${res.status}`);
    }

    const data = await res.json();
    return data.message;
}

export { sendGeminiReq }