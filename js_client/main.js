const net = require('net');

const requestContext = {
    meta: {
        path: 'status',
        method: 'POST',
    },
    body: 'Sample data',
};

const HOST = '127.0.0.1';
const PORT = 8080;

const client = new net.Socket();

client.connect(PORT, HOST, () => {
    console.log('Підключення встановлено');

    const jsonData = JSON.stringify(requestContext);
    const lengthBuffer = Buffer.alloc(4);
    lengthBuffer.writeUInt32BE(Buffer.byteLength(jsonData), 0);

    console.log('Відправлення даних:', lengthBuffer);
    client.write(lengthBuffer);

    console.log('Відправлення даних:', jsonData);
    client.write(jsonData);
});

client.on('data', (data) => {
    console.log('Отримано відповідь від сервера:', data.toString());

    client.destroy();
});

client.on('close', () => {
    console.log("З'єднання закрито");
});

client.on('error', (err) => {
    console.error('Помилка:', err.message);
});
