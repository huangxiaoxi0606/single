
#!/usr/bin/python
# -*- coding: utf-8 -*-
# author:liaozhj
# datetime:2018/11/1 10:28
# software: PyCharm
from locust import HttpLocust, task, TaskSet
import time
import sys
import auto_stop
import error_csv
import master
import pushgate
from redis import Redis

class WorkTaskSet(TaskSet):
    def __init__(self, parent):
        super(WorkTaskSet, self).__init__(parent)

    @task
    def getworktask(self):
        """GET类型的请求"""
        pass

class WorkLocust(HttpLocust):
    def __init__(self):
        super(WorkLocust, self).__init__()

    task_set = WorkTaskSet
    # 可修改默认多个任务间的等待时间
    min_wait = 10
    max_wait = 10
    # todo:修改成要测试的服务器
    host = "http://10.18.97.191:8080"


    #在redis中插入总共的用户数量
    record_user=False
    for i in sys.argv:
        if record_user:
            r = Redis(host='10.18.98.136', password="username_redis", decode_responses=True)
            r.set("gotpc_user_num",int(i))
            break
        if i=="-c":
            record_user=True







if __name__ == '__main__':
    w = WorkLocust()
    w.run()
