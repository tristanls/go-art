# art

_Stability: 1 - [Experimental](https://github.com/tristanls/stability-index#stability-1---experimental)_

Actor Run-Time is an abstract actor machine created by [Dale Schumacher](http://dalnefre.com). This is a Go language port of it. The JavaScript version is [available online](http://www.dalnefre.com/humus/sim/art.html)

The source for this description is a post from [Computational Actors Guild](https://groups.google.com/d/topic/computational-actors-guild/y-8DfmY4v0g/discussion).

## Description

### Machine Registers

#### Machine State Stack - `M`

#### Data Stack - `D`

#### Scope Stack - `S`

#### Code Stack - `C`

#### Event Queue - `Q`

#### Current Event - `E`

### Meta-Variables

Meta-variables are named with lower-case letters (`m`, `b`, `s`, `x`, `y`, `z`, ...).

### "undefined" -  `?`

### Empty list "NIL" - `()`

### Parenthesis

Parenthesis are used for grouping, usually of tuples.

### Tuples

Tuples are formed from pairs, grouped right-to-left, separated by commas `,`.

For example, if `((0, 1), 2, 3) = (x, y)` then `x = (0, 1)` and `y = (2, 3)`.

### Suffix operators

There are two suffix operators, `<` and `>`, used to access the head and tail of a pair.

For example, if `D = ((0, 1), 2, 3)` then `D< = (0, 1)` and `D> = (2, 3)` and `D<>` = 1.

### Literal Symbol - `#`

### Comments - `//`

## Initial configuration

    M=0,?   // Machine State Stack
    D=?     // Data Stack
    S=?,()  // Scope Stack
    C=?     // Code Stack
    Q=()    // Event Queue
    E=?     // Current Event
    E=(m,(s,b)) m=E< s=E>< b=E>>

In the initial machine state, symbols are extracted from the head of the code sequence, and actions are taken based on the meaning of the symbol. Essentially, symbols represent machine instructions and operate mostly on registers.

Registers are lists only if the operation changes them.

The changed contents of a register are represtend by a prime symbol `'` following the register name. For example, most operations include `C'=C>`, which means that the code sequence register become the tail of the current value in the code sequence register. Here are a few sample instructions:

    M< = 0: case C< of                              // top-level loop
      ?     M'=4,M> E'=Q< D'=E'>>,E'<,E'><,D Q'=Q>    // next event
      ()    C'=C> S'=S>                               // next code block
      ... case C<< of                                 // execute command
        ##    D'=C<><,D C'=C<>>,C>                      // literal
        #_    D'=#_,D C'=C<>,C>                         // wildcard
        #^    D'=E>,D C'=C<>,C>                         // self
        #@    D'=new(D<<,D<>),D> C'=C<>,C>              // (s,b) create
        #!    D'=D>> C'=C<>,C> Q`=Q+(D><,D<)            // m a send
        #%    D'=D> C'=C<>,C> E><'=D<< E>>'=D<>         // (s,b) become
        #$    D'=S<,D C'=C<>,C>                         // state/scope access
        #<    D'=D<<,D> C'=C<>,C>                       // first
        #>    D'=D<>,D> C'=C<>,C>                       // rest
        #(    M'=0,(M,D) C'=C<>,C>                      // begin tuple
        #,    C'=C<>,C>                                 // pair separator (ign)
        #)    M'=1,M> C'=C<>,C>                         // end tuple
        #{    *=S<,D M'=0,(M>,*) D=* C'=C<>,C>          // begin cases
        #;    D'=(D><,D<),D>> C'=C<>,C>                 // next case
        #}    M'=3,M> D'=(),D C'=C<>,C>                 // end cases
        #.    M'=4,M> D'=D><>,D<,D><<,D>> C'=C<>,C>     // apply match
    M< = 1: case D of                               // tuple first
      M>>   M'=M>< D'=(),M>>                          // empty
      ...   M'=2,M>                                   // non-empty
    M< = 2: case D> of                              // tuple rest
      M>>   M'=M>< D'=D<,M>>                          // found begin
      ...   D'=(D><,D<),D>>                           // push item
    M< = 3: case D> of                              // build cases
      M>>   M'=0,M>< D'=(M>><,D<),M>>>                // found begin
      ...   D'=(D><,D<),D>>                           // push item
    M< = 4: case D< of                              // match cases
      ()    M'=0,M> D'=D>>>                           // done
      ...   *=D<<>,D<>,D> M'=5,(M>,*) D'=D<<<,D><,*   // match case
    M< = 5: case D of                               // end cases
      M>>   M'=0,M>< D'=D>>>> S'=(D>><,D>>><),S C'=D<,C // execute match
      ... case D< of                                  // match case
        (_,_)   D'=D<<,D><<,D<>,D><>,D>>                // de-tuple
        #_      D'=D>>                                  // wildcard matched
        D><     D'=D>>                                  // constant matched
        ...     M'=4,M>< D'=M>>>                        // next case

Examples:

`{_(#$,#<,());}(#0,#1).` evaluates to `(#0,#1)`

`{(#2,_)(#$,#<,#>,());_(#$,#<,#<,());}(#0,#1).` evaluates to `#0`

The [Humus](http://www.dalnefre.com/humus/sim/humus.html) expression `SEND 2 TO NEW \x.[SEND x TO SELF]` translates to:

    #2{_(#$,#<,#^,#!,());}@!

which will keep sending the message `2` to itself forever.