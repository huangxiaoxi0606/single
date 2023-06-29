#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
@Author  : 丸子
@Email   : zhongj@yoozoo.com
@Time    : 2018/11/6 11:25
@File    : auto_stop.py
"""

from locust import events, runners

# 保存已经停止的clientid,当全部client都停止后，发出停止指令，控制全部client退出
stoped_salve = []


def auto_stop(client_id, data):
    if isinstance(runners.locust_runner, runners.MasterLocustRunner):
        if runners.locust_runner.state == runners.STATE_RUNNING and data['user_count'] == 0:
            if client_id not in stoped_salve:
                stoped_salve.append(client_id)
                if len(stoped_salve) == runners.locust_runner.slave_count:
                    print("auto client stoped. exiting.")
                    runners.locust_runner.quit()
            else:
                print("client: {} is stoped.".format(client_id))


events.slave_report += auto_stop