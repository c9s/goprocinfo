goprocinfo
===================

/proc information parser for Go.

Usage
---------------

```go
import linuxproc "github.com/c9s/goprocinfo/linux"

stat, err := linuxproc.ReadStat("proc/stat")
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
