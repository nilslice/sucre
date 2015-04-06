package sucre

import "github.com/go-gl/gl/v3.2-core/gl"
import "sort"

type deeperFirst []innerSquareData
func (a deeperFirst) Len() int {return len(a)}
func (a deeperFirst) Swap(i, j int) {a[i], a[j] = a[j], a[i]}
func (a deeperFirst) Less(i, j int) bool {return a[i].Depth > a[j].Depth}

func drawSquares(squares []innerSquareData, texId uint32, transparent bool) {
   count := int32(len(squares))
   if count == 0 {
      return
   }
   
   if transparent {
      sort.Sort(deeperFirst(squares))
   }
   
   gl.BindTexture(gl.TEXTURE_2D_ARRAY, texId)

   // Upload squares
   gl.BufferData(gl.ARRAY_BUFFER, int(count * 6 * 4), gl.Ptr(squares), gl.DYNAMIC_DRAW)
   
   // Draw
   enableGlCaps(transparent)
   gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, count)
   disableGlCaps()
}

func enableGlCaps(transparent bool) {    
   // Depth Test
   gl.Enable(gl.DEPTH_TEST)
   if transparent {
      gl.DepthFunc(gl.LEQUAL)
   } else {
      gl.DepthFunc(gl.LESS)
   }
   
   // Blending
   if transparent {
      gl.Enable(gl.BLEND)
      gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
   }
   
   // Backface Culling
   gl.Enable(gl.CULL_FACE)
   gl.FrontFace(gl.CW)
   gl.CullFace(gl.BACK)
}

func disableGlCaps() {
   gl.Disable(gl.DEPTH_TEST)
   gl.Disable(gl.BLEND)
   gl.Disable(gl.CULL_FACE)
}
