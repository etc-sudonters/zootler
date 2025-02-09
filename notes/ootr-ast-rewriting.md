# OOTR AST Rewriting && Disassembly

OOTR has the benefit of being written in Python and additionally using Python
syntax for its logic rules, which means it can take advantage of Python's
ability to parse and compile Python code at run time. Python also offers
standard library tools for manipulating it's intermediate abstract syntax tree
-- AST -- representation and for inspecting the op code of compiled Python
objects.

## Abstract Syntax Trees

AST is an intermediate representation of code, and despite its graphical nature
typically has a relatively high fidelity to source. Individual pieces of source
code are expanded into much larger, nested structures that recursively
_represent_ the relevant source. ASTs dispense with redundant information like
parens because precedence is encoded into the tree hierarchy, so they aren't a
fully recreation of source.[^fn2]

```
 (2+2)*4  |     2+2*4   
----------|--------------
      *   |       +
    /  \  |     /  \
   +    4 |    *    2
 /  \     |  /  \
2    2    | 2    4
```


## Python AST Primer

I highly recommend becoming familiar with the Python AST module documentation.
It's been several years since I'd done anything interesting with this module
and a refresher was worthwhile. The short of it is that once you get your hands
on some Python AST you can use either a `NodeVisitor` to walk the tree or
`NodeTransformer` to rewrite the AST.

There's a few more notes like using `fix_missing_locations`, using `compile`
and `eval` or `exec` to transform the AST into python code objects and execute
them. Another module to be aware of is the `dis` module which disassembles
python code objects into op codes[^fn1].

Interacting with Python's AST involves a lot of pomp and circumstance that
often overshadows what is actually happening.

```python
ast.Attribute(
    value=ast.Attribute(
        value=ast.Attribute( 
            value=ast.Name(
                id='state',
                ctx=ast.Load()
            ),
            attr='world',
            ctx=ast.Load()
        ), 
        attr='settings',
        ctx=ast.Load()
    ), 
    attr=child.id,
    ctx=ast.Load()
)
```


## Rule AST Transformer

This class handles the entire process of transforming some string of syntax
into a callable rule function. Macros are expanded, identities are established,
optimizations are applied, and the rule is stuffed inside the AST of a lambda
which is then compiled and evaluated -- producing the callable lambda rule.

Everything that follows is an organized summary and not 100% accurate.

## Tranformations

### visit_Name

1. If the name is a method on the transformer, call that method with the
   current node and use the output as the replacement AST. This handles `at`,
   `here`, and at time of day checks.
2. If the name is a "rule_alias" -- it originates from `helpers.json`:
    - If it does not take 0 arguments, throw
    - Otherwise, expand its body into AST, using this instance to transform it
3. If the name refers to an "escaped_item", replace it with `has(name, 1)`
4. If the name is a value on world or settings, resolve to that value
5. If the name resolves to a method on "State" replace with a call to that method
6. If the name matches a particular event regex, then ensure it exists as an
   event and replace with a call to `has(name, 1)`
7. Throw -- this name isn't known

### visit_Str

This replaces string literals of items with `has(item, 1)`. This is called
by the `visit_Constant` method which Python stdlib calls.

### visit_Tuple

Tuples are replaced with `has(item[0], item[1])` when `item[0]` resolves to an
item and `item[1]` resolves to a number. Otherwise this method throws.

### visit_Subscript

Indexes into nested setting block.

### visit_Call

1. If the caller isn't an identifier, return -- this doesn't come up in the
   logic files
2. If the caller is a method on the transformer, call it with the passed node
   and return its result as the replacement AST
3. If the caller is a "rule_alias" then:
    1. Check that the number of passed arguments is correct
    2. Evaluate the arguments into either identifiers or constants
    3. Regex the source of the alias, replacing parameter names with the
       evaluated arguments
    4. Parse the resulting source, using this transformer to transform it
4. Otherwise, resolve arguments and produce AST that calls a matching method
   name on State

### visit_Compare

This is primarily `==` and `!=`; however there is a single `<` used for
checking if child chicken collection is possible. Python allows chaining
comparators and the AST node looks _very unusual_ without that information.
This method eliminates redundant and constant comparisons. In particular
`can_use` is mentioned, which benefits a lot from branch elimination.

### visit_UnaryOp

Only `not` is used in the logic files which inverts a truthiness value of any
kind in Python. This attempts to statically evaluate the result.

### visit_BinOp

This handles any `setting_name in setting_block` look ups.

### visit_BoolOp

Handles `and`/`or` operations and eliminates redundant branches as much as
possible, and attempts to eliminate multiple branches at once.

