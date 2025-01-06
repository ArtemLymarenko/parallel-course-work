import client

if __name__ == "__main__":
    client = client.TcpSocketClient(host='127.0.0.1', port=8080, buff_size=2048)

    request = {
        "meta": {
            "path": '/index/search',
            "method": 'GET',
        },
        "body": {
            "query": "Bradly Brad Tom Harry Old Main man marry sue",
        },
        "connectionAlive": True,
    }

    client.connect()
    res = client.fetch_open_conn(request)
    print(res)
    res2 = client.fetch_open_conn(request)
    print(res2)
    client.close()

