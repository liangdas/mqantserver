from __future__ import absolute_import

from threading import Lock
import time

from opentracing import Span
from opentracing.ext import tags


class BasicSpan(Span):
    """BasicSpan is a thread-safe implementation of opentracing.Span.
    """

    def __init__(
            self,
            tracer,
            operation_name=None,
            context=None,
            parent_id=None,
            tags=None,
            start_time=None):
        super(BasicSpan, self).__init__(tracer, context)
        self._tracer = tracer
        self._lock = Lock()

        self.operation_name = operation_name
        self.start_time = start_time
        self.parent_id = parent_id
        self.tags = tags if tags is not None else {}
        self.duration = -1
        self.logs = []

    def set_operation_name(self, operation_name):
        with self._lock:
            self.operation_name = operation_name
        return super(BasicSpan, self).set_operation_name(operation_name)

    def set_tag(self, key, value):
        with self._lock:
            if key == tags.SAMPLING_PRIORITY:
                self.context.sampled = value > 0
            if self.tags is None:
                self.tags = {}
            self.tags[key] = value
        return super(BasicSpan, self).set_tag(key, value)

    def log_kv(self, key_values, timestamp=None):
        with self._lock:
            self.logs.append(LogData(key_values, timestamp))
        return super(BasicSpan, self).log_kv(key_values, timestamp)

    def finish(self, finish_time=None):
        with self._lock:
            finish = time.time() if finish_time is None else finish_time
            self.duration = finish - self.start_time
            self._tracer.record(self)

    def set_baggage_item(self, key, value):
        new_context = self._context.with_baggage_item(key, value)
        with self._lock:
            self._context = new_context
        return self

    def get_baggage_item(self, key):
        with self._lock:
            return self.context.baggage.get(key)


class LogData(object):
    def __init__(
            self,
            key_values,
            timestamp=None):
        self.key_values = key_values
        self.timestamp = time.time() if timestamp is None else timestamp