## Additional methods

### replace_subrule

Expands `at` and `here` macros by creating a new event  in the target location
-- either the one specified by `at` or the location currently being processed
-- and replacing the macro with a `has(event name, 1)` call. These event rules
are added to a queue that is handled after initial rule parsing has finished.
The generated events for the same target location have sequential names.

### create_delayed_rules

Handles the processing of the created subrules.

### make_access_rule

Handles compiling the ast and creating a callable from the code object.

## Dissambled Rules

Examining the op codes generated by the entire process is also pretty
interesting. It also gives a target to match in behavior. libzootr doesn't
intend to match Python semantics exactly but it should produce comparable byte
code -- albeit in its own instruction set.

```python
# ~/source/external/OoT-Randomizer
# $ ipython
from Settings import Settings
from Main import resolve_settings, build_world_graphs
from io import StringIO
import dis
from random import shuffle
# disable all the rom stuff
settings = Settings({})
settings.create_uncompressed_rom = settings.create_compressed_rom = settings.create_wad_file = False
settings.patch_without_output = settings.create_patch_file = False
# otherwise randomizer get mads
settings.create_spoiler = True
resolve_settings(settings)
worlds = build_world_graphs(settings)
world = worlds[0]
# pretty cool this works
shuffle(locations := world.get_locations())

get_location = iter(locations).__next__

def dis_rule(loc):
    print(f"{loc.parent_region} -> {loc.name}: {loc.rule_string}\n")
    for rule in loc.access_rules:
        dis.dis(rule)
        print("")

```

The very first location popped for me is one of the subrules!
```
Bottom of the Well -> Bottom of the Well Subrule 2: None

  1           0 LOAD_FAST                1 (age)
              2 LOAD_CONST               1 ('adult')
              4 COMPARE_OP               2 (==)
              6 RETURN_VALUE
```

An example of how the transformer folds constants:

```
Haunted Wasteland -> Wasteland Crate After Quicksand 1: can_break_crate

  1           0 LOAD_CONST               1 (True)
              2 RETURN_VALUE

GV Upper Stream -> GV Crate Near Cow: is_child and can_break_crate

  1           0 LOAD_FAST                1 (age)
              2 LOAD_CONST               1 ('child')
              4 COMPARE_OP               2 (==)
              6 RETURN_VALUE
```

A call into the state object:

```
HC Garden -> HC Zeldas Courtyard Mario Wonderitem: can_use(Slingshot)

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has_all_of)
              4 LOAD_GLOBAL              1 (Slingshot)
              6 BUILD_TUPLE              1
              8 CALL_METHOD              1
             10 JUMP_IF_FALSE_OR_POP     9 (to 18)
             12 LOAD_FAST                1 (age)
             14 LOAD_CONST               1 ('child')
             16 COMPARE_OP               2 (==)
        >>   18 RETURN_VALUE
```

I've included some longer dissaemblies below, but looking at the op codes isn't
the only interesting thing available to us:

```python
dis.show_code(get_location().access_rules[0])
```

Lovely spirit temple logic yum yum yum. The file name is generated by
`make_access_rule` by string dumping the AST, which is very interesting side
channel information. The AST is also deceptively compressed in it. But we're
able to find information like what tokens the rule interacts with.

```
Name:              <lambda>
Filename:          <Spirit Temple Central Chamber Pot 1: BoolOp(Or(), [Call(Attribute(Name('state', Load()), 'has_any_of', Load()), [Tuple([Name('Bombchus_10'
, Load()), Name('Bombchus', Load()), Name('Bomb_Bag', Load()), Name('Bombchus_5', Load()), Name('Bombchus_20', Load())], Load())], []), Call(Attribute(Name('s
tate', Load()), 'has', Load()), [Name('Small_Key_Spirit_Temple', Load()), Constant(3)], []), Compare(Constant('Spirit Temple'), [In()], [List([], Load())]), C
all(Attribute(Name('state', Load()), 'has', Load()), [Name('Small_Key_Spirit_Temple', Load()), Constant(2)], [])])>
Argument count:    1
Positional-only arguments: 0
Kw-only arguments: 3
Number of locals:  4
Stack size:        7
Flags:             OPTIMIZED, NEWLOCALS, NOFREE, 0x1000000
Constants:
   0: None
   1: 3
   2: 'Spirit Temple'
   3: ()
   4: 2
Names:
   0: has_any_of
   1: Bombchus_10
   2: Bombchus
   3: Bomb_Bag
   4: Bombchus_5
   5: Bombchus_20
   6: has
   7: Small_Key_Spirit_Temple
Variable names:
   0: state
   1: age
   2: spot
   3: tod
```

