import random
import os
import time

# A basictracer-specific instance of guid_rng. See _fork_guard_pid.
guid_rng = random.Random()

# The current pid. If the process forks (which happens, for instance, in
# uwsgi), we consult _fork_guard_pid and re-seed guid_rng accordingly.
_fork_guard_pid = 0


def generate_id():
    global _fork_guard_pid

    # Microbenchmarks suggest that os.getpid() takes less than 0.1 microsecond.
    pid = os.getpid()
    if (_fork_guard_pid == 0) or (_fork_guard_pid != pid):
        _fork_guard_pid = pid
        guid_rng.seed(int(1000000 * time.time()) ^ pid)
    return guid_rng.getrandbits(64) - 1
