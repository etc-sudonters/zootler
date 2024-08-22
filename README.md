For funsies attempt at an alt-frontend for ootr, mostly I work on whatever I
find interesting at the moment

Minimal to no dependencies is a goal. Substrate is an exception because that
exists primarily as a place for me to put pieces I'm a little _too_ focused on
AND makes sense as a reusable module. 
 
## Storage

All of this used to be a big ol `map[reflect.Type]map[int][]interface{}` but
now `internal/table` and `internal/query` form the heart of this system.

`internal/table` is columnar-esque storage system. The table is divided into
independent columns that store components. A component is essentially any type
that is used to describe a row. A row is a collection of column entries joined
by a common rowid. There is no fixed schema for a row, rather projections are
assembled as needed. However, this means a row can only possess a component
either zero or one times.

There are several options for storing columns:
- `columns.Bit`: produces a singleton value if the rowid is present in its bitset
- `columns.Hashmap`: Stores components in a `map[table.RowId]table.Value`
- `columns.Slice`: Stores components in a `[]table.Value` indexed by
  `table.RowId`[^fn1]

Other columns backed by a sparse sets or trees are possible but unimplemented
currently.

A rudimentary indexing system is present to assist finding components with
specific characteristics. This falls back to a column scan and typically
`reflect.DeepEqual` if an index isn't present on a scanned column. Every column
effectively has a bitset index tracking membership.

`internal/query` provides some abstraction over the table, primarily it
provides an interface for gathering columns from the table and iterating over
the matching rows. For a row to match it must:

- be present in every column named by `Query.Exists` or `Query.Load`
- not be present in every column named by `Query.NotExists`

Rows are not matched by any particular property of its column value -- that is
handled by the similar `Engine.Lookup` method. Rather all that matters is if a
column has a value for a row or not.

`query.Engine` also provides facilities for creating columns, inserting and
removing rows from columns, and most importantly provides a mapping between
types and column ids. `table.Table` doesn't make efforts to ensure its storing
an appropriate type in a column it just adds the value to specified
column[^fn2].

`internal/bundle` provides iterators over `query.Engine` queries and lookups.

### Pending

Investigate using "archetypes":

```go
type SongLocation struct {
    Name     components.Name
    Song     components.Song
    Location components.Location
}
var arch SongLocation

q := engine.CreateQuery()
q.LoadArchectype(arch) // not this but the idea
rows, _ := engine.Retrieve(q)

for rowid, values := range rows.All {
    (&arch).Init(rowid, values)
    fmt.Printl("%+v", arch)
    // SongLocation{Name: components.Name("first"), Song: components.Song{}, Location: components.Location{}}
}
```



## Logic

The necessary files can be dumped from a local copy of OOTR's source code via
the `dump-zootr.py` helper. This will copy over the logic files and dump item
and location representations to json files. Most of the loading of these files
is the responsibility of the calling application, but `internal/rules` handles
transforming the logic json files into bytecode. 

The logic files describe connections between game world locations -- edge rules
-- using a (subset) of Python. `rules/parser` produces an AST from the edge
rule or helper passed. `rules/runtime` accepts this AST and produces chunks of
bytecode, constants and names. `rules/runtime` also provides a VM to execute
this chunk. The VM/bytecode is currently a mostly 1:1 mapping from the AST
leading to awkwardness in decisions like not supporting `in` or subscripting
in the VM. The VM supports calling both the compiled OOTR helper functions
and Go.

The game world is modeled by a directed graph built from the locations, exit
edge rules, and check edge rules.[^fn3] I haven't quite decided on where/how to
stash these chunks but throwing them into the table is appealing. 

### Filling



---

[^fn1]: Note this means the length of the array is _at least_ the highest rowid
    _ever_ tracked by the column

[^fn2]: This is pretty intentional since there should never be an incorrect
    placement when operating the table via the query interface.

[^fn3]: OOTR distinguishes between checks and events however I haven't found
    any benefit for this yet.
