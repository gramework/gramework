#include "textflag.h"

TEXT ·call(SB), NOSPLIT, $0-32
	JMP runtime·reflectcall(SB)
