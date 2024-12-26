import socket
import struct
import json

class TcpSocketClient:
    def __init__(self, host, port, buff_size):
        self.host = host
        self.port = port
        self.buff_size = buff_size
        self.sock = None

    def connect(self):
        if not self.sock:
            self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.sock.connect((self.host, self.port))

    def close(self):
        if self.sock:
            self.sock.close()
            self.sock = None

    def write_int32_to_buffer(self, value):
        return struct.pack('>I', value)

    def parse_buffered_int32(self, buffer):
        return struct.unpack('>I', buffer)[0]

    def send_request_raw(self, sock, data):
        request_bin = json.dumps(data).encode('utf-8')
        request_len = len(request_bin)
        total_chunks = (request_len + self.buff_size - 1) // self.buff_size

        sock.sendall(self.write_int32_to_buffer(self.buff_size))
        sock.sendall(self.write_int32_to_buffer(total_chunks))

        offset = 0
        while offset < request_len:
            end = min(offset + self.buff_size, request_len)
            sock.sendall(request_bin[offset:end])
            offset = end

    def read_response(self, sock):
        chunk_size_buff = sock.recv(4)
        if not chunk_size_buff:
            raise ConnectionError("Failed to read chunk size.")
        chunk_size = self.parse_buffered_int32(chunk_size_buff)

        total_chunks_buff = sock.recv(4)
        if not total_chunks_buff:
            raise ConnectionError("Failed to read total chunks.")
        total_chunks = self.parse_buffered_int32(total_chunks_buff)

        response_data = bytearray()
        for _ in range(total_chunks):
            chunk = sock.recv(chunk_size)
            if not chunk:
                raise ConnectionError("Failed to read a chunk.")
            response_data.extend(chunk)

        return json.loads(response_data.decode('utf-8'))

    def fetch(self, request):
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
            sock.connect((self.host, self.port))
            self.send_request_raw(sock, request)
            return self.read_response(sock)

    def fetch_open_conn(self, request):
        self.send_request_raw(self.sock, request)
        return self.read_response(self.sock)