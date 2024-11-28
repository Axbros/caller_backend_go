import asyncio
import websockets
import json

async def connect_and_disconnect():
    uri = "ws://localhost:8080/ws/"
    try:
        async with websockets.connect(uri) as websocket:
            print("WebSocket连接已成功建立")
            # 构造要发送的消息
            message = {
                "event": "heartbeat",
                "message": "",
                "type": "user",
                "data": "88",
                "key": "RYZ4SM",
                "to": "",
                "from": "88"
            }
            # 将字典转换为JSON字符串格式后发送
            await websocket.send(json.dumps(message))
            print("已向服务器发送心跳消息")
            await asyncio.sleep(20)  # 等待10秒
            await websocket.close()  # 关闭WebSocket连接
            print("WebSocket连接已关闭")
    except websockets.exceptions.ConnectionError as e:
        print(f"连接出现错误: {e}")

asyncio.run(connect_and_disconnect())