#!/usr/bin/python
import os
import time
import json
import re

open_falcon_message = [{
        "endpoint" : "",
        "metric": "service_monitor",
        "timestamp": 0,
        "step": 60,
        "value": 1,
        "counterType": "GAUGE",
        "tags": "",
        }]

st1 = os.popen("cat /root/test.log")
ts = int(time.time())
open_falcon_message[0]["timestamp"] = ts
n = 0
while True:
    i = st1.readline()
    if i:
        m = re.search(r"stop",i)
        if m:
            n = n + 1
    else:
        break
open_falcon_message[0]["value"] = n
print (json.dumps(open_falcon_message))
