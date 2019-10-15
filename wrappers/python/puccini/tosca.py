import os.path, ctypes
from . import go
from ruamel.yaml import YAML

yaml = YAML()

library_path = os.path.join(os.path.dirname(__file__), 'libpuccini.so')
lib = ctypes.cdll.LoadLibrary(library_path)

lib.Compile.argtypes = [ctypes.c_char_p]
lib.Compile.restype = ctypes.c_char_p

def compile(url):
    return yaml.load(lib.Compile(go.to_c_char_p(url)))
