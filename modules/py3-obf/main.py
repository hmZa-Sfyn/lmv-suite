#!/usr/bin/env python3
import ast
import random
import string
import sys
import os
import base64
import re
import zlib

class UltraObfuscator(ast.NodeTransformer):
    def __init__(self, level=1):
        self.level = level
        self.var_map = {}
        self.func_map = {}
        self.used_names = set(list(globals().keys()) + ['base64', 're', 'random', 'zlib'])

    def generate_random_name(self, length=25):  # Even longer for ultra confusion
        while True:
            name = '_' + ''.join(random.choices(string.ascii_letters + string.digits, k=length))
            if name not in self.used_names:
                self.used_names.add(name)
                return name

    def visit_Name(self, node):
        if isinstance(node.ctx, ast.Store):
            if node.id not in self.var_map:
                self.var_map[node.id] = self.generate_random_name()
            node.id = self.var_map[node.id]
        elif isinstance(node.ctx, ast.Load):
            node.id = self.var_map.get(node.id, node.id)
        return node

    def visit_FunctionDef(self, node):
        if node.name not in self.func_map:
            self.func_map[node.name] = self.generate_random_name()
        node.name = self.func_map[node.name]
        for arg in node.args.args:
            if arg.arg not in self.var_map:
                self.var_map[arg.arg] = self.generate_random_name()
            arg.arg = self.var_map[arg.arg]
        self.generic_visit(node)
        return node

    def visit_Call(self, node):
        if isinstance(node.func, ast.Name):
            node.func.id = self.func_map.get(node.func.id, node.func.id)
        self.generic_visit(node)
        return node

    def encode_string(self, s):
        if not s:
            return ast.Constant(value='')
        encoded = base64.b64encode(s.encode('utf-8')).decode('utf-8')
        b64decode = ast.Attribute(
            value=ast.Name(id='base64', ctx=ast.Load()),
            attr='b64decode',
            ctx=ast.Load()
        )
        call_decode = ast.Call(
            func=b64decode,
            args=[ast.Constant(value=encoded)],
            keywords=[]
        )
        return ast.Call(
            func=ast.Attribute(value=call_decode, attr='decode', ctx=ast.Load()),
            args=[ast.Constant(value='utf-8')],
            keywords=[]
        )

    def visit_Constant(self, node):
        if self.level >= 2 and isinstance(node.value, str) and len(node.value) > 3:
            return self.encode_string(node.value)
        return node

    def visit_JoinedStr(self, node):
        if self.level < 2:
            return node
        new_values = []
        for val in node.values:
            if isinstance(val, ast.Constant) and isinstance(val.value, str):
                new_values.append(self.encode_string(val.value))
            elif isinstance(val, ast.FormattedValue):
                val.value = self.visit(val.value)
                new_values.append(val)
            else:
                new_values.append(val)
        return ast.JoinedStr(values=new_values)

    def add_junk(self, body):
        if self.level >= 3 and len(body) > 1:
            num_junk = random.randint(10, 20)  # Ultra junk
            for _ in range(num_junk):
                junk_options = [
                    # Empty print
                    ast.Expr(value=ast.Call(
                        func=ast.Name(id='print', ctx=ast.Load()),
                        args=[ast.Constant(value='')],
                        keywords=[]
                    )),
                    # Dummy assign
                    ast.Assign(
                        targets=[ast.Name(id=self.generate_random_name(), ctx=ast.Store())],
                        value=ast.BinOp(
                            left=ast.Constant(value=random.randint(1,10000)),
                            op=random.choice([ast.Add(), ast.Sub(), ast.Mult(), ast.Div()]),
                            right=ast.Constant(value=random.randint(1,10000))
                        )
                    ),
                    # False if
                    ast.If(
                        test=ast.Compare(
                            left=ast.Call(func=ast.Attribute(value=ast.Name(id='random', ctx=ast.Load()), attr='random', ctx=ast.Load()), args=[], keywords=[]),
                            ops=[ast.Gt()],
                            comparators=[ast.Constant(value=2)]
                        ),
                        body=[ast.Pass()],
                        orelse=[]
                    ),
                    # Regex compile
                    ast.Expr(value=ast.Call(
                        func=ast.Attribute(value=ast.Name(id='re', ctx=ast.Load()), attr='compile', ctx=ast.Load()),
                        args=[ast.Constant(value=''.join(random.choices(string.ascii_letters, k=10)))],
                        keywords=[]
                    )),
                    # Junk function def
                    ast.FunctionDef(
                        name=self.generate_random_name(),
                        args=ast.arguments(posonlyargs=[], args=[], vararg=None, kwonlyargs=[], kw_defaults=[], kwarg=None, defaults=[]),
                        body=[ast.Pass(), ast.Expr(value=ast.Constant(value=random.randint(1,100)))],
                        decorator_list=[],
                        returns=None
                    )
                ]
                body.insert(random.randint(0, len(body)-1), random.choice(junk_options))
        return body

    def visit_Module(self, node):
        node.body = self.add_junk(node.body)
        self.generic_visit(node)
        return node

def minify_internal(code):
    lines = code.split('\n')
    minified_lines = []
    for line in lines:
        stripped = line.strip()
        if stripped and not stripped.startswith('#'):  # Skip empty and comments
            indent = len(line) - len(line.lstrip())
            minified_lines.append(' ' * indent + stripped)
    return '\n'.join(minified_lines)

def obfuscate_code(source_code, level):
    tree = ast.parse(source_code)
    obf = UltraObfuscator(level=level)
    new_tree = obf.visit(tree)
    ast.fix_missing_locations(new_tree)
    obfuscated = ast.unparse(new_tree)
    minified = minify_internal(obfuscated)
    return minified

def main():
    target = os.getenv('ARG_TARGET', '').strip()
    output = os.getenv('ARG_OUTPUT', 'obfuscated.py').strip()
    try:
        level = int(os.getenv('ARG_LEVEL', '1'))
        if level not in [1, 2, 3]:
            level = 1
    except:
        level = 1

    if not target:
        print('[!] Target file is required')
        sys.exit(1)

    if not os.path.isfile(target):
        print('[!] Input file does not exist or is not a file')
        sys.exit(1)

    try:
        with open(target, 'r', encoding='utf-8') as f:
            code = f.read()

        obfuscated_code = obfuscate_code(code, level)

        # Always compress for no visible newlines/spaces
        compressed = zlib.compress(obfuscated_code.encode('utf-8'), level=9)
        encoded = base64.b64encode(compressed).decode('utf-8')

        # Wrapper: one-liner runnable
        wrapper = f"import zlib,base64;exec(zlib.decompress(base64.b64decode('{encoded}')).decode('utf-8'))"

        if level >= 3:
            wrapper = f"import zlib,base64,re,random;exec(zlib.decompress(base64.b64decode('{encoded}')).decode('utf-8'))"

        with open(output, 'w', encoding='utf-8') as f:
            f.write(wrapper)

        print(f'[+] UltraSuperCool obfuscation successful (level {level})')
        print(f'[+] Obfuscated file saved to: {output} (runnable one-liner!)')

    except Exception as e:
        print(f'[!] Error during obfuscation: {str(e)}')
        sys.exit(1)

if __name__ == '__main__':
    main()