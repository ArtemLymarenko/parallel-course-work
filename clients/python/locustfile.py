import json
import time
import random

from locust import TaskSet, task, User, events, between
from client import TcpSocketClient

class TcpSocketResponse:
    details = None
    request_name = None
    response_time = None
    content_length = None
    data = None

    def __init__(self, request_details):
        self.details = request_details
        self.request_name = request_details['name']
        self.response_time = request_details['time']
        self.content_length = request_details['length']
        self.data = request_details['bytes']

    def success(self):
        events.request_success.fire(
            request_type=self.details['type'],
            name=self.details['name'],
            response_time=self.details['time'],
            response_length=self.details['length']
        )

    def failure(self, message):
        events.request_failure.fire(
            request_type=self.details['type'],
            name=self.details['name'],
            response_time=self.details['time'],
            exception=message
        )

class TcpLocustClient:
    def __init__(self, host, port, buff_size):
        self.client = TcpSocketClient(host, port, buff_size)

    def send_request(self, request_name, request, catch_response=False):
        start_time = time.time()
        try:
            response_data = self.client.fetch(request)
            response_time = int((time.time() - start_time) * 1000)
            content_length = len(json.dumps(response_data).encode('utf-8'))

            response = TcpSocketResponse({
                'type': 'tcp',
                'name': request_name,
                'time': response_time,
                'length': content_length,
                'bytes': response_data
            })

            if catch_response:
                return response
            else:
                events.request.fire(
                    request_type='tcp',
                    name=request_name,
                    response_time=response_time,
                    response_length=content_length,
                    exception=None
                )
                return response_data
        except Exception as e:
            response_time = int((time.time() - start_time) * 1000)
            if catch_response:
                response = TcpSocketResponse({
                    'type': 'tcp',
                    'name': request_name,
                    'time': response_time,
                    'length': 0,
                    'bytes': None
                })
                response.failure(str(e))
                return response
            else:
                events.request.fire(
                    request_type='tcp',
                    name=request_name,
                    response_time=response_time,
                    response_length=0,
                    exception=e
                )


def generate_random_sentence(word_count=6):
    words = [
        "apple", "banana", "grape", "orange", "pear", "cherry", "melon", "kiwi", "peach", "plum",
        "dog", "cat", "bird", "fish", "rabbit", "elephant", "tiger", "lion", "zebra", "giraffe",
        "car", "bike", "bus", "train", "plane", "boat", "truck", "motorcycle", "scooter", "subway",
        "happy", "sad", "angry", "excited", "nervous", "joyful", "bored", "confused", "surprised", "calm",
        "computer", "keyboard", "monitor", "mouse", "screen", "laptop", "tablet", "phone", "internet", "website",
        "music", "movie", "book", "magazine", "newspaper", "song", "album", "concert", "theater", "stage",
        "coffee", "tea", "water", "juice", "milk", "soda", "beer", "wine", "cocktail", "whiskey",
        "mountain", "river", "ocean", "forest", "desert", "lake", "island", "beach", "volcano", "canyon",
        "city", "town", "village", "suburb", "neighborhood", "street", "avenue", "park", "plaza", "square",
        "school", "university", "college", "library", "bookstore", "museum", "gallery", "theater", "stadium", "gym",
        "family", "friend", "team", "group", "partner", "colleague", "boss", "employee", "student", "teacher",
        "dog", "cat", "bird", "fish", "rabbit", "hamster", "snake", "lizard", "turtle", "frog"
    ]

    sentence = ' '.join(random.choice(words) for _ in range(word_count))

    return sentence

def generate_random_file():
    files = [
        "resources/test/0_2.txt", "resources/test/1_3.txt", "resources/test/2_3.txt",
        "resources/test/3_4.txt", "resources/test/4_4.txt", "resources/test/5_4.txt",
        "resources/test/6_3.txt", "resources/test/7_1.txt", "resources/test/8_2.txt",
        "resources/test/9_4.txt", "resources/test/10_3.txt", "resources/test/11_3.txt",
        "resources/test/15_2.txt", "resources/test/20_1.txt", "resources/test/29_3.txt",
    ]

    return random.choice(files)

class UserBehavior(TaskSet):
    def on_start(self):
        self.tcp_client = TcpLocustClient(self.user.host, self.user.port, buff_size=2048)

    @task(4)
    def search_index(self):
        request = {
            "meta": {
                "path": "/index/search",
                "method": "GET"
            },
            "body": {
                "query": generate_random_sentence(10)
            },
        }
        self.tcp_client.send_request("Search Index", request)

    @task(1)
    def add_file(self):
        request = {
            "meta": {
                "path": "/index/file",
                "method": "POST"
            },
            "body": {
                "fileName": generate_random_file(),
            },
        }
        self.tcp_client.send_request("Add File", request)

    @task(1)
    def remove_file(self):
        request = {
            "meta": {
                "path": "/index/file",
                "method": "DELETE"
            },
            "body": {
                "fileName": generate_random_file()
            },
        }
        self.tcp_client.send_request("Remove File", request)

class TcpUser(User):
    tasks = [UserBehavior]
    host = "127.0.0.1"
    port = 8080
    wait_time = between(1, 2)
