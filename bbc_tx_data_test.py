#!/usr/bin/python
# -*- coding: UTF-8 -*-

import time
import uuid
import requests
import json
import sys

def str_to_hex(s):
    return ''.join([hex(ord(c)).replace('0x', '') for c in s])

def hex_to_str(s):
    result = ""
    arr16 = []
    for index in range(len(s)/2):
        arr16.append(int(s[index*2:index*2+2],16))
    return ''.join([chr(i) for i in arr16])

str_ = 'hello 大棒客为武汉加油'
str_16 = str_to_hex(str_)
#print str_16
#print hex_to_str(str_16)
str_time = hex(int(time.time())).replace('0x','')
str_uuid = str(uuid.uuid1()).replace('-','')
data = str_uuid+str_time+"00"+str_16
url = 'http://127.0.0.1:9902'
data = '{"id":1,"method":"sendfrom","jsonrpc":"2.0","params":{"from":"20g00pbrgnhrjp8m2w4p1h2gwpn6w4s1rsst7avwc4av02580n6m70wjm","to":"1965p604xzdrffvg90ax9bk0q3xyqn5zz2vc9zpbe3wdswzazj7d144mm","amount":100.00000000,"data":"%s"}}' % data
response = requests.post(url, data=data)
res = json.loads(response.text)
txid = res["result"]
data = '{"id":1,"method":"gettransaction","jsonrpc":"2.0","params":{"txid":"%s"}}' % txid
response = requests.post(url, data=data)
res = json.loads(response.text)
data_str = res["result"]["transaction"]["data"]

print "uuid1", data_str[0:8] + "-" + data_str[8:12] + "-" + data_str[12:16] + "-" + data_str[16:20] + "-" + data_str[20:32]
timeArray = time.localtime(int(data_str[32:40],16))
otherStyleTime = time.strftime("%Y-%m-%d %H:%M:%S", timeArray)
print "时间:", otherStyleTime
print "内容:", hex_to_str(data_str[42:])
