// Author : Philippe Trottier (atrika@github)

package sucre

import "github.com/go-gl/gl/v3.2-core/gl"
import "os"
import "image"
import "image/draw"
import "path/filepath"
import "regexp"
import _ "image/png"
import "math"

// Load all textures of a directory into the context's dictionnary
func (this *Context) loadTextures(textureLocation string) {

   isImgFile, _ := regexp.Compile("\\.png$")
   
   this.opaqueTexs = make(map[string]uint32, 32)
   this.transTexs  = make(map[string]uint32, 32)
   
   opaques := make([]*image.RGBA, 0, 32)
   transp  := make([]*image.RGBA, 0, 32)
   
   visit := func(path string, meta os.FileInfo, err error) error {
      if err == nil && isImgFile.MatchString(path) {
      
         rgba, invalid := getRGBA(path)
         if invalid != nil {
            return nil
         }
         
         if rgba.Opaque() {
            this.opaqueTexs[filepath.Base(path)] = uint32(len(opaques))
            opaques = append(opaques, rgba)
         } else {            
            this.transTexs[filepath.Base(path)]  = uint32(len(transp))
            transp  = append(transp,  rgba)
         }
      }
      return nil
   }

   filepath.Walk(textureLocation, visit)
   
   this.theOpaqueTex = upload(opaques, gl.RGB8)
   this.theTransTex  = upload(transp,  gl.RGBA8)
}

func upload(rgbas []*image.RGBA, internalFormat uint32) uint32 {
   // We assume every texture is of the same size, and a square
   texCount := int32(len(rgbas))
   if texCount == 0 {
      return 0
   }
   size := int32(rgbas[0].Rect.Size().X)
   
   // We build mipmaps down to 1x1
   var mmCount = int32(math.Log2(float64(size)) + 1);
   
   var theTexture uint32
   gl.GenTextures(1, &theTexture)
   gl.BindTexture(gl.TEXTURE_2D_ARRAY, theTexture)
      
   gl.TexStorage3D(gl.TEXTURE_2D_ARRAY, mmCount, internalFormat, size, size, texCount);
      
   for i, rgba := range rgbas {
      gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, // target
                       0,                   // LoD level
                       0,                   // x offset
                       0,                   // y offset
                       int32(i),            // z offset (used as texture id)
                       size,                // width
                       size,                // height
                       1,                   // depth (number of layers)
                       gl.RGBA,             // pixel layout
                       gl.UNSIGNED_BYTE,    // pixel data type
                       gl.Ptr(rgba.Pix))    // data
   }
   
   gl.GenerateMipmap(gl.TEXTURE_2D_ARRAY);
   gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR);
   gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR);
   gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
   gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
   
   return theTexture
}

// Get an RGBA image from disk
func getRGBA(imgPath string) (*image.RGBA, error) {
   imgFile, err := os.Open(imgPath)
   if err != nil {
      return nil, err
   }
   
   img, _, err := image.Decode(imgFile)
   if err != nil {
      return nil, err
   }
   
   rgba := image.NewRGBA(img.Bounds())
   
   draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
   
   return rgba, nil
}

