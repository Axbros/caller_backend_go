import websocket
import json
import time
# WebSocket服务器地址
websocket_url = "ws://127.0.0.1:8080/ws/"

# 发送单个心跳请求的函数
def send_heartbeat(url):
    try:
        # 创建WebSocket连接
        ws = websocket.create_connection(url)
        while True:
            # 创建消息
            message = {
                "event": "heartbeat",
                "message": "",
                "data": "888",
                "key": "A8KpZu",
                "to": "",
                "type": "user"
            }
            # 发送消息
            ws.send(json.dumps(message))
            
            try:
                response = ws.recv()
                print(f"Received message from server: {response}")
            except websocket.WebSocketConnectionClosedException:
                print("Connection closed by server.")
                break
            time.sleep(8)
            print("开心下一轮心跳")
    except Exception as e:
        print(f"Error sending heartbeat: {e}")


if __name__ == "__main__":
    send_heartbeat(websocket_url)