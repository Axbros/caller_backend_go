import websocket
import json
import threading

# WebSocket服务器地址
websocket_url = "ws://127.0.0.1:8080/ws"

# 发送单个心跳请求的函数
def send_heartbeat(url, data):
    try:
        # 创建WebSocket连接
        ws = websocket.create_connection(url)
        # 创建消息
        message = {
            "event": "heartbeat",
            "message": "",
            "data": "10000"+str(data),
            "key": "A8KpZu",
            "to": "",
            "type": "user"
        }
        # 发送消息
        ws.send(json.dumps(message))
        # ws.close()
        # 接收服务器响应（如果有的话）
        response = ws.recv()
        print(f"Received message from server for data {data}: " + response)
    except Exception as e:
        print(f"Error sending heartbeat for data {data}: {e}")
    finally:
        # 关闭连接
        print("关闭")
        

# 创建并启动线程来发送心跳请求
threads = []
for i in range(1):
    thread = threading.Thread(target=send_heartbeat, args=(websocket_url, i))
    threads.append(thread)
    thread.start()

# 等待所有线程完成
for thread in threads:
    thread.join()
