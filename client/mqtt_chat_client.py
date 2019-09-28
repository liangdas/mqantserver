# -*- coding: utf-8 -*-
'''
Created on 17/2/16.
@author: love
'''
import paho.mqtt.client as mqtt
import json
import ssl

def on_connect(client, userdata, flags, rc):
    print("Connected with result code %d"%rc)
    client.publish("Login/HD_Login/1", json.dumps({"userName": user, "passWord": "Hello,anyone!"}),qos=0,retain=False)



def on_message(client, userdata, msg):
    print ('---------------')
    print ("topic   :"+msg.topic)
    print ("payload :"+msg.payload)
    client.subscribe([("chat",2),("aaa",2)])
    client.unsubscribe(["chat"])
    #client.publish("login/addUser", json.dumps({"user": user, "say": "Hello,anyone!"}),qos=2,retain=False)
    #print(msg.topic+":"+str(msg.payload.decode()))
    #print(msg.topic+":"+msg.payload.decode())
    #payload = json.loads(msg.payload.decode())
    #print(payload.get("user")+":"+payload.get("say"))
def mylog(self,userdata,level, buf):
    print buf

if __name__ == '__main__':
    client = mqtt.Client(protocol=mqtt.MQTTv31)
    client.username_pw_set("admin", "password")  # 必须设置，否则会返回「Connected with result code 4」
    client.on_connect = on_connect
    client.on_message = on_message
    #链接测试服务器 需要用tls请求 python  tls功能比较弱。
    # 需要一个证书，这里使用的这个网站提供的证书https://curl.haxx.se/docs/caextract.html
    HOST = "mqant.com"
    # client.tls_set(ca_certs="caextract.pem", certfile=None, keyfile=None, cert_reqs=ssl.CERT_REQUIRED,
    #                tls_version=ssl.PROTOCOL_TLSv1, ciphers=None)
    client.connect(HOST, 3563, 60)
    #client.loop_forever()

    user = raw_input("请输入用户名:")
    client.user_data_set(user)

    client.loop_start()

    while True:
        s = raw_input("请先输入'join'加入房间,然后输入任意聊天字符:\n")
        if s:
            if s=="join":
                client.publish("Chat/HD_JoinChat/2", json.dumps({"roomName": "mqant"}),qos=0,retain=False)
            elif s=="start":
                    client.publish("Master/HD_Start_Process/2", json.dumps({"ProcessID": "001"}),qos=0,retain=False)
            elif s=="stop":
                client.publish("Master/HD_Stop_Process/2", json.dumps({"ProcessID": "001"}),qos=0,retain=False)
            else:
                client.publish("Chat/HD_Say/2", json.dumps({"roomName": "mqant","from":user,"target":"*","content": s}),qos=0,retain=False)