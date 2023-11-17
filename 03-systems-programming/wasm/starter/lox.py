import sys
from parser import Parser

from error import Handler
from scanner import Scanner
from wasm_compiler import WasmCompiler

compiler = WasmCompiler()


def run_file(path):
    with open(path) as f:
        run(f.read())

    if Handler.had_error:
        sys.exit(65)


def run_prompt():
    while True:
        try:
            line = input("> ")
            run(line)
            Handler.had_error = False  # Recover from user mistakes
        except EOFError:
            break


def run(source):
    scanner = Scanner(source)
    tokens = scanner.scan_tokens()
    parser = Parser(tokens)
    statements = parser.parse()

    if Handler.had_error:
        return

    compiler.compile(statements)


if __name__ == "__main__":
    if len(sys.argv) > 2:
        print(f"Usage: {sys.argv[0]} [script]")
        sys.exit(1)
    elif len(sys.argv) == 2:
        run_file(sys.argv[1])
    else:
        run_prompt()
