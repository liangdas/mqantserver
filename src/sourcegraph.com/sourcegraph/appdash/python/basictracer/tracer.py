import time

import opentracing
from opentracing import Format, Tracer
from opentracing import UnsupportedFormatException
from .context import SpanContext
from .recorder import SpanRecorder, DefaultSampler
from .span import BasicSpan
from .util import generate_id


class BasicTracer(Tracer):

    def __init__(self, recorder=None, sampler=None):
        """Initialize a BasicTracer instance.

        Note that the returned BasicTracer has *no* propagators registered. The
        user should either call register_propagator() for each needed
        inject/extract format and/or the user can simply call
        register_required_propagators().

        The required formats are opt-in because of protobuf version conflicts
        with the binary carrier.
        """

        super(BasicTracer, self).__init__()
        self.recorder = NoopRecorder() if recorder is None else recorder
        self.sampler = DefaultSampler(1) if sampler is None else sampler
        self._propagators = {}

    def register_propagator(self, format, propagator):
        """Register a propagator with this BasicTracer.

        :param string format: a Format identifier like Format.TEXT_MAP
        :param Propagator propagator: a Propagator instance to handle
            inject/extract calls involving `format`
        """
        self._propagators[format] = propagator

    def register_required_propagators(self):
        from .text_propagator import TextPropagator
        from .binary_propagator import BinaryPropagator
        self.register_propagator(Format.TEXT_MAP, TextPropagator())
        self.register_propagator(Format.HTTP_HEADERS, TextPropagator())
        self.register_propagator(Format.BINARY, BinaryPropagator())

    def start_span(
            self,
            operation_name=None,
            child_of=None,
            references=None,
            tags=None,
            start_time=None):

        start_time = time.time() if start_time is None else start_time

        # See if we have a parent_ctx in `references`
        parent_ctx = None
        if child_of is not None:
            parent_ctx = (
                child_of if isinstance(child_of, opentracing.SpanContext)
                else child_of.context)
        elif references is not None and len(references) > 0:
            # TODO only the first reference is currently used
            parent_ctx = references[0].referenced_context

        # Assemble the child ctx
        ctx = SpanContext(span_id=generate_id())
        if parent_ctx is not None:
            if parent_ctx._baggage is not None:
                ctx._baggage = parent_ctx._baggage.copy()
            ctx.trace_id = parent_ctx.trace_id
            ctx.sampled = parent_ctx.sampled
        else:
            ctx.trace_id = generate_id()
            ctx.sampled = self.sampler.sampled(ctx.trace_id)

        # Tie it all together
        return BasicSpan(
            self,
            operation_name=operation_name,
            context=ctx,
            parent_id=(None if parent_ctx is None else parent_ctx.span_id),
            tags=tags,
            start_time=start_time)

    def inject(self, span_context, format, carrier):
        if format in self._propagators:
            self._propagators[format].inject(span_context, carrier)
        else:
            raise UnsupportedFormatException()

    def extract(self, format, carrier):
        if format in self._propagators:
            return self._propagators[format].extract(carrier)
        else:
            raise UnsupportedFormatException()

    def record(self, span):
        self.recorder.record_span(span)


class NoopRecorder(SpanRecorder):
    def record_span(self, span):
        pass
