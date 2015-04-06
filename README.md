# sucre
An easy to use library made for 2D drawing using OpenGL 3.2

## Installation
go get github.com/atrika/sucre

## Documentation
https://godoc.org/github.com/atrika/sucre

## Requirements
The textures used must be squares, and they must all have the same sizes.
For example :
```
   /home/username/Images/textures
      texture1.png (512x512)
      texture2.png (512x512)
```
## Example using glfw

![result of example](https://i.imgur.com/PP6uuuj.png)

```go
package main

import "github.com/atrika/sucre"
import "github.com/go-gl/glfw/v3.1/glfw"
import "github.com/go-gl/gl/v3.2-core/gl"
import "runtime"
import "math"

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
   context.Initialize("/home/username/Images/textures/512")
   context.SetCameraPosition(0.0, 0.0)
   context.SetCameraAngle(0.0)
   context.SetCameraSize(10.0, 10.0)
   context.SetClearColor(sucre.Color{0.4, 0.1, 0.1})
      
   tex7, _  := context.GetTextureId("7.png")
   tex5, _  := context.GetTextureId("5.png")
   
   for !window.ShouldClose() {
      context.ClearScene()
      
      basic := sucre.BasicRectData{
         Width:  3.0,
         Height: 3.0,
         Angle:  0.0,
         PosX:   3.0,
         PosY:   3.0,
         Depth:  0.5,
      }      
      rect := sucre.RectData{basic, tex7}
      
      // Square of size 3 at (3,3) using 7.png
      context.AddRect(rect)
      
      rect.PosX    = -3.0
      rect.PosY    = -3.0
      rect.Height  =  2.0
      rect.Angle   = float32(math.Pi / 6.0)
      rect.Texture = tex5
      
      // Rotated 3x2 rectangle at (-3, -3) using 5.png
      context.AddRect(rect)
      
      rect.Width   = 1.0
      rect.Height  = 1.0
      rect.Depth   = 0.6 
      rect.Texture = tex7
      
      // Rotated square of size 1 at (-3, -3), behind the
      //  3x2 rectangle and using 7.png
      context.AddRect(rect)
      
      context.Draw()
      
      window.SwapBuffers()
      glfw.PollEvents()
   }
}
```