Here's an example of a fictional location the randomizer uses:
```
Name:              <lambda>
Filename:          <Root Exits -> Child Spawn: Compare(Name('age', Load()), [Eq()], [Constant('child')])>
Argument count:    1
Positional-only arguments: 0
Kw-only arguments: 3
Number of locals:  4
Stack size:        2
Flags:             OPTIMIZED, NEWLOCALS, NOFREE, 0x1000000
Constants:
   0: None
   1: 'child'
Variable names:
   0: state
   1: age
   2: spot
   3: tod
```

### at/here expansions

```
Market Mask Shop Storefront -> Mask of Truth Access from Market Mask Shop Storefront: complete_mask_quest or (at('Kakariko Village', is_child and Keaton_Mask) and at('Lost Woods', is_child and can_play(Sarias_Song) and Skull_Mask) and at('Graveyard', is_child and at_day and Spooky_Mask) and at('Hyrule Field', is_child and has_all_stones and Bunny_Hood))


  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has_all_of)
              4 LOAD_GLOBAL              1 (Kakariko_Village_Subrule_1)
              6 LOAD_GLOBAL              2 (Hyrule_Field_Subrule_2)
              8 LOAD_GLOBAL              3 (Lost_Woods_Subrule_3)
             10 LOAD_GLOBAL              4 (Graveyard_Subrule_1)
             12 BUILD_TUPLE              4
             14 CALL_METHOD              1
             16 RETURN_VALUE

Name:              <lambda>
Filename:          <Mask of Truth Access from Market Mask Shop Storefront: Call(Attribute(Name('state', Load()), 'has_all_of', Load()), [Tuple([Name('Kakariko
_Village_Subrule_1', Load()), Name('Hyrule_Field_Subrule_2', Load()), Name('Lost_Woods_Subrule_3', Load()), Name('Graveyard_Subrule_1', Load())], Load())], []
)>
Argument count:    1
Positional-only arguments: 0
Kw-only arguments: 3
Number of locals:  4
Stack size:        6
Flags:             OPTIMIZED, NEWLOCALS, NOFREE, 0x1000000
Constants:
   0: None
Names:
   0: has_all_of
   1: Kakariko_Village_Subrule_1
   2: Hyrule_Field_Subrule_2
   3: Lost_Woods_Subrule_3
   4: Graveyard_Subrule_1
Variable names:
   0: state
   1: age
   2: spot
   3: tod

```


## More disassembles

### Multi rule locations

Let's find some locations with multiple rules attached to them:

```python
multi_rule_locations = [l for l in locations if len(l.access_rules) > 1]
len(multi_rule_locations) # 62 for me
get_location = iter(multi_rule_locations).__next__
```

This is one that expands from a constant rule into needing multiple checks:

```
Market Bombchu Shop -> Market Bombchu Shop Item 5: True

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has)
              4 LOAD_GLOBAL              1 (Progressive_Wallet)
              6 CALL_METHOD              1
              8 RETURN_VALUE

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has_any_of)
              4 LOAD_GLOBAL              1 (Bombchus_10)
              6 LOAD_GLOBAL              2 (Bombchus_5)
              8 LOAD_GLOBAL              3 (Bombchus)
             10 LOAD_GLOBAL              4 (Bombchus_20)
             12 BUILD_TUPLE              4
             14 CALL_METHOD              1
             16 RETURN_VALUE
```

