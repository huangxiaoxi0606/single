#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
@Author  : 丸子
@Email   : zhongj@yoozoo.com
@Time    : 2018/11/14 17:49
@File    : pushgate.py
"""
import gevent
from locust import events
from locust import runners
import six
from prometheus_client import push_to_gateway, Gauge, CollectorRegistry, Counter, delete_from_gateway

push_gateway = 'http://10.18.97.253:9091'

SAVE_INFO_PUSHGATE = 1
push_g = None
push_reportid = LOCUSTREPORTID
push_cr = CollectorRegistry()

push_user_count = Gauge('locust_user_count', '当前用户数', ['push_reportid'], registry=push_cr)
push_num_requests = Gauge('locust_num_requests', '总次数', ['push_reportid', 'testname'], registry=push_cr)
push_num_failures = Gauge('locust_num_failures', '总错误数', ['push_reportid', 'testname'], registry=push_cr)
push_current_rps = Gauge('locust_current_rps', '当前RPS', ['push_reportid', 'testname'], registry=push_cr)
push_avg_response_time = Gauge('locust_avg_response_time', '平均响应时间', ['push_reportid', 'testname'], registry=push_cr)
push_median_response_time = Gauge('locust_median_response_time', '中值响应时间', ['push_reportid', 'testname'],
                                  registry=push_cr)
push_max_response_time = Gauge('locust_max_response_time', '最大响应时间', ['push_reportid', 'testname'], registry=push_cr)
push_min_response_time = Gauge('locust_min_response_time', '最小响应时间', ['push_reportid', 'testname'], registry=push_cr)


def pushgate_sort_stats(stats):
    return [stats[key] for key in sorted(six.iterkeys(stats))]


# 测试开始状态，开始后才上报数据。
def send_pushgate():
    global push_cr
    global push_user_count
    global push_num_requests
    global push_num_failures
    global push_current_rps
    global push_avg_response_time
    global push_median_response_time
    global push_max_response_time
    global push_min_response_time
    while True:
        for s in pushgate_sort_stats(runners.locust_runner.request_stats):
            push_num_requests.labels(push_reportid, s.name).set(s.num_requests or 0)
            push_num_failures.labels(push_reportid, s.name).set(s.num_failures or 0)
            push_current_rps.labels(push_reportid, s.name).set(s.current_rps or 0)
            push_avg_response_time.labels(push_reportid, s.name).set(s.avg_response_time or 0)
            push_median_response_time.labels(push_reportid, s.name).set(s.median_response_time or 0)
            push_max_response_time.labels(push_reportid, s.name).set(s.max_response_time or 0)
            push_min_response_time.labels(push_reportid, s.name).set(s.min_response_time or 0)
        push_user_count.labels(push_reportid).set(runners.locust_runner.user_count)
        print(push_user_count)
        push_to_gateway(push_gateway, job="oldlocusttest"+str(push_reportid), registry=push_cr)
        gevent.sleep(SAVE_INFO_PUSHGATE)


def pushgate_start():
    if isinstance(runners.locust_runner, runners.MasterLocustRunner):
        print("master start save info to pushgate")
        global push_g
        push_g = gevent.spawn(send_pushgate)


events.master_start_hatching += pushgate_start


def pushgate_quit():
    if isinstance(runners.locust_runner, runners.MasterLocustRunner):
        print("master stop save info pushgate")
        global push_g
        push_g.kill(block=True)
        gevent.sleep(2)
        delete_from_gateway(push_gateway, 'oldlocusttest'+str(push_reportid))


events.quitting += pushgate_quit