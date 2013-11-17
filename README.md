goprocinfo
===================

/proc information parser for Go.

Usage
---------------

```go
import "github.com/c9s/goprocinfo"

stat, err := goprocinfo.ReadStat("proc/stat")
if err != nil {
    t.Fatal("stat read fail")
}

// stat.CPUStatAll
// stat.CPUStats
// stat.Processes
// ... etc
```


Reference
------------

* http://www.mjmwired.net/kernel/Documentation/filesystems/proc.txt
