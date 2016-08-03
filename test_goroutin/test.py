import os
import time


def get_pid_status(fn):
    ps = {}
    with open(fn) as f:
        for line in f:
            fields = line.split(":")
            if fields[0] == "Name":
                ps["Name"] = fields[1].strip()
            if fields[0] == "State":
                ps["State"] = fields[1].strip()
    return ps


def all_procs():
    src = "/tmp/test_data"
    dirs = os.listdir(src)
    procs = []
    for d in dirs:
        pd = os.path.join(src, d)
        if os.path.isdir(pd):
            sf = os.path.join(pd, "status")
            if os.path.exists(sf):
                procs.append(get_pid_status(sf))
    return procs

if __name__ == '__main__':
    start_time = time.time()
    procs = all_procs()
    print time.time() - start_time
    print len(procs)
