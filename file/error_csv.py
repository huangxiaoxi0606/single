#!/usr/bin/python
# -*- coding: utf-8 -*-
from locust import events, runners
import six
from six.moves import StringIO
import csv
import gevent

print("load error_csv file")

SAVE_INFO_ERROR = 1

def save_error():
    data = StringIO()
    writer = csv.writer(data)
    writer.writerow(["occurences", "Error"])
    i = 0
    for error in six.itervalues(runners.locust_runner.stats.errors):
        i = i+1
        if i >100 :
            break
        writer.writerow([error.occurences, error.to_name()])
    data.flush()
    data.seek(0)
    d = data.getvalue()
    with open('/data/load/master/locust_error.csv', 'w', encoding='utf-8') as ll:
        ll.write(d)
    data.close()





def error_start():
    if isinstance(runners.locust_runner, runners.MasterLocustRunner):
        print("master start save info to pushgate")
        global error_g
        error_g = gevent.spawn(send_error)


def send_error():
    while True:
        save_error()
        gevent.sleep(SAVE_INFO_ERROR)


def error_quit():
    if isinstance(runners.locust_runner, runners.MasterLocustRunner):
        print("master stop save info pushgate")
        global error_g
        error_g.kill(block=True)
        gevent.sleep(2)
        save_error()

events.master_start_hatching += error_start
events.quitting += error_quit