```
GC Grotto -> GC Deku Scrub Grotto Left: can_stun_deku

  1           0 LOAD_FAST                1 (age)
              2 LOAD_CONST               1 ('adult')
              4 COMPARE_OP               2 (==)
              6 JUMP_IF_TRUE_OR_POP     30 (to 60)
              8 LOAD_FAST                0 (state)
             10 LOAD_METHOD              0 (has_any_of)
             12 LOAD_GLOBAL              1 (Bombchus_10)
             14 LOAD_GLOBAL              2 (Boomerang)
             16 LOAD_GLOBAL              3 (Slingshot)
             18 LOAD_GLOBAL              4 (Deku_Stick_Drop)
             20 LOAD_GLOBAL              5 (Bombchus)
             22 LOAD_GLOBAL              6 (Buy_Deku_Shield)
             24 LOAD_GLOBAL              7 (Deku_Shield_Drop)
             26 LOAD_GLOBAL              8 (Buy_Deku_Stick_1)
             28 LOAD_GLOBAL              9 (Deku_Nut_Drop)
             30 LOAD_GLOBAL             10 (Buy_Deku_Nut_5)
             32 LOAD_GLOBAL             11 (Bomb_Bag)
             34 LOAD_GLOBAL             12 (Buy_Deku_Nut_10)
             36 LOAD_GLOBAL             13 (Kokiri_Sword)
             38 LOAD_GLOBAL             14 (Bombchus_5)
             40 LOAD_GLOBAL             15 (Bombchus_20)
             42 BUILD_TUPLE             15
             44 CALL_METHOD              1
             46 JUMP_IF_TRUE_OR_POP     30 (to 60)
             48 LOAD_FAST                0 (state)
             50 LOAD_METHOD             16 (has_all_of)
             52 LOAD_GLOBAL             17 (Magic_Meter)
             54 LOAD_GLOBAL             18 (Dins_Fire)
             56 BUILD_TUPLE              2
             58 CALL_METHOD              1
        >>   60 RETURN_VALUE

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has)
              4 LOAD_GLOBAL              1 (Progressive_Wallet)
              6 LOAD_CONST               1 (0)
              8 CALL_METHOD              2
             10 RETURN_VALUE

Dodongos Cavern Torch Room -> Dodongos Cavern Deku Scrub Side Room Near Dodongos:  can_blast_or_smash or Progressive_Strength_Upgrade

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has_any_of)
              4 LOAD_GLOBAL              1 (Progressive_Strength_Upgrade)
              6 BUILD_TUPLE              1
              8 CALL_METHOD              1
             10 JUMP_IF_TRUE_OR_POP     25 (to 50)
             12 LOAD_FAST                0 (state)
             14 LOAD_METHOD              0 (has_any_of)
             16 LOAD_GLOBAL              2 (Bombchus_10)
             18 LOAD_GLOBAL              3 (Bombchus)
             20 LOAD_GLOBAL              4 (Bomb_Bag)
             22 LOAD_GLOBAL              5 (Bombchus_5)
             24 LOAD_GLOBAL              6 (Bombchus_20)
             26 BUILD_TUPLE              5
             28 CALL_METHOD              1
             30 JUMP_IF_TRUE_OR_POP     25 (to 50)
             32 LOAD_FAST                0 (state)
             34 LOAD_METHOD              7 (has_all_of)
             36 LOAD_GLOBAL              8 (Megaton_Hammer)
             38 BUILD_TUPLE              1
             40 CALL_METHOD              1
             42 JUMP_IF_FALSE_OR_POP    25 (to 50)
             44 LOAD_FAST                1 (age)
             46 LOAD_CONST               1 ('adult')
             48 COMPARE_OP               2 (==)
        >>   50 RETURN_VALUE

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has)
              4 LOAD_GLOBAL              1 (Progressive_Wallet)
              6 LOAD_CONST               1 (0)
              8 CALL_METHOD              2
             10 RETURN_VALUE
```

Bottle logic:
```
Market Potion Shop -> Market Potion Shop Item 2: True

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has)
              4 LOAD_GLOBAL              1 (Progressive_Wallet)
              6 LOAD_CONST               1 (2)
              8 CALL_METHOD              2
             10 RETURN_VALUE

 82           0 LOAD_FAST                0 (self)
              2 LOAD_METHOD              0 (has_any_of)
              4 LOAD_GLOBAL              1 (ItemInfo)
              6 LOAD_ATTR                2 (bottle_ids)
              8 CALL_METHOD              1
             10 JUMP_IF_TRUE_OR_POP     11 (to 22)
             12 LOAD_FAST                0 (self)
             14 LOAD_METHOD              3 (has)
             16 LOAD_GLOBAL              4 (Rutos_Letter)
             18 LOAD_CONST               1 (2)
             20 CALL_METHOD              2
        >>   22 RETURN_VALUE
```

