#!/usr/bin/env python

# Imports
import time
import appdash
from opentracing import Format, Tracer
# Appdash: Socket Collector
from appdash.sockcollector import RemoteCollector

# Create a remote appdash collector.
collector = RemoteCollector(debug=True)
collector.connect(host="localhost", port=7701)

# Create a tracer
tracer = appdash.create_new_tracer(collector)
tracer.register_required_propagators()
for i in range(0, 7):
    # Generate a few spans with some annotations.
    span = None
    # Name the span.
    if i == 0:
        span = tracer.start_span("TEST")
    else:
        span = tracer.start_span("SQL TEST")
    carrier={}
    tracer.inject(span.context,Format.TEXT_MAP,carrier)
    print carrier
    span.set_tag("query", "SELECT * FROM table_name;")
    span.set_tag("foo", "bar")

    if i % 2 == 0:
        span.log_event("Hello world!")

    span_context=tracer.extract(Format.TEXT_MAP,carrier)
    child_span = tracer.start_span("child", child_of=span_context)
    child_span.finish()

    span.finish(finish_time=time.time()+2)

# Close the collector's connection.
collector.close()
