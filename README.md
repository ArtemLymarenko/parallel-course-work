# Parallel Computing Course Work

This repo contains UML diagrams and source code for the parallel computing course work.

# Project Setup Guide

## Prerequisites

### Clone the repository to your local machine
```bash
   git clone https://github.com/ArtemLymarenko/parallel-course-work.git 
   cd parallel-course-work
```

## Python Installation (>=3.12.3)
1. **Download Python** from the official website: [python.org](https://www.python.org/downloads/)
    - Ensure you have Python version 3.12.3 or higher installed.

2. **Install `pip` Package Manager:**
    - For **Mac/Linux**:
      ```bash
      python -m ensurepip --default-pip
      ```
    - For **Windows**:
      ```bash
      py -m ensurepip --default-pip
      ```

3. **Install Python Dependencies:**
    - Navigate to the Python client directory:
      ```bash
      cd clients/python
      ```
    - Install the required packages:
      ```bash
      pip install -r requirements.txt
      ```

4. **Running Locust Client:**
    - To start the Locust client:
      ```bash
      locust
      ```
    - To run the server in a different mode:
      ```bash
      locust -f locustfile_open.py
      ```

5. **Additional Flags:**
   You can add `--master` and `--worker` flags to create separate processes for testing. These should be run in different console windows:
    - For the master process:
      ```bash
      locust --master
      ```
    - For worker processes:
      ```bash
      locust --worker
      ```

---

## Golang Installation (>=go1.22)
Download Golang from the official website: [go.dev](https://go.dev/dl/)

### Linux

1. Remove previous Go installation and extract archive:
   ```bash
     rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
   ```
   Note: May require root/sudo access. Do not extract into existing /usr/local/go directory.

2. Add to PATH in `$HOME/.profile` or `/etc/profile`:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   ```
   Note: Log out and back in for changes to take effect, or run `source $HOME/.profile`.

3. Verify installation:
   ```bash
   go version
   ```

### macOS
   1. Open and run the downloaded package installer
   2. Go will be installed to `/usr/local/go`
   3. `/usr/local/go/bin` will be added to PATH automatically
   4. Restart Terminal sessions for changes to apply
   5. Verify installation:
   ```bash
   go version
   ```

### Windows
   1. Run the downloaded MSI installer
   2. Default installation: Program Files or Program Files (x86)
   3. Close and reopen command prompts after installation
   4. Verify installation:
       - Click Start menu
       - Search for "cmd" and press Enter
       - Run:
   ```bash
   go version
   ```
---

# Setting Up Go Project
## Locally
1. **Install Go Dependencies:**
    - For **Linux/macOS**:
      ```bash
      cd server && go mod tidy && cd ../pkg && go mod tidy && cd ../clients/golang && go mod tidy && cd ../..
      ```
    - For **Windows**:
      ```bash
      cd server; go mod tidy; cd ../pkg; go mod tidy; cd ../clients/golang; go mod tidy; cd ../..
      ```

2. **Start Server and Client:**
    - For **Linux/macOS**:
        - To start the server:
          ```bash
          cd server && go run cmd/main.go
          ```
        - To start the client:
          ```bash
          cd clients/golang && go run cmd/main.go
          ```
    - For **Windows**:
        - To start the server:
          ```bash
          cd server; go run cmd/main.go
          ```
        - To start the client:
          ```bash
          cd clients/golang; go run cmd/main.go
          
          ```

## Using Docker
1. **Run Docker Compose in the root folder:**
   - To start the containers:  
   ```bash
     docker-compose up
     ```
   - To stop the containers:
   ```bash
     docker-compose down
     ```

>⚠️ **IMPORTANT:**
>
> Files for building an inverted index should be stored in the ```server/resources/data/``` folder.
>
> File for logs should be stored in the ```server/resources/logs``` foler and named ```logs.txt```.
---
## Open the clients

> To open Golang client use localhost:3000

> To open Locust client use localhost:8089
---

# Application Protocol Documentation

## Overview
This document describes the protocol for communicating with the server via TCP connection. All requests and responses use JSON serialization.

## Connection Protocol
## How to communicate?
- First, you need to serialize your request into special format using JSON, described below.
- Then you need to calculate `chunkSize` and `totalChunks` and before sending each request you need to send size of chunk and their total number.
1. chunkSize   - `4 bytes`
2. totalChunks - `4 bytes`
- To get response from the server, you first need to retrieve `chunkSize` and `totalChunks` and then all other data splitted by chunks.
- Finally, deserialize retrieved data from JSON.

### Request Format
All requests must include a `meta` object and may optionally include `body` and `connectionAlive` fields.

**Allowed Methods:** `GET`, `POST`, `DELETE`

```json
{
    "meta": {
        "path": { "type": "string" },      
        "method": { "type": "string" }    
    },
    "body": {},             
    "connectionAlive": { "type": "boolean" }
}
```

### Response Format
```json
{
    "status": "ResponseStatus",
    "body": {}
}
```

### Status Codes
```go
StatusOK                  ResponseStatus = 0
StatusProcessing          ResponseStatus = 1
StatusNotFound            ResponseStatus = 2
StatusBadRequest          ResponseStatus = 3
StatusInternalServerError ResponseStatus = 4
```

### Error Response
When an error occurs, the response body will contain a message:
```json
{
    "message": { "type": "string" }
}
```

## API Endpoints

### 1. Search Files
Search for files matching a specific query.

- **Path:** `index/search`
- **Method:** `GET`
- **Request Body:**
```json 
{
    "query": "string"
}
```
- **Response Body:**
```json 
{
    "files": ["string"]
}
```

### 2. Search Any
Search for files with more flexible matching.

- **Path:** `index/search-any`
- **Method:** `GET`
- **Request Body:**
```json 
{
    "query": "string"
}
```
- **Response Body:**
```json 
{
    "files": ["string"]
}
```

### 3. Get File Content
Retrieve the content of a specific file.

- **Path:** `/index/file`
- **Method:** `GET`
- **Request Body:**
```json 
{
    "fileName": "string"
}
```
- **Response Body:**
```json 
{
    "fileContent": "string"
}
```

### 4. Add File
Add a new file to the index.

- **Path:** `/index/file`
- **Method:** `POST`
- **Request Body:**
```json 
{
    "fileName": "string"
}
```

### 5. Remove File
Delete a file from the index.

- **Path:** `/index/file`
- **Method:** `DELETE`
- **Request Body:**
```json 
{
   "fileName": "string"
}
```


