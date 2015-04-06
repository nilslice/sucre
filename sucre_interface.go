// Author : Philippe Trottier (atrika@github)

package sucre

import "math"
import "github.com/go-gl/gl/v3.2-core/gl"

// ----------------------------------------------------------------
// ------------------------ Data Structures -----------------------
// ----------------------------------------------------------------

type Color struct {
   R, G, B float32
}

type BasicRectData struct {   
   PosX, PosY    float32
   Width, Height float32
   Angle         float32       // in radians
   Depth         float32       // 0.0 <= Depth < 1.0
}

type innerRectData struct {
   BasicRectData
   TextureId  uint32
}

type RectData struct {
   BasicRectData
   Texture
}

type Texture struct {
   Id          uint32
   Transparent bool
}

// ----------------------------------------------------------------
// ------------------------ General Stuff -------------------------
// ----------------------------------------------------------------

type Context struct {
   
   // Textures
   opaqueTexs map[string]uint32
   transTexs  map[string]uint32
   
   // Instance stuff
   transRects  []innerRectData
   opaqueRects []innerRectData
   
   // Buffers
   theOpaqueTex   uint32
   theTransTex    uint32
   instanceBuffer uint32
   
   // Shader stuff
   zoomUni, rotUni, cameraUni int32
   theProgram, theVAO uint32
   
   bg Color   
}

// Initializes OpenGL and loads the textures from disk
func (this *Context) Initialize(textureLocation string) {
   
   // Create the shader program
   this.theProgram = createProgram()
   
   // Upload the Rect mesh and set up the program inputs
   this.initVAO()
   
   // Load the textures from disk
   this.loadTextures(textureLocation)
   
   this.transRects  = make([]innerRectData, 0, 32)
   this.opaqueRects = make([]innerRectData, 0, 32)
   this.SetClearColor(Color{0, 0, 0})
   
   // Initialize Camera
   this.SetCameraPosition(0.0, 0.0)
   this.SetCameraSize(1.0, 1.0)
   this.SetCameraAngle(0.0)
}

// ----------------------------------------------------------------
// ------------------------- Camera Setup -------------------------
// ----------------------------------------------------------------

// Sets the position of the camera (default is [0;0])
func (this *Context) SetCameraPosition(posX, posY float32) { 
   gl.UseProgram(this.theProgram)    
   gl.Uniform2f(this.cameraUni, posX, posY)
}

// Sets the size of the camera (default 1.0 for w and h)
func (this *Context) SetCameraSize(width, height float32) {   
   xZoom := 2.0 / width
   yZoom := 2.0 / height  
   
   zoom := [...]float32{
      xZoom,   0.0, 0.0, 0.0,
        0.0, yZoom, 0.0, 0.0,
        0.0,   0.0, 1.0, 0.0,
        0.0,   0.0, 0.0, 1.0 }
         
   gl.UseProgram(this.theProgram)  
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
   
   gl.UseProgram(this.theProgram)  
   gl.UniformMatrix4fv(this.rotUni, 1, true, &rot[0])
}

// ----------------------------------------------------------------
// ------------------------- Rect Setup -------------------------
// ----------------------------------------------------------------

// Gets the ID of a texture
func (this *Context) GetTextureId(name string) (Texture, bool){
   textureId, ok := this.opaqueTexs[name]   
   if ok {
      return Texture{textureId, false}, ok
   }
   textureId, ok = this.transTexs[name]
   return Texture{textureId, true}, ok
}

// Adds a Rect to be drawn in the next Draw call
func (this *Context) AddRect(data RectData) {
   inner := innerRectData{data.BasicRectData, data.Id}
   if data.Transparent {
      this.transRects  = append(this.transRects, inner)
   } else {
      this.opaqueRects = append(this.opaqueRects, inner)
   }  
}

// ----------------------------------------------------------------
// ------------------------- Scene Control ------------------------
// ----------------------------------------------------------------


// Sets the color used to clear the screen (default is black)
func (this *Context) SetClearColor(bg Color) {
   this.bg = bg
}

// Clears the scene of all Rects
func (this *Context) ClearScene() {
   gl.ClearColor(this.bg.R, this.bg.G, this.bg.B, 1.0)
   gl.ClearDepth(1.0)
   gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// Draws the Rects
func (this *Context) Draw() {

   // Bind what's needed to draw
   gl.UseProgram(this.theProgram) 
   gl.BindVertexArray(this.theVAO) 
   gl.BindBuffer(gl.ARRAY_BUFFER, this.instanceBuffer)
   
   // Start with the opaques  
   drawRects(this.opaqueRects, this.theOpaqueTex, false)
   // Then draw the transparents
   drawRects(this.transRects, this.theTransTex, true)
   
   // Clear the Rects
   this.opaqueRects = this.opaqueRects[:0]
   this.transRects = this.transRects[:0]
}
