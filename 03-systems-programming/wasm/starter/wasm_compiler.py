from error import LoxCompileError
from expr import Expr
from stmt import Stmt
from token_type import TokenType


class WasmCompiler(Expr.Visitor, Stmt.Visitor):
    def compile(self, statements):
        print("TODO: compile to WAT and print to stdout")