This is a deku scrub that sells a potion for more than 99 rupees:
```
Ganons Castle Deku Scrubs -> Ganons Castle Deku Scrub Right: can_stun_deku

  1           0 LOAD_FAST                1 (age)
              2 LOAD_CONST               1 ('adult')
              4 COMPARE_OP               2 (==)
              6 JUMP_IF_TRUE_OR_POP     30 (to 60)
              8 LOAD_FAST                0 (state)
             10 LOAD_METHOD              0 (has_any_of)
             12 LOAD_GLOBAL              1 (Bombchus_10)
             14 LOAD_GLOBAL              2 (Boomerang)
             16 LOAD_GLOBAL              3 (Slingshot)
             18 LOAD_GLOBAL              4 (Deku_Stick_Drop)
             20 LOAD_GLOBAL              5 (Bombchus)
             22 LOAD_GLOBAL              6 (Buy_Deku_Shield)
             24 LOAD_GLOBAL              7 (Deku_Shield_Drop)
             26 LOAD_GLOBAL              8 (Buy_Deku_Stick_1)
             28 LOAD_GLOBAL              9 (Deku_Nut_Drop)
             30 LOAD_GLOBAL             10 (Buy_Deku_Nut_5)
             32 LOAD_GLOBAL             11 (Bomb_Bag)
             34 LOAD_GLOBAL             12 (Buy_Deku_Nut_10)
             36 LOAD_GLOBAL             13 (Kokiri_Sword)
             38 LOAD_GLOBAL             14 (Bombchus_5)
             40 LOAD_GLOBAL             15 (Bombchus_20)
             42 BUILD_TUPLE             15
             44 CALL_METHOD              1
             46 JUMP_IF_TRUE_OR_POP     30 (to 60)
             48 LOAD_FAST                0 (state)
             50 LOAD_METHOD             16 (has_all_of)
             52 LOAD_GLOBAL             17 (Magic_Meter)
             54 LOAD_GLOBAL             18 (Dins_Fire)
             56 BUILD_TUPLE              2
             58 CALL_METHOD              1
        >>   60 RETURN_VALUE

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has)
              4 LOAD_GLOBAL              1 (Progressive_Wallet)
              6 LOAD_CONST               1 (0)
              8 CALL_METHOD              2
             10 RETURN_VALUE

 82           0 LOAD_FAST                0 (self)
              2 LOAD_METHOD              0 (has_any_of)
              4 LOAD_GLOBAL              1 (ItemInfo)
              6 LOAD_ATTR                2 (bottle_ids)
              8 CALL_METHOD              1
             10 JUMP_IF_TRUE_OR_POP     11 (to 22)
             12 LOAD_FAST                0 (self)
             14 LOAD_METHOD              3 (has)
             16 LOAD_GLOBAL              4 (Rutos_Letter)
             18 LOAD_CONST               1 (2)
             20 CALL_METHOD              2
        >>   22 RETURN_VALUE
```

Tuple expansion:
```
Kak House of Skulltula -> Kak 50 Gold Skulltula Reward: (Gold_Skulltula_Token, 50)

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has)
              4 LOAD_GLOBAL              1 (Gold_Skulltula_Token)
              6 LOAD_CONST               1 (50)
              8 CALL_METHOD              2
             10 RETURN_VALUE

  1           0 LOAD_CONST               1 (True)
              2 RETURN_VALUE
```

Zora Tunic Logic:
```
ZD Shop -> ZD Shop Item 1: True

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has)
              4 LOAD_GLOBAL              1 (Progressive_Wallet)
              6 LOAD_CONST               1 (2)
              8 CALL_METHOD              2
             10 RETURN_VALUE

  1           0 LOAD_FAST                1 (age)
              2 LOAD_CONST               1 ('adult')
              4 COMPARE_OP               2 (==)
              6 RETURN_VALUE
```

### Spirit Temple

Let's turn up some of the infamous spirit key logic.

```python
spirit = [l for l in locations if "spirit temple" in l.name]
get_location = iter(spirit).__next__
```

For the first one, try to figure out what `34 LOAD_CONST               3 (())`
refers to

```
Spirit Temple Central Chamber -> Spirit Temple Sun Block Room Silver Rupee Center Front:  (Small_Key_Spirit_Temple, 3) or spirit_temple_shortcuts or has_explosives or ((Small_Key_Spirit_Temple, 2) and free_bombchu_drops)

  1           0 LOAD_FAST                0 (state)
              2 LOAD_METHOD              0 (has_any_of)
              4 LOAD_GLOBAL              1 (Bombchus_10)
              6 LOAD_GLOBAL              2 (Bombchus)
              8 LOAD_GLOBAL              3 (Bomb_Bag)
             10 LOAD_GLOBAL              4 (Bombchus_5)
             12 LOAD_GLOBAL              5 (Bombchus_20)
             14 BUILD_TUPLE              5
             16 CALL_METHOD              1
             18 JUMP_IF_TRUE_OR_POP     25 (to 50)
             20 LOAD_FAST                0 (state)
             22 LOAD_METHOD              6 (has)
             24 LOAD_GLOBAL              7 (Small_Key_Spirit_Temple)
             26 LOAD_CONST               1 (3)
             28 CALL_METHOD              2
             30 JUMP_IF_TRUE_OR_POP     25 (to 50)
             32 LOAD_CONST               2 ('Spirit Temple')
             34 LOAD_CONST               3 (())
             36 CONTAINS_OP              0
             38 JUMP_IF_TRUE_OR_POP     25 (to 50)
             40 LOAD_FAST                0 (state)
             42 LOAD_METHOD              6 (has)
             44 LOAD_GLOBAL              7 (Small_Key_Spirit_Temple)
             46 LOAD_CONST               4 (2)
             48 CALL_METHOD              2
        >>   50 RETURN_VALUE
```

