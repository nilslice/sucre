// Author : Philippe Trottier (atrika@github)

package sucre

import "github.com/go-gl/gl/v3.2-core/gl"
import "os"
import "image"
import "image/draw"
import "path/filepath"
import "regexp"
import _ "image/png"

// Load all textures of a directory into the context's dictionnary
func (this *Context) loadTextures(textureLocation string) {

   isImgFile, _ := regexp.Compile("\\.png$")
   
   this.texturesByName = make(map[string]uint32, 10)
   
   visit := func(path string, meta os.FileInfo, err error) error {
      if err == nil && isImgFile.MatchString(path) {
      
         rgba, invalid := getRGBA(path)
         if invalid != nil {
            return nil
         }
         
         tex := uploadRGBA(rgba)
         
         this.texturesByName[filepath.Base(path)] = tex
      }
      return nil
   }

   filepath.Walk(textureLocation, visit)
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

// Upload an RGBA image to the GPU as a texture
func uploadRGBA(rgba *image.RGBA) uint32 {
   var tex uint32
   gl.GenTextures(1, &tex)
   gl.BindTexture(gl.TEXTURE_2D, tex)
   gl.TexImage2D(gl.TEXTURE_2D,             // target
                 0,                         // mipmap level
                 gl.RGBA,                   // number of color components
                 int32(rgba.Rect.Size().X), // width
		           int32(rgba.Rect.Size().Y), // height
		           0,                         // border
		           gl.RGBA,                   // data format
		           gl.UNSIGNED_BYTE,          // data type of pixel
		           gl.Ptr(rgba.Pix))          // raw data
		           
   gl.GenerateMipmap(gl.TEXTURE_2D);
   gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
   gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
   gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR);
   gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR);
   
   gl.BindTexture(gl.TEXTURE_2D, 0)
   return tex
}


