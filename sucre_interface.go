// Author : Philippe Trottier (atrika@github)

package sucre

import "math"
import "github.com/go-gl/gl/v3.2-core/gl"
import "sort"

// ----------------------------------------------------------------
// ------------------------ Data Structures -----------------------
// ----------------------------------------------------------------

type Color struct {
   R, G, B float32
}

type SquareData struct {
   PosX, PosY float32
   Size       float32
   Angle      float32       // in radians
   Depth      float32       // 0.0 <= Depth < 1.0
   TextureId  uint32
}

// ----------------------------------------------------------------
// ------------------------ General Stuff -------------------------
// ----------------------------------------------------------------

type Context struct {
   
   // Textures
   texturesByName map[string]uint32
   transparencyEnabled bool
   
   // Instance stuff
   squares []SquareData
   
   // Shader stuff
   instanceBuffer uint32
   zoomUni, rotUni, cameraUni int32
   theProgram, theVAO uint32
}

// Initializes OpenGL and loads the textures from disk
func (this *Context) Initialize(textureLocation string, transparencyEnabled bool) error {
   if err := gl.Init(); err != nil {
      return err
   }
   
   this.transparencyEnabled = transparencyEnabled
   
   // Depth Test
   gl.Enable(gl.DEPTH_TEST)
   if transparencyEnabled {   
      gl.DepthFunc(gl.LEQUAL)
   } else {       
      gl.DepthFunc(gl.LESS)
   }
   gl.ClearDepth(1.0)
   
   // Blending
   if transparencyEnabled {
      gl.Enable(gl.BLEND)
      gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
   }
   
   // Backface Culling
   gl.Enable(gl.CULL_FACE)
   gl.FrontFace(gl.CW)
   gl.CullFace(gl.BACK)
   
   // Create the shader program
   this.theProgram = createProgram()
   
   // Upload the square mesh and set up the program inputs
   this.initVAO()
   
   // Use the VAO and Program
   gl.BindBuffer(gl.ARRAY_BUFFER, this.instanceBuffer)
   gl.UseProgram(this.theProgram)
   gl.BindVertexArray(this.theVAO)   
   
   this.loadTextures(textureLocation)
   
   this.squares = make([]SquareData, 0, 32)
   
   this.SetClearColor(Color{0, 0, 0})
   this.SetCameraPosition(0.0, 0.0)
   this.SetCameraSize(1.0, 1.0)
   this.SetCameraAngle(0.0)
   
   return nil
}

// ----------------------------------------------------------------
// ------------------------- Camera Setup -------------------------
// ----------------------------------------------------------------

// Sets the position of the camera (default is [0;0])
func (this *Context) SetCameraPosition(posX, posY float32) {   
   gl.Uniform2f(this.cameraUni, posX, posY)
}

// Sets the size of the camera (default 1.0 for w and h)
func (this *Context) SetCameraSize(width, height float32) {   
   xZoom := 1.0 / width
   yZoom := 1.0 / height  
   
   zoom := [...]float32{
      xZoom,   0.0, 0.0, 0.0,
        0.0, yZoom, 0.0, 0.0,
        0.0,   0.0, 1.0, 0.0,
        0.0,   0.0, 0.0, 1.0 }
         
   gl.UniformMatrix4fv(this.zoomUni, 1, true, &zoom[0])  
}

// Sets the angle of the camera in radians (default is 0)
func (this *Context) SetCameraAngle(rad float64) {  
   rad = -rad;
 
   sinT := float32(math.Sin(rad))
   cosT := float32(math.Cos(rad))
   
   rot := [...]float32{
      cosT, -sinT, 0.0, 0.0,
      sinT,  cosT, 0.0, 0.0,
       0.0,   0.0, 1.0, 0.0,
       0.0,   0.0, 0.0, 1.0 }
         
   gl.UniformMatrix4fv(this.rotUni, 1, true, &rot[0])  
}

// ----------------------------------------------------------------
// ------------------------- Square Setup -------------------------
// ----------------------------------------------------------------

// Gets the ID of a texture
func (this *Context) GetTextureId(name string) (uint32, bool){
   textureId, ok := this.texturesByName[name]
   return textureId, ok
}

// Adds a square to be drawn in the next Draw call
func (this *Context) AddSquare(data SquareData) {   
   this.squares = append(this.squares, data)
}

// ----------------------------------------------------------------
// --------------------- Transparency Sorting ---------------------
// ----------------------------------------------------------------

type byDepth []SquareData

func (a byDepth) Len() int {
   return len(a)
}

func (a byDepth) Swap(i, j int) {
   a[i], a[j] = a[j], a[i]
}

func (a byDepth) Less(i, j int) bool {
   return a[i].Depth > a[j].Depth
}


// ----------------------------------------------------------------
// ------------------------- Scene Control ------------------------
// ----------------------------------------------------------------


// Sets the color used to clear the screen (default is black)
func (this *Context) SetClearColor(bg Color) {
   gl.ClearColor(bg.R, bg.G, bg.B, 1.0)
}

// Clears the scene of all squares
func (this *Context) ClearScene() {
   gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Draws the squares
func (this *Context) Draw() {
   count := int32(len(this.squares))
   if count == 0 {
      return
   }
   
   if this.transparencyEnabled {
      sort.Sort(byDepth(this.squares))
   }

   // Upload squares and draw call
   gl.BufferData(gl.ARRAY_BUFFER, int(count * 6 * 4), gl.Ptr(this.squares), gl.DYNAMIC_DRAW)
   gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, count)
      
   // Clear the square buffer
   this.squares = this.squares[:0]
}
