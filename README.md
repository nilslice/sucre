# sucre
An easy to use library made for 2D drawing using OpenGL 3.2

## Usage with glfw

```go
package main

import "github.com/atrika/sucre"
import "github.com/go-gl/glfw/v3.1/glfw"
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

   var context sucre.Context   
   context.Initialize("/home/username/Images/textures")
   context.SetCameraPosition(0.0, 0.0)
   context.SetCameraAngle(0.0)
   context.SetCameraSize(20.0, 20.0)
   context.SetClearColor(sucre.Color{0.4, 0.1, 0.1})
   
   tex1, _ := context.GetTextureId("3.png")
   
   for !window.ShouldClose() {
      context.ClearScene()
      
      data := sucre.SquareData{PosX: 8.0, 
                               PosY: 8.0,
                               Depth: 0.5,
                               Angle: 0.0,
                               Size: 3.0}
                               
      context.AddSquare(tex1, data)
      
      context.Draw()
      
      window.SwapBuffers()
      glfw.PollEvents()
   }
}
```
