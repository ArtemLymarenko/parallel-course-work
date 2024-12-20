const net = require('net');

class TCPClient {
    constructor(host = '127.0.0.1', port = 8080) {
        this.host = host;
        this.port = port;
        this.client = new net.Socket();
        this.buffer = Buffer.alloc(0);
        this.requests = []

        this.client.connect(port, host, () => {
            console.log(`Connected to server at ${host}:${port}`);
        });

        this.client.on('data', (data) => {
            try {
                const response = JSON.parse(data.toString());
                const { resolve } = this.requests.shift();
                resolve(response);
            } catch (error) {
                const { reject } = this.requests.shift();
                reject(error.message);
            }
        });

        this.client.on('error', (err) => {
            console.error('Connection error:', err.message);
        });

        this.client.on('close', () => {
            console.log('Connection closed');
        });
    }

    sendRequest(requestContext) {
        return new Promise((resolve, reject) => {
            const jsonData = JSON.stringify(requestContext);
            const lengthBuffer = Buffer.alloc(4);
            lengthBuffer.writeUInt32BE(Buffer.byteLength(jsonData), 0);

            this.requests.push({ resolve, reject });
            this.client.write(lengthBuffer);
            this.client.write(jsonData);
        });
    }

    close() {
        this.client.end()
    }
}

const client = new TCPClient();

const requestContext = {
    meta: {
        path: '/index/search',
        method: 'GET',
    },
    body: {
        query: "mr",
    },
    connectionAlive: true,
};

const requestContext2 = {
    meta: {
        path: '/index/search',
        method: 'GET',
    },
    body: {
        query: "Hello asdlas as",
    },
    connectionAlive: true,
};

(async function () {
    try {
        const response1 = await client.sendRequest(requestContext);
        console.log('Received response 1:', response1);

        const response2 = await client.sendRequest(requestContext);
        console.log('Received response 2:', response2);

        const response3 = await client.sendRequest(requestContext2);
        console.log('Received response 3:', response3);

        const response4 = await client.sendRequest(requestContext);
        console.log('Received response 4:', response4);

        const response5 = await client.sendRequest(requestContext);
        console.log('Received response 5:', response5);

        const response6 = await client.sendRequest(requestContext2);
        console.log('Received response 6:', response6);
    } catch (error) {
        console.error('Error:', error);
    } finally {
        client.close()
    }
})();