from misc.factory import Factory
from abc import ABCMeta, abstractmethod

class BuilderGen(metaclass=ABCMeta):
    """
    Builder interface base class.
    """

    def __str__(self):
        return self.__class__.__name__

    @abstractmethod
    def acq(self):
       """document's data acquisition"""

    @abstractmethod
    def write(self):
        """writes the document"""

    @abstractmethod
    def rel(self):
        """release resources previously gotten"""


class DocBuilderError(Exception):

    def __init__(self, message=None):
        self.message = message

    def __str__(self):
        return self.message


class DocBuilderStepError(DocBuilderError):

    def __init__(self, message=None, step_data=None):
        self.step_data = step_data
        super().__init__(message=message)


class PipeLine(Factory):

    def __init__(self, *args):
        for t in args:
            self.subscribe(*t)

    def __call__(self, barchetype, **kwargs):
        """runs pipeline to create a document"""

        ic = self.inceptors.get(barchetype, None)

        if ic is None:
            msg = "{0} is an unsupported builder".format(barchetype)
            raise DocBuilderError(msg)

        if not issubclass(ic, BuilderGen):
            msg = "{0} is an unsupported class specification".format(builder)
            raise DocBuilderError(msg)

        builder = ic(**kwargs)

        for s in ("acq", "write", "rel"):
            try:
                step = getattr(builder, s)
                step()
            except DocBuilderStepError:
                raise
