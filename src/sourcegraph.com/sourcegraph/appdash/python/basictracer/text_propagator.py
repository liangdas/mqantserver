from __future__ import absolute_import

from opentracing import SpanContextCorruptedException
from .context import SpanContext
from .propagator import Propagator

prefix_tracer_state = 'ot-tracer-'
prefix_baggage = 'ot-baggage-'
field_name_trace_id = prefix_tracer_state + 'traceid'
field_name_span_id = prefix_tracer_state + 'spanid'
field_name_sampled = prefix_tracer_state + 'sampled'
field_count = 3


class TextPropagator(Propagator):
    """A BasicTracer Propagator for Format.TEXT_MAP."""

    def inject(self, span_context, carrier):
        carrier[field_name_trace_id] = '{0:x}'.format(span_context.trace_id)
        carrier[field_name_span_id] = '{0:x}'.format(span_context.span_id)
        carrier[field_name_sampled] = str(span_context.sampled).lower()
        if span_context.baggage is not None:
            for k in span_context.baggage:
                carrier[prefix_baggage+k] = span_context.baggage[k]

    def extract(self, carrier):  # noqa
        count = 0
        span_id, trace_id, sampled = (0, 0, False)
        baggage = {}
        for k in carrier:
            v = carrier[k]
            k = k.lower()
            if k == field_name_span_id:
                span_id = int(v, 16)
                count += 1
            elif k == field_name_trace_id:
                trace_id = int(v, 16)
                count += 1
            elif k == field_name_sampled:
                if v == str(True).lower():
                    sampled = True
                elif v == str(False).lower():
                    sampled = False
                else:
                    raise SpanContextCorruptedException()
                count += 1
            elif k.startswith(prefix_baggage):
                baggage[k[len(prefix_baggage):]] = v

        if count != field_count:
            raise SpanContextCorruptedException()

        return SpanContext(
            span_id=span_id,
            trace_id=trace_id,
            baggage=baggage,
            sampled=sampled)