This one shrinks considerably:
```
Spirit Temple Central Chamber -> Spirit Temple GS Lobby:  (is_child and logic_spirit_lobby_gs and Boomerang and (Small_Key_Spirit_Temple, 5)) or (is_adult and (Hookshot or Hover_Boots or logic_spirit_lobby_jump) and ((Small_Key_Spirit_Temple, 3) or spirit_temple_shortcuts)) or (logic_spirit_lobby_gs and Boomerang and (Hookshot or Hover_Boots or logic_spirit_lobby_jump) and (has_explosives or ((Small_Key_Spirit_Temple, 2) and free_bombchu_drops)))

  1           0 LOAD_FAST                1 (age)
              2 LOAD_CONST               1 ('adult')
              4 COMPARE_OP               2 (==)
              6 JUMP_IF_FALSE_OR_POP    20 (to 40)
              8 LOAD_FAST                0 (state)
             10 LOAD_METHOD              0 (has_any_of)
             12 LOAD_GLOBAL              1 (Progressive_Hookshot)
             14 LOAD_GLOBAL              2 (Hover_Boots)
             16 BUILD_TUPLE              2
             18 CALL_METHOD              1
             20 JUMP_IF_FALSE_OR_POP    20 (to 40)
             22 LOAD_FAST                0 (state)
             24 LOAD_METHOD              3 (has)
             26 LOAD_GLOBAL              4 (Small_Key_Spirit_Temple)
             28 LOAD_CONST               2 (3)
             30 CALL_METHOD              2
             32 JUMP_IF_TRUE_OR_POP     20 (to 40)
             34 LOAD_CONST               3 ('Spirit Temple')
             36 LOAD_CONST               4 (())
             38 CONTAINS_OP              0
        >>   40 RETURN_VALUE
```

