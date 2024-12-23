import client

if __name__ == "__main__":
    client = client.TcpSocketClient(host='127.0.0.1', port=8080, buff_size=2048)

    request = {
        "meta": {
            "path": '/index/search',
            "method": 'GET',
        },
        "body": {
            "query": "mr",
        },
        "connectionAlive": True,
    }

    client.fetch(request)
