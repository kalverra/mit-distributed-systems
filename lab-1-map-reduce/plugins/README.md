# Plugins

Plugins for the Map Reduce function implementations.

Must follow this interface:

```go
func Map(string, string) []KeyValue {

}

func Reduce(string, []string) string {

}

```
