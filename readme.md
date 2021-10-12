
Minetest map parser library for go

![](https://github.com/minetest-go/mapparser/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/minetest-go/mapparser/badge.svg)](https://coveralls.io/github/minetest-go/mapparser)

# Features

* Extracts metadata/inventories
* NodeID/Param1/Param2 handling
* Supports zlib compressed mapblocks
* zstd compressed mapblocks (minetest 5.5+)

# Example

```go
// read mapblock data from db/file/somewhere else
data, err := ioutil.ReadFile("mapblock.bin")
if err != nil {
    panic(err)
}

// parse
mapblock, err := Parse(data)
if err != nil {
    panic(err)
}

// mapblock version
fmt.Printf("Version: %d", mapblock.Version)

// nodes
fmt.Printf("%s", mapblock.GetNodeName(10,0,2)) // node-name
fmt.Printf("%s", mapblock.GetNodeId(10,0,2)) //raw nodeid
fmt.Printf("%s", mapblock.GetParam2(10,0,2)) //param2

// node-id mapping
for id, name := range mapblock.BlockMapping {
    fmt.Printf("%d = %s", id, name)
}

// inventories
invMap := mapblock.Metadata.GetInventoryMapAtPos(10,0,0)
mainInv := invMap["main"]
for _, item := range mainInv {
    fmt.Printf("%s %d", item.Name, item.Count)
}

// metadata
md := mapblock.Metadata.GetMetadata(10,0,1)
for key, value := range md {
    fmt.Printf("%s = %s", key, value)
}
```