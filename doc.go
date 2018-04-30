// Package queue contains many different concurrent queue algorithms and tests
//
// All names follow  a convention: "[SM]P[SM]C[sw]?i?p?[rnacq]<variant>""
//
//     [SM]P:
//        supports either single `S` or multiple `M` concurrent producers
//
//     [SM]C:
//        supports either single `S` or multiple `M` concurrent consumers
//
//     [rnacq]: buffer implementation
//        `r` dynamically sized ring buffer
//        `n` node based,
//        `a` fixed size array based,
//        `c` channel based,
//        `q` dynamically sized ring buffer with sequence number.
//
//     [sw]?: waiting behavior
//        `s` when it is a spinning implementation,
//        `w` when it is partially spinning and partially waiting
//
//     i?: intrusiveness
//        `i` when it is an intrusive implementation (caller must allocate a node)
//
//     p?: buffer value padding
//        `p` when values are padded in the internal queue
//
//     <variant>:
//        special variant identifier for a particular implementation,
//        which indicates either base implementation author / paper / code.
//
package queue
