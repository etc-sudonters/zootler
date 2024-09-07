# Macros && Funcs

Zootler has two kinds of macros:

1. Scripted macros, aka `helpers.json`, which are C preprocessor style lexical
   macros. The token body of the macros is copy and pasted in place of the
   macro's invocation, and if there are parameters the literal tokens of the
   arguments are pasted over the macro's parameters.
2. Built in, i.e. `at` and `here`, which are also lexical in nature but run
   arbitrary code when expanded rather than copying token streams.

## Prior Art

I looked at three different implementations of macros for inspiration:

1. OOTR's implementation
2. cpp documentation
3. rustc documentation

### OOTR

OOTR also has a sort of multitiered approach to its rule parsing and compiling.
Scripted macros have parameters textually replaced with regexes and then that
resulting string is sent into the AST engine, recursively expanding macros.
There are also several macros that are procedurally expanded, primarily `at`
and `here` but there are a few others such as the time of day checks. 

### cpp

The C_ _P_re_p_rocessor[^fn10] is a _lexical_ macro system only[^fn1]. It has specific
rules around what, how, and when macros are expanded. When the preprocessor
encounters an identifier that corresponds to an object-like macro, or a call
that corresponds to a function-like macro then it replaces the tokens _of the
entire invocation_ with whatever the body of the macros is. After this token
replacement, these new tokens are then parsed and possibly expanded. However,
cpp will not recursively expand a macro -- that is a macro that attempts to
expand itself. However, macro arguments are expanded before the macro call
itself so stacked calls to the same macro -- `MIN(MIN(a,b), MAX(c,d))` would
expand `MIN` twice as expected.

```c
//object like 
#define FOO 1
// expands to literally BAR, will not recursively expand
#define BAR BAR

//function like
#define MIN(a, b) (a < b) ? (a) : (b)

// (1 < 2) ? (1) : (2)
MIN(1,2)
// expands inside out so ARG1 and ARG2 are expanded first, and then those
// expansions are pasted into the outer min, the full expansion is quite big
MIN( MIN(1,2), MIN(3,4) )
```

### rust/rustc

rustc also has a multitiered approach to macros. I primarily took notes on
procedural macros, or procmacros, and macros by example/`macro_rules!`[^fn2].
Both options are hygienic -- the macro generated code is inserted in a way that
doesn't clobber the surrounding environment. `macro_rules!` has its own parser
machinery and the macros are lexical in nature -- the macro is handed some
typed tokens and then generates valid rust syntax code with it.[^fn3]
Procmacros are external code that is compiled and invoked by the compiler and
is able to spit out either tokens or ast -- at least definitely tokens. These
macros include things like custom derives. Most of the notes I have taken are
around its work queue for macro expansion and its specific fragment integration
rules that are interesting but not particularly relevant here.[^fn9]

## Implementation

Before diving into the macro system implementation, a quick summary of the
process without macros involved.

### Parsing, Compiling, and some run time

Since the initial inspiration for this project was Robert Nystrom's Crafting
Interpreters the parser is a Pratt style parser. The parser itself accepts a
`TokenStream` and a `Grammar[T]` and then does this kind of recursive,
iterative dance to transform the sequential, flat token space into deeply and
properly nested parse trees. Parsing a string rule creates a
`peruse.StringLexer` and then feeds the resulting tokens into a
`peruse.Parser[T]`[^fn5]. The resulting parse tree has enough fidelity to
source that it could be written directly back into its originating source.[^fn6]

The parse tree is then lowered into AST. The AST has fewer nodes, and ergo
fewer features. Mostly these concern the several ways settings and token
quantities can be queried by the rules. The AST lower will attempt to eliminate
constant operations such as compares to self and redundant arms for `and`/`or`. 
Since we eliminate string literal aliases for tokens in this pass, the lowering
also begins the process of interning strings, names and constants. 

1. `(Token, Qty)` are explicitly rewritten as `has(Token, Qty)` 
2. `setting_block[sub_setting]` and `sub_setting in setting_block` are
   rewritten to function calls
3. Standalone identifers that correspond to setting names or trick names are
   also rewritten to calls
4. `==`/`!=` attempt to eliminate themselves to constants, `and`/`or` attempt
   to eliminate arms
5. Standalone identifiers that correspond to "tokens" are rewritten to
   `has(Token, 1)` calls

This AST is then written to an assembly consisting of:

1. String table
2. Constant table
3. Name table
4. Instruction table

The constants table holds actual values as 64bit packed quiet NANs.[^fn8] The
strings table holds all the string data for the assembly in a single []uint8
and provides 3 byte offset+length handles. The names table also these 3byte
pointers into the strings but this string is then used to query the run time.
The instruction table stores rules in a single []uint32 and provides string
aliases to actual Go slices of its contents. Instructions are 32bit linear
"assembly" -- zasm -- with 8bit operations and 24bit arbitrary payloads[^fn7]. 

