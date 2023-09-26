import socket

PORT = 8000

if __name__ == "__main__":
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.bind(("", PORT))
        sock.listen(1)
        client_sock, client_addr = sock.accept()
        while True:
            payload, addr = client_sock.recvfrom(4096)
            client_sock.sendto(payload, ("", PORT))
    finally:
        sock.close()
