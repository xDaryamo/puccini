import os.path, ctypes, c
from ruamel.yaml import YAML

yaml = YAML()

# See: https://medium.com/learning-the-go-programming-language/calling-go-functions-from-other-languages-4c7d8bcc69bf

library_path = os.path.join(os.path.dirname(__file__), '..', '..', 'dist', 'puccini-tosca.so')
lib = ctypes.cdll.LoadLibrary(library_path)

lib.Compile.argtypes = [ctypes.c_char_p]
lib.Compile.restype = ctypes.c_char_p

def compile(url):
    return yaml.load(lib.Compile(c.to_c_char_p(url)))
