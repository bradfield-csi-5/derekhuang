from token_type import TokenType


class LoxCompileError(Exception):
    pass


class Handler:
    had_error = False

    @classmethod
    def error(cls, line, message):
        cls.report(line, "", message)

    @classmethod
    def report(cls, line, where, message):
        print(f"[line {line}] Error{where}: {message}")
        cls.had_error = True

    @classmethod
    def parse_error(cls, token, message):
        if token.type == TokenType.EOF:
            cls.report(token.line, " at end", message)
        else:
            cls.report(token.line, f" at {token.lexeme}", message)
