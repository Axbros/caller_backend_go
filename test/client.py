import asyncio
import websockets

async def send_message():
    uri = "ws://localhost:8080/ws"
    async with websockets.connect(uri) as websocket:
        #挂断电话 速度很快 话机向甲方传递数据
        message = '{"event":"heartbeat","message":"16600229957","data":"1","key":"testkey","type":"client"}'
        await websocket.send(message)
        response = await websocket.recv()
        print(f"Sent: {message}, Received: {response}")

asyncio.get_event_loop().run_until_complete(send_message())