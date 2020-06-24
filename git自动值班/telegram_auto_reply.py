import random
from telethon import TelegramClient
import asyncio
import time
from telethon import events

# 个人密钥私密
api_id = xxxxxx
api_hash = "xxxxxxxxxxxxxxxxxxxx"
myTgId = 个人telegram账号

client = TelegramClient('mm1ofauto', api_id, api_hash)
message = ['你好，我正在休假中']

# 全局缓存
time_cache = {}

# 设定多久过期
max_expire_sec = 120

# 刷缓存
def refresh_cache(key):
    time_cache[key] = time.time()
    pass

def main():
    client.start()
    print("已启动监控...")

    @client.on(events.NewMessage(outgoing=True))
    async def _(event):
        # print("发消息:", event.message)
        if event.is_private:
            # 私聊，如果是主动发出，则刷新缓存
            refresh_cache(event.message.to_id.user_id)
        else:
            # 群聊，如果是主动发出，则刷新缓存
            refresh_cache(event.message.to_id.channel_id)
        pass

    @client.on(events.NewMessage(incoming=True))
    async def _(event):
        # print("收消息:", event.message)
        print("{},消息类型:{},发送ID:{},接受ID:{}\n消息内容:\n{}".format(
            time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(time.time())),
            event.is_private,
            event.message.to_id, event.message.from_id, event.message.message)
        )
        # print("缓存：", time_cache)

        # systemHour = time.localtime().tm_hour
        # if systemHour >= 10 and systemHour <= 24:
        #     print("工作时间中，关闭监控")
        #     return

        if event.is_private:
            # 效果就是回复之后，30s之后才开启自动回复
            leftS = time.time() - time_cache.get(event.message.from_id, 0)
            if int(leftS) < max_expire_sec:
                print("【{}】已自动回复过，剩余【{}】s".format(event.message.from_id, max_expire_sec - int(leftS)))
                return

            # 记录时间点
            refresh_cache(event.message.from_id)

            time.sleep(3)
            await client.send_message(event.message.from_id, message[random.randint(0, len(message) - 1)])
        else:
            # 这里一般是群组沟通
            # 判断是否提及某人
            if event.message.mentioned == True:
                for e in event.message.entities:
                    # print("群指定用户id：", e.user_id)
                    if e.user_id == myTgId:
                        # 效果就是回复之后，30s之后才开启自动回复
                        leftS = time.time() - time_cache.get(event.message.to_id.channel_id, 0)
                        if int(leftS) < max_expire_sec:
                            print("【{}】已自动回复过，剩余【{}】s".format(event.message.to_id.channel_id, max_expire_sec - int(leftS)))
                            return

                        refresh_cache(event.message.to_id.channel_id)

                        time.sleep(3)
                        await client.send_message(event.message.to_id.channel_id, message[random.randint(0, len(message) - 1)])
                        pass
            pass

    client.run_until_disconnected()

if __name__ == '__main__':
    main()
