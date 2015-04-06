# sucre
An easy to use library made for 2D drawing using OpenGL 3.2

## Installation
go get github.com/atrika/sucre

## Documentation
https://godoc.org/github.com/atrika/sucre

## Requirements
The textures used must be squarish, and they must all have the same sizes.
For example :
```
   /home/username/Images/textures
      texture1.png (512x512)
      texture2.png (512x512)
```

## Usage with glfw

```go
package main

import "github.com/atrika/sucre"
import "github.com/go-gl/glfw/v3.1/glfw"
import "github.com/go-gl/gl/v3.2-core/gl"
import "runtime"

func init() {
   runtime.LockOSThread()
}

func main() {
   err := glfw.Init()
   if err != nil {
      panic(err)
   }
   defer glfw.Terminate()

   glfw.WindowHint(glfw.ContextVersionMajor, 3)
   glfw.WindowHint(glfw.ContextVersionMinor, 2)
   glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

   window, err := glfw.CreateWindow(800, 800, "Testing", nil, nil)
   if err != nil {
      panic(err)
   }

   window.MakeContextCurrent()

   if err := gl.Init(); err != nil {
      panic(err)
   }

   var context sucre.Context
   context.Initialize("/home/username/Images/textures")
   context.SetCameraPosition(0.0, 0.0)
   context.SetCameraAngle(0.0)
   context.SetCameraSize(20.0, 20.0)
   context.SetClearColor(sucre.Color{0.4, 0.1, 0.1})

   tex1, _ := context.GetTextureId("texture1.png")

   for !window.ShouldClose() {
      context.ClearScene()

      basic := sucre.BasicSquareData{PosX:  8.0, 
                                     PosY:  8.0,
                                     Depth: 0.5,
                                     Angle: 0.0,
                                     Size:  3.0}
      data := sucre.SquareData{basic, tex1}

      context.AddSquare(data)

      context.Draw()

      window.SwapBuffers()
      glfw.PollEvents()
   }
}
```
