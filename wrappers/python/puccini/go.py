
import ctypes

# See: https://medium.com/learning-the-go-programming-language/calling-go-functions-from-other-languages-4c7d8bcc69bf

class GoString(ctypes.Structure):
    _fields_ = (('p', ctypes.c_char_p), ('n', ctypes.c_longlong))

    def __init__(self, s):
        self.p = str.encode(s)
        self.n = len(s)

    def __str__(self):
        return self.p[:self.n].decode() if self.p else ''


def to_c_char_p(s):
    return ctypes.c_char_p(str.encode(s))
