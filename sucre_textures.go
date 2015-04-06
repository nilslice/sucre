// Author : Philippe Trottier (atrika@github)

package sucre

import "github.com/go-gl/gl/v3.2-core/gl"
import "os"
import "image"
import "image/draw"
import "path/filepath"
import "regexp"
import _ "image/png"
import "fmt"
import "math"

// Load all textures of a directory into the context's dictionnary
func (this *Context) loadTextures(textureLocation string) {

   isImgFile, _ := regexp.Compile("\\.png$")
   
   this.texturesByName = make(map[string]uint32, 10)
   
   rgbas := make([]*image.RGBA, 0, 32)
   
   visit := func(path string, meta os.FileInfo, err error) error {
      if err == nil && isImgFile.MatchString(path) {
      
         rgba, invalid := getRGBA(path)
         if invalid != nil {
            return nil
         }
         
         this.texturesByName[filepath.Base(path)] = uint32(len(rgbas))
         rgbas = append(rgbas, rgba)
      }
      return nil
   }

   filepath.Walk(textureLocation, visit)
   
   // We assume every texture is of the same size, and a square
   texCount := int32(len(rgbas))
   if texCount == 0 {
      fmt.Printf("No textures found at %s" + "\n", textureLocation)
      return
   }
   size := int32(rgbas[0].Rect.Size().X)
   
   // We build mipmaps up to 32x32
   var mmCount = int32(math.Log2(float64(size)) - 4);
   if mmCount <= 0 {
      mmCount = 1
   }
      
   var theTexture uint32
   gl.GenTextures(1, &theTexture)
   gl.BindTexture(gl.TEXTURE_2D_ARRAY, theTexture)
   gl.TexStorage3D(gl.TEXTURE_2D_ARRAY, mmCount, gl.RGBA8, size, size, texCount);
      
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

