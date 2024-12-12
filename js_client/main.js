const net = require('net');

const requestContext = {
    meta: {
        path: '/search',
        method: 'GET',
    },
    body: {
        query: "mr",
    },
};

const HOST = '127.0.0.1';
const PORT = 8080;

const client = new net.Socket();

client.connect(PORT, HOST, () => {
    console.log('Connection established');

    const jsonData = JSON.stringify(requestContext);
    const lengthBuffer = Buffer.alloc(4);
    lengthBuffer.writeUInt32BE(Buffer.byteLength(jsonData), 0);

    console.log('Sending:', lengthBuffer);
    client.write(lengthBuffer);

    console.log('Sending:', jsonData);
    client.write(jsonData);
});

client.on('data', (data) => {
    console.log('Received answer from server:', JSON.parse(data.toString()));

    client.destroy();
});

client.on('close', () => {
    console.log("Connection closed");
});

client.on('error', (err) => {
    console.error('Error occurred:', err.message);
});
