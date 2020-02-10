# GoMDSlides

====

# Writing

- Write markdown into contents.md
- wasmserve
- open http://localhost:5000

====

# Deploy Example

- make deploy
- cd dist && python3 -m http.server
- open http://localhost:8000

====

# Usage

- `<-` preve page
- `->` next page
- `s` spotloght
- `f` fullscreen
- `r` reload contents.md

====

# Code

```go
package main
import "fmt"
func main() {
  fmt.Println("Hello!")
}
```

====

# fragment

```markdown
## Item1 {.fragment}

## Item2 {.fragment}
```

## Item1 {.fragment}

## Item2 {.fragment}

====

# End
