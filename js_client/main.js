const net = require('net');

function fetchData(requestContext, host = '127.0.0.1', port = 8080) {
    return new Promise((resolve, reject) => {
        const client = new net.Socket();

        client.connect(port, host, () => {
            const jsonData = JSON.stringify(requestContext);
            const lengthBuffer = Buffer.alloc(4);
            lengthBuffer.writeUInt32BE(Buffer.byteLength(jsonData), 0);

            console.log('Sending request:', requestContext);
            client.write(lengthBuffer);
            client.write(jsonData);
        });

        client.on('data', (data) => {
            try {
                const parsedData = JSON.parse(data.toString());
                resolve(parsedData);
            } catch (error) {
                reject('Error parsing data: ' + error.message);
            }
            client.destroy();
        });

        client.on('close', () => {});

        client.on('error', (err) => {
            reject('Error occurred: ' + err.message);
        });
    });
}

const requestContext = {
    meta: {
        path: '/search',
        method: 'GET',
    },
    body: {
        query: "mr",
    },
};

fetchData(requestContext)
    .then(response => {
        console.log('Received response:', response);
    })
    .catch(error => {
        console.error('Error:', error);
    });