This assembly isn't intended to be executed directly. Instead the assembly is
meant to be a cache of the mostly compiled rules. The last mile happens when
the assembly and specific settings are combined by the compiler to generate the
final tables and op codes for the run time VM. The compiler outputs 8bit
instructions with variable length encoding -- VLE -- to facilitate very common
operations like checking token quantities. Instead of being encoding as:

1. OP_LOAD_IDENT
2. OP_LOAD_CONST
3. OP_CALL_2 [u16lo, u16hi]

The compiler emits: `OP_CHK_QTY [u16lo, u16hi, u8]` which is actually a 1:1
translation of what the assembler emits in 1 instruction, just unpacked into
4 uint8s.

### Macros

Macros are inserted at the lexical level, which makes sense they're lexical
macros.[^fn11] The short of it is that specific syntax parsers are "annointed"
and can divert the flow of tokens into the macro infrastructure. However, when
it can't, regular parsing happens. 

Specifically, `RulesParser` overwrites the rules for parsing identities and
calls to first check if the token, or sequence of tokens, is eligible for
expansion. Eligibility is determined by a single rule: if its a macro currently
being expanded, then it isn't eligible for expansion. This prevents recursive
expansions of macros, and is similar to cpp's eligibility rules in effect and
implementation.

To handle expansion, the macro invocation and its arguments are passed to an
`Expander` which is responsible for outputting a parser tree fragment. For
script macros the expansion is yield macro body tokens until a parameter stub
token is encountered, in which case yield all the tokens for the relevant
argument and then go back to the macro's body. "Built in" macros -- maybe
procmacros is a better name -- don't have a macro body _per se_. Instead these
are typically executing arbitrary code to effect other parts of the application
and then outputting a fragment that reflects how the rule interacts with that.

Here, expansion is handled recursively. When the parser rule diverts into macro
expansion, any parser driven expansions of code use the same parser that
invoked the macro expansion. The parser handles this by managing both a stack
of token streams and a stack of macro expansions. Expansion begins with pushing
a new expansion context, copying the expanding macro into it, copying the
tokens AND AST for the macro arguments into it and marking this macro as
ineligible for expansion.[^fn11] At the end of expansion, the context is
destructed and the macro is marked eligible for expansion again.

---
[ln1]: github.com/etc-sudonters/substrate/peruse 
[ln2]: https://rustc-dev-guide.rust-lang.org/overview.html
[ln3]: https://rustc-dev-guide.rust-lang.org/macro-expansion.html
[ln4]: https://gcc.gnu.org/onlinedocs/cpp/Macros.html
---
[^fn1]: as far as I can tell, I'm not really a C person

[^fn2]: chanting "macro rules" as my car spins out on a banana peel and off a
    cliff

[^fn3]: I'm also not really a rust person, but this has been my experience at
    least

[^fn5]: `peruse` comes from [etc-sudonters/substrate][ln1] which is where stuff
    goes to prevent too much fiddling with. I did have to fiddle some to make
    macros work but primarily it was abstracting StringLexer into TokenStream.

[^fn6]: I usually find `s-expr` format to be friendlier, especially with some
    kind of visual paren/depth matcher. In any case, it is possible to a
    literal translation from rule source to some other target.


[^fn7]: The astute will note that [3]uint8 is 24bits which is the size of the
    string and name pointers. The final result is intentional, but it took some
    playing around with to realize that. 

[^fn8]: the quick and dirty is that there are 9 types encodable with this
    scheme. Either it's not NAN, ergo a float64 value. Or it is a NAN with bits
    63, 49 and 48 set respectively in the table below. This leaves a whooping 6
    bytes that can be freely used. The assembly pays a very inflated cost for
    small values like individual booleans but it is a consistent way of dealing
    with values and eliminates at least portions of the string and name tables.
    ```
    63
    |49
    ||48
    |||
    mask | type
    000  | Bool
    001  | location literal -- hypothetical
    010  | I48 (negative) or U47, op dependent
    011  | short string  -- up to 6 ascii chars, otherwise nul terminated
    100  | identity
    101  | token literal
    110  | I48
    111  | [6]uint8 -- payload is operation defined
    ```

[^fn9]: The [overview][ln2] and [macro expansion][ln3] pages on the rustc
    compiler documentation where my primary sources. 

[^fn10]: The [gnu cpp docs][ln4] were my primary reference here

[^fn11]: This only became obvious over time and trying many things that didn't
    work.


[^fn11]: The implementation also provide an opaque deferable to destruct the context.
