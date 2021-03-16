import os.path, ctypes
from . import go
from io import StringIO
from contextlib import closing
from ruamel.yaml import YAML

yaml = YAML()

library_path = os.path.join(os.path.dirname(__file__), 'libpuccini.so')
library = ctypes.cdll.LoadLibrary(library_path)

library.Compile.argtypes = [ctypes.c_char_p]
library.Compile.restype = ctypes.c_char_p

class Problems(Exception):
    def __init__(self, problems):
        self.message = 'problems'
        self.problems = problems

def compile(url, inputs={}):
    inputs = _yaml_dumps(inputs)
    result = yaml.load(library.Compile(go.to_c_char_p(url), go.to_c_char_p(inputs)))
    if 'problems' in result:
        raise Problems(result['problems'])
    elif 'error' in result:
        raise Exception(result['error'])
    return result['clout']

def _yaml_dumps(data):
    with closing(StringIO()) as s:
        yaml.dump(data, s)
        return s.getvalue()
