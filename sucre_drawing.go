package sucre

import "github.com/go-gl/gl/v3.2-core/gl"

// ----------------------------------------------------------------
// --------------------- Transparency Sorting ---------------------
// ----------------------------------------------------------------

type deeperFirst []SquareData

func (a deeperFirst) Len() int {
   return len(a)
}

func (a deeperFirst) Swap(i, j int) {
   a[i], a[j] = a[j], a[i]
}

func (a deeperFirst) Less(i, j int) bool {
   return a[i].Depth > a[j].Depth
}

// ----------------------------------------------------------------
// -------------------- OpenGL State Modifiers --------------------
// ----------------------------------------------------------------

func (this *Context) bindBuffersForDrawing() {
   gl.UseProgram(this.theProgram) 
   gl.BindVertexArray(this.theVAO) 
   gl.BindTexture(gl.TEXTURE_2D_ARRAY, this.theTexture)
}

func (this *Context) enableGlCaps() {    
   // Depth Test
   if !this.transparencyEnabled {
      gl.Enable(gl.DEPTH_TEST)
      gl.DepthFunc(gl.LESS)
      gl.ClearDepth(1.0)
   }
   
   // Blending
   if this.transparencyEnabled {
      gl.Enable(gl.BLEND)
      gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
   }
   
   // Backface Culling
   gl.Enable(gl.CULL_FACE)
   gl.FrontFace(gl.CW)
   gl.CullFace(gl.BACK)
}

func (this *Context) disableGlCaps() {
   gl.Disable(gl.DEPTH_TEST)
   gl.Disable(gl.BLEND)
   gl.Disable(gl.CULL_FACE)
}
