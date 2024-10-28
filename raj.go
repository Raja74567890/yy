import socket
import sys
import threading
import time
import signal
from datetime import datetime, timedelta

BUFFER_SIZE = 16384
EXPIRATION_DATE = datetime(2024, 12, 31)
TARIFF_RATE = 0.05
stop_flag = threading.Event()

def check_expiration():
    if datetime.now() > EXPIRATION_DATE:
        print("This tool has expired.")
        sys.exit(1)

def calculate_tariff(data_size_kb):
    return data_size_kb * TARIFF_RATE

def send_udp_traffic(ip, port, duration):
    try:
        with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as sock:
            sock.connect((ip, port))
            data_sent_kb = 0
            start_time = time.time()

            while time.time() - start_time < duration and not stop_flag.is_set():
                sock.send(b"x" * BUFFER_SIZE)
                data_sent_kb += BUFFER_SIZE / 1024

            print(f"Total UDP data sent: {data_sent_kb} KB, Tariff: ${calculate_tariff(data_sent_kb):.2f}")
    except Exception as e:
        print(f"UDP traffic error: {e}")

def send_tcp_traffic(ip, port, duration):
    try:
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
            sock.connect((ip, port))
            data_sent_kb = 0
            start_time = time.time()

            while time.time() - start_time < duration and not stop_flag.is_set():
                sock.send(b"x" * BUFFER_SIZE)
                data_sent_kb += BUFFER_SIZE / 1024

            print(f"Total TCP data sent: {data_sent_kb} KB, Tariff: ${calculate_tariff(data_sent_kb):.2f}")
    except Exception as e:
        print(f"TCP traffic error: {e}")

def signal_handler(signum, frame):
    stop_flag.set()

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print(f"Usage: {sys.argv[0]} <IP> <PORT> <DURATION>")
        sys.exit(1)

    check_expiration()

    ip = sys.argv[1]
    port = int(sys.argv[2])
    duration = int(sys.argv[3])

    signal.signal(signal.SIGINT, signal_handler)
    signal.signal(signal.SIGTERM, signal_handler)

    udp_thread = threading.Thread(target=send_udp_traffic, args=(ip, port, duration))
    tcp_thread = threading.Thread(target=send_tcp_traffic, args=(ip, port, duration
