from __future__ import absolute_import

from abc import ABCMeta, abstractmethod


class Propagator(object):
    __metaclass__ = ABCMeta

    @abstractmethod
    def inject(self, span_context, carrier):
        pass

    @abstractmethod
    def extract(self, carrier):
        pass
