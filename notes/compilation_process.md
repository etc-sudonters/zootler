# Compilation

OOTR's logic files are written in a subset of Python. OOTR is Python itself and
benefits from standard library offerings for parsing, transforming and
compiling Python code at run time[^fn1]. Zootler uses bespoke parsing,
optimization, and compilation.

## Parsing

Lexing and parsing are straightforward affairs, there's no back tracking
involved and all the rules produce booleans. Since this entire project
was inspired by Robert Nystrom's Crafting Interpreters it uses a Pratt parser.
The parser produces a concrete syntax tree that is immediately lowered into a
much sparser AST, where many syntax features are transformed into function
calls. This AST is then transformed in steps:

1. At/Here rules are yanked and replaced with token queries
2. Function calls are expanded in place where possible
3. Constant branches and comparisons are eliminated
4. Comparisons to setting values are rewritten into function calls
5. Standalone identifiers are expanded into function calls for checking
   settings or run time values

Where possible generic function calls are replaced with more specific function
calls. For example, the rule `skipped_trials[Forest]` is represented with a
`parser.Subscript` node in the CST but this is rewritten into an
`ast.Call{Callee: "load_setting_2", Args: ...}` and then into `ast.Call{Callee:
"is_trial_skipped", Args: []{Forest}}`.

## Compiling

After all rules, including extracted at/here, have been through this process
the AST is converted into 32bit linear IR -- "zasm". This produces a single
assembly that contains all the rules, and blocks for names, strings and
constants to produce a symbol table from. This assembly is intended to be
cached and reused for multiple generations.

The compiler combines a this assembly with settings to produce the final rule
set for the world. Settings are evaluated at compile and eliminated from the
final code, nested "and"/"or" and neighboring "has(TOK, 1)" are collapsed into
more efficient constructions.


[^fn1]: [Some notes here](./ootr-ast-rewriting.md)