Things quickly go down hill for us though:
```
Spirit Temple Central Chamber -> Spirit Temple Sun Block Room Chest:  (is_child and Sticks and (logic_spirit_sun_chest_no_rupees or (Silver_Rupee_Spirit_Temple_Sun_Block, 5)) and (Small_Key_Spirit_Temple, 5)) or (is_adult and (has_fire_source or (logic_spirit_sun_chest_bow and (Silver_Rupee_Spirit_Temple_Sun_Block, 5) and Bow)) and ((Small_Key_Spirit_Temple, 3) or spirit_temple_shortcuts)) or ((can_use(Dins_Fire) or (((Magic_Meter and Fire_Arrows) or (logic_spirit_sun_chest_bow and (Silver_Rupee_Spirit_Temple_Sun_Block, 5))) and Bow and Sticks and (logic_spirit_sun_chest_no_rupees or (Silver_Rupee_Spirit_Temple_Sun_Block, 5)))) and (has_explosives or ((Small_Key_Spirit_Temple, 2) and free_bombchu_drops)))

  1           0 LOAD_FAST                1 (age)
              2 LOAD_CONST               1 ('child')
              4 COMPARE_OP               2 (==)
              6 POP_JUMP_IF_FALSE       23 (to 46)
              8 LOAD_FAST                0 (state)
             10 LOAD_METHOD              0 (has_any_of)
             12 LOAD_GLOBAL              1 (Buy_Deku_Stick_1)
             14 LOAD_GLOBAL              2 (Deku_Stick_Drop)
             16 BUILD_TUPLE              2
             18 CALL_METHOD              1
             20 POP_JUMP_IF_FALSE       23 (to 46)
             22 LOAD_FAST                0 (state)
             24 LOAD_METHOD              3 (has)
             26 LOAD_GLOBAL              4 (Silver_Rupee_Spirit_Temple_Sun_Block)
             28 LOAD_CONST               2 (5)
             30 CALL_METHOD              2
             32 POP_JUMP_IF_FALSE       23 (to 46)
             34 LOAD_FAST                0 (state)
             36 LOAD_METHOD              3 (has)
             38 LOAD_GLOBAL              5 (Small_Key_Spirit_Temple)
             40 LOAD_CONST               2 (5)
             42 CALL_METHOD              2
             44 JUMP_IF_TRUE_OR_POP     99 (to 198)
        >>   46 LOAD_FAST                1 (age)
             48 LOAD_CONST               3 ('adult')
             50 COMPARE_OP               2 (==)
             52 POP_JUMP_IF_FALSE       56 (to 112)
             54 LOAD_FAST                0 (state)
             56 LOAD_METHOD              6 (has_all_of)
             58 LOAD_GLOBAL              7 (Magic_Meter)
             60 LOAD_GLOBAL              8 (Dins_Fire)
             62 BUILD_TUPLE              2
             64 CALL_METHOD              1
             66 POP_JUMP_IF_TRUE        46 (to 92)
             68 LOAD_FAST                0 (state)
             70 LOAD_METHOD              6 (has_all_of)
             72 LOAD_GLOBAL              7 (Magic_Meter)
             74 LOAD_GLOBAL              9 (Bow)
             76 LOAD_GLOBAL             10 (Fire_Arrows)
             78 BUILD_TUPLE              3
             80 CALL_METHOD              1
             82 POP_JUMP_IF_FALSE       56 (to 112)
             84 LOAD_FAST                1 (age)
             86 LOAD_CONST               3 ('adult')
             88 COMPARE_OP               2 (==)
             90 POP_JUMP_IF_FALSE       56 (to 112)
        >>   92 LOAD_FAST                0 (state)
             94 LOAD_METHOD              3 (has)
             96 LOAD_GLOBAL              5 (Small_Key_Spirit_Temple)
             98 LOAD_CONST               4 (3)
            100 CALL_METHOD              2
            102 JUMP_IF_TRUE_OR_POP     99 (to 198)
            104 LOAD_CONST               5 ('Spirit Temple')
            106 LOAD_CONST               6 (())
            108 CONTAINS_OP              0
            110 JUMP_IF_TRUE_OR_POP     99 (to 198)
        >>  112 LOAD_FAST                0 (state)
            114 LOAD_METHOD              6 (has_all_of)
            116 LOAD_GLOBAL              7 (Magic_Meter)
            118 LOAD_GLOBAL              8 (Dins_Fire)
            120 BUILD_TUPLE              2
            122 CALL_METHOD              1
            124 POP_JUMP_IF_TRUE        84 (to 168)
            126 LOAD_FAST                0 (state)
            128 LOAD_METHOD              6 (has_all_of)
            130 LOAD_GLOBAL              7 (Magic_Meter)
            132 LOAD_GLOBAL              9 (Bow)
            134 LOAD_GLOBAL             10 (Fire_Arrows)
            136 BUILD_TUPLE              3
            138 CALL_METHOD              1
            140 JUMP_IF_FALSE_OR_POP    99 (to 198)
            142 LOAD_FAST                0 (state)
            144 LOAD_METHOD              0 (has_any_of)
            146 LOAD_GLOBAL              1 (Buy_Deku_Stick_1)
            148 LOAD_GLOBAL              2 (Deku_Stick_Drop)
            150 BUILD_TUPLE              2
            152 CALL_METHOD              1
            154 JUMP_IF_FALSE_OR_POP    99 (to 198)
            156 LOAD_FAST                0 (state)
            158 LOAD_METHOD              3 (has)
            160 LOAD_GLOBAL              4 (Silver_Rupee_Spirit_Temple_Sun_Block)
            162 LOAD_CONST               2 (5)
            164 CALL_METHOD              2
            166 JUMP_IF_FALSE_OR_POP    99 (to 198)
        >>  168 LOAD_FAST                0 (state)
            170 LOAD_METHOD              0 (has_any_of)
            172 LOAD_GLOBAL             11 (Bombchus_10)
            174 LOAD_GLOBAL             12 (Bombchus)
            176 LOAD_GLOBAL             13 (Bomb_Bag)
            178 LOAD_GLOBAL             14 (Bombchus_5)
            180 LOAD_GLOBAL             15 (Bombchus_20)
            182 BUILD_TUPLE              5
            184 CALL_METHOD              1
            186 JUMP_IF_TRUE_OR_POP     99 (to 198)
            188 LOAD_FAST                0 (state)
            190 LOAD_METHOD              3 (has)
            192 LOAD_GLOBAL              5 (Small_Key_Spirit_Temple)
            194 LOAD_CONST               7 (2)
            196 CALL_METHOD              2
        >>  198 RETURN_VALUE


Child Spirit Temple Climb -> Spirit Temple Child Climb North Chest:  (is_child and has_projectile(child) and (Small_Key_Spirit_Temple, 5)) or (is_adult and has_projectile(adult) and ((Small_Key_Spirit_Temple, 3) or spirit_temple_shortcuts or ((Small_Key_Spirit_Temple, 2) and free_bombchu_drops))) or has_projectile(both)

  1           0 LOAD_FAST                1 (age)
              2 LOAD_CONST               1 ('child')
              4 COMPARE_OP               2 (==)
              6 POP_JUMP_IF_FALSE       22 (to 44)
              8 LOAD_FAST                0 (state)
             10 LOAD_METHOD              0 (has_any_of)
             12 LOAD_GLOBAL              1 (Bombchus_10)
             14 LOAD_GLOBAL              2 (Boomerang)
             16 LOAD_GLOBAL              3 (Slingshot)
             18 LOAD_GLOBAL              4 (Bombchus)
             20 LOAD_GLOBAL              5 (Bomb_Bag)
             22 LOAD_GLOBAL              6 (Bombchus_5)
             24 LOAD_GLOBAL              7 (Bombchus_20)
             26 BUILD_TUPLE              7
             28 CALL_METHOD              1
             30 POP_JUMP_IF_FALSE       22 (to 44)
             32 LOAD_FAST                0 (state)
             34 LOAD_METHOD              8 (has)
             36 LOAD_GLOBAL              9 (Small_Key_Spirit_Temple)
             38 LOAD_CONST               2 (5)
             40 CALL_METHOD              2
             42 JUMP_IF_TRUE_OR_POP     77 (to 154)
        >>   44 LOAD_FAST                1 (age)
             46 LOAD_CONST               3 ('adult')
             48 COMPARE_OP               2 (==)
             50 POP_JUMP_IF_FALSE       54 (to 108)
             52 LOAD_FAST                0 (state)
             54 LOAD_METHOD              0 (has_any_of)
             56 LOAD_GLOBAL              1 (Bombchus_10)
             58 LOAD_GLOBAL              4 (Bombchus)
             60 LOAD_GLOBAL             10 (Bow)
             62 LOAD_GLOBAL              5 (Bomb_Bag)
             64 LOAD_GLOBAL             11 (Progressive_Hookshot)
             66 LOAD_GLOBAL              6 (Bombchus_5)
             68 LOAD_GLOBAL              7 (Bombchus_20)
             70 BUILD_TUPLE              7
             72 CALL_METHOD              1
             74 POP_JUMP_IF_FALSE       54 (to 108)
             76 LOAD_FAST                0 (state)
             78 LOAD_METHOD              8 (has)
             80 LOAD_GLOBAL              9 (Small_Key_Spirit_Temple)
             82 LOAD_CONST               4 (3)
             84 CALL_METHOD              2
             86 JUMP_IF_TRUE_OR_POP     77 (to 154)
             88 LOAD_CONST               5 ('Spirit Temple')
             90 LOAD_CONST               6 (())
             92 CONTAINS_OP              0
             94 JUMP_IF_TRUE_OR_POP     77 (to 154)
             96 LOAD_FAST                0 (state)
             98 LOAD_METHOD              8 (has)
            100 LOAD_GLOBAL              9 (Small_Key_Spirit_Temple)
            102 LOAD_CONST               7 (2)
            104 CALL_METHOD              2
            106 JUMP_IF_TRUE_OR_POP     77 (to 154)
        >>  108 LOAD_FAST                0 (state)
            110 LOAD_METHOD              0 (has_any_of)
            112 LOAD_GLOBAL              1 (Bombchus_10)
            114 LOAD_GLOBAL              4 (Bombchus)
            116 LOAD_GLOBAL              5 (Bomb_Bag)
            118 LOAD_GLOBAL              6 (Bombchus_5)
            120 LOAD_GLOBAL              7 (Bombchus_20)
            122 BUILD_TUPLE              5
            124 CALL_METHOD              1
            126 JUMP_IF_TRUE_OR_POP     77 (to 154)
            128 LOAD_FAST                0 (state)
            130 LOAD_METHOD              0 (has_any_of)
            132 LOAD_GLOBAL              2 (Boomerang)
            134 LOAD_GLOBAL              3 (Slingshot)
            136 BUILD_TUPLE              2
            138 CALL_METHOD              1
            140 JUMP_IF_FALSE_OR_POP    77 (to 154)
            142 LOAD_FAST                0 (state)
            144 LOAD_METHOD              0 (has_any_of)
            146 LOAD_GLOBAL             10 (Bow)
            148 LOAD_GLOBAL             11 (Progressive_Hookshot)
            150 BUILD_TUPLE              2
            152 CALL_METHOD              1
        >>  154 RETURN_VALUE
```

---

[ln1]: https://github.com/ponyorm/pony/

[^fn1]: The [pony orm][ln1]  uses dis to take apart generator expressions and
    rewrite them into SQL queries.

[^fn2]: It's a real "the map is not the territory" situation.
