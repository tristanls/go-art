package main

import "fmt"
import "github.com/tristanls/go-art"

func main() {
  runtime := art.CreateArt()
  fmt.Println(runtime)
  runtime = runtime.Load("(#x,#y,())><")
  fmt.Println(runtime)
}