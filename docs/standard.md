# SeeGoals 𓅰 - Go Coding standard
Go is an opinionated language. This means there are substantially fewer coding standard decisions compared to other languages. Future holy wars will have to decide whether this is a good approach. As for this project, our aim is to follow any and all official guidelines for idiomatic Go as best as possible. In cases where there are no official guidelines, this document will be updated based on team decisions.

## Idiomatic Go
To follow Go's official guidelines for idiomatic Go, we have two main rules:

1. Always use `gofmt` when writing code. This ensures formatting style consistency across the project. In vscode this can be configured to run automatically (with an extension).

1. Always refer to Go's official [Effective Go](https://go.dev/doc/effective_go) as the baseline for idiomatic Go.

### Naming
Naming conventions are covered by [Effective Go](https://go.dev/doc/effective_go), and this project intends to follow that. However, it tends to be a common question in any project, so here is it explicitly.

Generally, the Go philosophy is to keep names short and concise. Longer names do **not** automatically make things clearer. Added context should instead be inferred (and provided) through things like types, comments, etc.

#### Mixed caps
The convention in Go is to use **mixed caps** for functions and variables. That means when there are multiple words in a function or variable, it should be either `camelCase` or `UpperCamelCase`. Since Go uses the first letter to decide visibility, whether or not capitalize is not a style choice, but a functional one.

#### Packages
Package names are typically single word, all lower case, and to the point. It's generally preferred to use directory context to provide additional sematic context. For example, a package `base64` could be located in `encoding/base64`. The directory structure makes it clear that `base64` specifically relates to encoding - no need to name the package `encoding_base64`, for example.

When creating packages we **don't** worry about naming collisions. If a collision does arise, it should be handled at the import by renaming the package identifier
```go
import customname "mypackage"
```

## SeeGoals standard 𓅰
Idiomatic Go does not cover everything needed for a standard. The following are specific decisions made by the SeeGoals team.

### Documentation & comments
Most developers agree that documentation is important. Despite this, documentation is often an afterthought or forgotten about entirely. To strike a balance between workload and usefulness, this project has two guidelines:

1. [Go Doc comments](https://tip.golang.org/doc/comment) MUST be used wherever possible. The only exception is packages, for which it should be considered on case by case basis.

1. Line comments are not a requirement, and should generally only be used when code is otherwise unclear or ambiguous.

#### Package definitions
Doc comments for packages are not a strict requirement. They should be decided on a case by case basis. Typically, we try to use it when there is a clear benefit to have a high level overview of package.

#### Type definitions
Type definitions should always have Doc comments.

For structs, this typically looks like this:
```go
// Buffer object used by client.
//
// A reference to this object can be obtained
// through the client, however care must be
// taken to ensure that access is synchronized.
type Buffer struct {
    // Holds underlying buffer data.
    //
    // This is not directly exported, so please
    // see the associated access methods.
    data []byte
}
```

Normal type definitions:
```go
// Type alias used for readability
//
// Use is mostly internal for this package,
// but it can be import and used elsewhere if needed.
type Index = int64
```

Lastly, interfaces are similar to structs:
```go
// Interface for buffer types
//
// Should be implemented by any buffer-like
// object to abstract the specific implementation.
type BufferObject interface {
    // Get a raw pointer to buffer data.
    //
    // This method does NOT guarantee a
    // non-nil pointer.
    GetDataPtr() *[]byte

    // Get a byte from at index.
    //
    // May return error when index
    // is out of range.
    GetByteAt(i Index) (byte, error)
}
```

#### Function definitions
A popular choice in many languages is to use Doxygen for function comments. Unfortunately, there is no Go support in Doxygen. Nor is there in-built support for `@param`-style comments. As there isn't good support for parameter level comments, this project will not be using it. Instead we rely on two things:

1. Good names for parameters and type definitions. Clear and informative naming should be considered as way to document code.

1. Informative Doc comments. If some detail of a function is hard to grasp, it should be explained in text by the function's Doc comment.

Example:
```go
// Try to connect client to address at some port.
//
// Address should NOT contain a port number - ex. "123.123.123.123"
// May return an error if the connection fails.
//
// Example:
//
//	c := NewClient()
//	if err := c.Connect("123.123.123.123", 8080); err != nil {
//		...
//	}
func (c *Client) Connect(address string, port int) error {
    ...
}
```
The above example highlights three important aspects of good function comments. First, it's not immediately clear what format `address` should have, hence it should be explained explicitly. Second, the function may return an error, so the comment should make clear why this may be returned. Third, a simple example of how to use this function. This not generally necessary, but often greatly appreciated.

### Misc
This section contains miscellaneous guidelines for this project. When this document is updated, and new information doesn't fit naturally anywhere else, it should be added here.

#### Renaming package identifiers
It's okay to rename package identifiers when imported as way to make provide added semantics. This should be used when the identifier is otherwise unclear in the import context:
```go
import goodname "dumbname"
```

#### File naming
WIP