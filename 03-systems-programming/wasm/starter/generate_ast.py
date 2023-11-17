import sys


def define_ast(output_dir, base_name, types):
    path = f"{output_dir}/{base_name.lower()}.py"
    type_map = {}
    for type in types.split("\n"):
        if not type:
            continue
        parts = type.split(":")
        class_name = parts[0].strip()
        fields = parts[1].strip().split()
        type_map[class_name] = fields

    with open(path, "w") as f:
        f.write(f"class {base_name}:\n")
        f.write("    class Visitor:\n")
        for class_name in type_map:
            f.write(
                f"        def visit_{class_name.lower()}_{base_name.lower()}(self, {base_name.lower()}):\n"
            )
            f.write(f"            pass\n\n")
        for class_name, fields in type_map.items():
            define_type(f, base_name, class_name, fields)


def define_type(f, base_name, class_name, fields):
    f.write(
        f"""    class {class_name}:
        def __init__(self, {', '.join(fields)}):
"""
    )
    for field in fields:
        f.write(f"            self.{field} = {field}\n")
    f.write(
        f"""
        def accept(self, visitor):
            return visitor.visit_{class_name.lower()}_{base_name.lower()}(self)

"""
    )


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <output directory>")
        sys.exit(64)

    output_dir = sys.argv[1]
    define_ast(
        output_dir,
        "Expr",
        """
Assign   : name value
Binary   : left operator right
Call     : callee paren arguments
Grouping : expression
Literal  : value
Logical  : left operator right
Unary    : operator right
Variable : name
""",
    )

    define_ast(
        output_dir,
        "Stmt",
        """
Block      : statements
Expression : expression
Function   : name params body
If         : condition then_branch else_branch
Print      : expression
Return     : keyword value
Var        : name initializer
While      : condition body
""",
    )
