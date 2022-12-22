
import os.path, ctypes, ard
from . import go

library_path = os.path.join(os.path.dirname(__file__), 'libpuccini.so')
library = ctypes.cdll.LoadLibrary(library_path)

library.Compile.argtypes = (ctypes.c_char_p,)
library.Compile.restype = ctypes.c_char_p


class Problems(Exception):
    def __init__(self, problems):
        self.message = 'problems'
        self.problems = problems


def compile(url, inputs=None, quirks=None):
    inputs = ard.encode(inputs or {})
    quirks = ard.encode(quirks or [])
    result = ard.read(library.Compile(go.to_c_char_p(url), go.to_c_char_p(inputs), go.to_c_char_p(quirks)))
    if 'problems' in result:
        raise Problems(result['problems'])
    elif 'error' in result:
        raise Exception(result['error'])
    return result['clout']
