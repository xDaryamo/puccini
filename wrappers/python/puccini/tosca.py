import os.path, ctypes
from . import go
from ruamel.yaml import YAML

yaml = YAML()

library_path = os.path.join(os.path.dirname(__file__), 'libpuccini.so')
library = ctypes.cdll.LoadLibrary(library_path)

library.Compile.argtypes = [ctypes.c_char_p]
library.Compile.restype = ctypes.c_char_p

def compile(url):
    return yaml.load(library.Compile(go.to_c_char_p(url)))
