from error import LoxCompileError
from expr import Expr
from stmt import Stmt
from token_type import TokenType

# Everything is a signed i32
# Expressions add a single item to the stack
# Statements leave the stack unchanged

BINARY_OPS = {
    TokenType.BANG_EQUAL: "ne",
    TokenType.EQUAL_EQUAL: "eq",
    TokenType.GREATER: "gt_s",
    TokenType.GREATER_EQUAL: "ge_s",
    TokenType.LESS: "lt_s",
    TokenType.LESS_EQUAL: "le_s",
    TokenType.MINUS: "sub",
    TokenType.PLUS: "add",
    TokenType.SLASH: "div_s",
    TokenType.STAR: "mul",
}


class WasmCompiler(Expr.Visitor, Stmt.Visitor):
    def __init__(self):
        self.indent = 0

    def print(self, s, end=None):
        print(" " * self.indent + s, end=end)

    def compile(self, statements):
        self.print("(module")
        self.indent += 2
        self.print('(import "console" "log" (func $print (param i32)))')
        # TODO: Require that all top-level statements are function declarations
        for statement in statements:
            statement.accept(self)
        self.indent -= 2
        self.print(")")

    def visit_block_stmt(self, stmt):
        for statement in stmt.statements:
            statement.accept(self)

    def visit_expression_stmt(self, stmt):
        # TODO:
        # - Don't put extra junk on the stack after an expression statement
        # - Make sure for loops still work
        # - Make sure nested assignments work
        stmt.expression.accept(self)

    def visit_function_stmt(self, stmt):
        # TODO: Disallow nested functions
        self.print(f'(func ${stmt.name.lexeme} (export "{stmt.name.lexeme}")', end="")
        for param in stmt.params:
            # No indent
            print(f" (param ${param.lexeme} i32)", end="")
        # No indent
        print(f" (result i32)")
        self.indent += 2
        for statement in stmt.body:
            statement.accept(self)
        self.indent -= 2
        # TODO: Is there a way to avoid this hack for returns inside blocks?
        self.print("(unreachable)")
        self.print(")")

    def visit_if_stmt(self, stmt):
        """
        (block
          (block
            <condition>
            (i32.eqz)
            (br-if 0)
            <then-branch>
            (br 1)
          )
          <else-branch>
        )
        """
        if stmt.else_branch is not None:
            self.print("(block")
            self.indent += 2

        self.print("(block")
        self.indent += 2

        stmt.condition.accept(self)
        self.print("(i32.eqz)")
        self.print("(br_if 0)")
        stmt.then_branch.accept(self)

        if stmt.else_branch is not None:
            self.print("(br 1)")

        self.indent -= 2
        self.print(")")

        if stmt.else_branch is not None:
            stmt.else_branch.accept(self)
            self.indent -= 2
            self.print(")")

    def visit_print_stmt(self, stmt):
        stmt.expression.accept(self)
        self.print(f"(call $print)")

    def visit_return_stmt(self, stmt):
        # TODO: Require that return statements have values
        stmt.value.accept(self)
        self.print(f"(return)")

    def visit_var_stmt(self, stmt):
        # TODO:
        # - Require that all declarations are at top of function
        # - Disallow initializers (declaration only)
        self.print(f"(local ${stmt.name.lexeme} i32)")

    def visit_while_stmt(self, stmt):
        """
        (block
          (loop
            <condition>
            (i32.eqz)
            (br_if 1)
            <body>
            (br 0)
          )
        )
        """
        self.print("(block")
        self.indent += 2
        self.print("(loop")
        self.indent += 2
        stmt.condition.accept(self)
        self.print("(i32.eqz)")
        self.print("(br_if 1)")
        stmt.body.accept(self)
        self.print("(br 0)")
        self.indent -= 2
        self.print(")")
        self.indent -= 2
        self.print(")")

    def visit_assign_expr(self, expr):
        expr.value.accept(self)
        self.print(f"(local.set ${expr.name.lexeme})")

    def visit_binary_expr(self, expr):
        if expr.operator.type not in BINARY_OPS:
            raise LoxCompileError(f"Unknown binary operator {expr.token}")
        expr.left.accept(self)
        expr.right.accept(self)
        self.print(f"(i32.{BINARY_OPS[expr.operator.type]})")

    def visit_call_expr(self, expr):
        for argument in expr.arguments:
            argument.accept(self)
        # TODO: Require that the callee is an identifier
        self.print(f"(call ${expr.callee.name.lexeme})")

    def visit_grouping_expr(self, expr):
        expr.expression.accept(self)

    def visit_literal_expr(self, expr):
        self.print(f"(i32.const {int(expr.value)})")

    def visit_logical_expr(self, expr):
        if expr.operator.type == TokenType.AND:
            op = "and"
        elif expr.operator.type == TokenType.OR:
            op = "or"
        else:
            raise LoxCompileError(f"Unknown logical operator {expr.token}")
        expr.left.accept(self)
        expr.right.accept(self)
        self.print(f"(i32.{op})")

    def visit_unary_expr(self, expr):
        if expr.operator.type == TokenType.BANG:
            expr.right.accept(self)
            self.print(f"(i32.eqz)")
        elif expr.operator.type == TokenType.MINUS:
            self.print(f"(i32.const 0)")
            expr.right.accept(self)
            self.print(f"(i32.sub)")

    def visit_variable_expr(self, expr):
        self.print(f"(local.get ${expr.name.lexeme})")
