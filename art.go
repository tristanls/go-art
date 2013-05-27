// Go language port of Actor Run-Time (an abstract actor machine created by Dale Schumacher (http://dalnefre.com))
package art

import "fmt"
import "strings"

import "github.com/anisus/queue"
import "github.com/noll/samling/stack"

// Art holds the state of the Actor Run-Time
type Art struct {
  MachineState *stack.Stack
  Data *stack.Stack
  Scope *stack.Stack
  Code *stack.Stack
  Event *queue.Queue
  CurrentEvent *stack.Stack
}
// Load loads a program into the runtime
func (runtime *Art) Load(program string) *Art {
  runtime.Scope.Push("()")
  runtime.Code.Push(")")
  runtime.Code.Push("()")
  tokanized := Tokanize(program)
  for i := len(tokanized) - 1; i >= 0; i-- {
    runtime.Code.Push(tokanized[i])
  }
  runtime.Code.Push("(")
  return runtime
}
// String pretty-print current state of actor runtime
func (runtime *Art) String() string {
  stringValue := "{\n"
  stringValue += "  Machine State Stack: "
  stringValue += StackToString(runtime.MachineState) + "\n"
  stringValue += "  Data Stack:          "
  stringValue += StackToString(runtime.Data) + "\n"
  stringValue += "  Scope Stack:         "
  stringValue += StackToString(runtime.Scope) + "\n"
  stringValue += "  Code Stack:          "
  stringValue += StackToString(runtime.Code) + "\n"
  stringValue += "  Event Queue:         "
  stringValue += QueueToString(runtime.Event) + "\n"
  stringValue += "  Current Event Stack: "
  stringValue += StackToString(runtime.CurrentEvent) + "\n"
  stringValue += "\n}"
  return stringValue
}

// CreateArt creates a new Actor Run-Time in an initial configuration
func CreateArt() *Art {
  machineState := stack.New()
  machineState.Push("?")
  machineState.Push("0")

  data := stack.New()
  data.Push("?")

  scope := stack.New()
  scope.Push("?")

  code := stack.New()
  code.Push("?")

  event := queue.New()
  event.Enqueue("()")

  currentEvent := stack.New()
  currentEvent.Push("?")
  return &Art{
    MachineState: machineState,
    Data: data,
    Scope: scope,
    Code: code,
    Event: event,
    CurrentEvent: currentEvent,
  }
}

// QueueToString pretty-prints queue contents
func QueueToString(_queue *queue.Queue) string {
  stringValue := ""
  replacement := queue.New()

  for value, ok := _queue.Dequeue(); ok; value, ok = _queue.Dequeue() {
    stringValue += fmt.Sprintf("%s", value) + ","
    replacement.Enqueue(value)
  }
  stringValue = strings.TrimRight(stringValue, ",")

  for value, ok := replacement.Dequeue(); ok; value, ok = replacement.Dequeue(){
    _queue.Enqueue(value)
  }
  return stringValue
}

// StackToString pretty-prints stack contents
func StackToString(_stack *stack.Stack) string {
  stringValue := ""
  temporary := stack.New()

  for _stack.Len() > 0 {
    value := _stack.Pop()
    stringValue += fmt.Sprintf("%s", value)
    if (_stack.Len() > 0) {
      stringValue += ","
    }
    temporary.Push(value)
  }

  for temporary.Len() > 0 {
    _stack.Push(temporary.Pop())
  }
  return stringValue
}

// Tokanize tokanizes a program by appending '#' in front of each literal
func Tokanize(program string) []string {
  runes := strings.Split(program, "")
  tokenized := make([]string, len(runes))
  for i := 0; i < len(runes); i++ {
    tokenized[i] = "#" + runes[i]
  }
  return tokenized
}