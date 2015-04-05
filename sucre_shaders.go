// Author : Philippe Trottier (atrika@github)

package sucre

import "github.com/go-gl/gl/v3.2-core/gl"
import "fmt"

// Sets up the pointers to the params of the vertex shader
func (this *Context) initVAO() {
   // The VAO
   var theVAO uint32
   gl.GenVertexArrays(1, &theVAO)
   gl.BindVertexArray(theVAO)
   this.theVAO = theVAO
   
   // Instancing VBO
   var instanceVBO uint32
   gl.GenBuffers(1, &instanceVBO)
   gl.BindBuffer(gl.ARRAY_BUFFER, instanceVBO)
   this.instanceBuffer = instanceVBO
   
   worldPos := uint32(gl.GetAttribLocation(this.theProgram, gl.Str("world_pos" + "\x00")))
   gl.EnableVertexAttribArray(worldPos)
   gl.VertexAttribPointer(worldPos, 2, gl.FLOAT, false, 4 * 5, gl.PtrOffset(0))
   gl.VertexBindingDivisor(worldPos, 1)
   
   scalePos := uint32(gl.GetAttribLocation(this.theProgram, gl.Str("scale" + "\x00")))
   gl.EnableVertexAttribArray(scalePos)
   gl.VertexAttribPointer(scalePos, 1, gl.FLOAT, false, 4 * 5, gl.PtrOffset(2 * 4))
   gl.VertexBindingDivisor(scalePos, 1)
   
   anglePos := uint32(gl.GetAttribLocation(this.theProgram, gl.Str("angle" + "\x00")))
   gl.EnableVertexAttribArray(anglePos)
   gl.VertexAttribPointer(anglePos, 1, gl.FLOAT, false, 4 * 5, gl.PtrOffset(3 * 4))
   gl.VertexBindingDivisor(anglePos, 1)
   
   depthPos := uint32(gl.GetAttribLocation(this.theProgram, gl.Str("depth" + "\x00")))
   gl.EnableVertexAttribArray(depthPos)
   gl.VertexAttribPointer(depthPos, 1, gl.FLOAT, false, 4 * 5, gl.PtrOffset(4 * 4))
   gl.VertexBindingDivisor(depthPos, 1)   
   
   // Square Mesh
   squareVBO := uploadSquareMesh()
   gl.BindBuffer(gl.ARRAY_BUFFER, squareVBO)
   
   vertexPos := uint32(gl.GetAttribLocation(this.theProgram, gl.Str("mesh_pos" + "\x00")))
   gl.EnableVertexAttribArray(vertexPos)
   gl.VertexAttribPointer(vertexPos, 2, gl.FLOAT, false, 4 * 4, gl.PtrOffset(0))
   
   textPos   := uint32(gl.GetAttribLocation(this.theProgram, gl.Str("vertex_texture_coords" + "\x00")))
   gl.EnableVertexAttribArray(textPos)
   gl.VertexAttribPointer(textPos, 2, gl.FLOAT, false, 4 * 4, gl.PtrOffset(2 * 4))
   
   // Vertex Shader Uniforms
   this.zoomUni      = gl.GetUniformLocation(this.theProgram, gl.Str("zoom_mat" + "\x00"))
   this.rotUni       = gl.GetUniformLocation(this.theProgram, gl.Str("rot_mat" + "\x00"))
   this.cameraUni    = gl.GetUniformLocation(this.theProgram, gl.Str("camera_pos" + "\x00"))

   // Clean up
   gl.BindBuffer(gl.ARRAY_BUFFER, 0)
   gl.BindVertexArray(0)
}

// Fills a buffer with a square mesh
func uploadSquareMesh() uint32 {
   
   theSquare := []float32 {
      // X  Y  U  V
      
      // TOP LEFT
      -1.0,  1.0, 0.0, 0.0,
      // TOP RIGHT
       1.0,  1.0, 1.0, 0.0,
      // BOTTOM LEFT
      -1.0, -1.0, 0.0, 1.0,
      
      // BOTTOM LEFT
      -1.0, -1.0, 0.0, 1.0,
      // TOP RIGHT
       1.0,  1.0, 1.0, 0.0,
      // BOTTOM RIGHT
       1.0, -1.0, 1.0, 1.0}
   
   var theVBO uint32
   gl.GenBuffers(1, &theVBO)
   gl.BindBuffer(gl.ARRAY_BUFFER, theVBO)
   gl.BufferData(gl.ARRAY_BUFFER, len(theSquare) * 4, gl.Ptr(theSquare), gl.STATIC_DRAW)
   gl.BindBuffer(gl.ARRAY_BUFFER, 0)
   
   return theVBO
}

// Creates a very simple shader program
func createProgram() uint32 {
   theProgram := gl.CreateProgram()
   
   vertexSrc := `#version 150
   
                 uniform mat4 zoom_mat;
                 uniform mat4 rot_mat;
                 uniform vec2 camera_pos;

                 in float depth;
                 in vec2 world_pos;
                 in float scale;
                 in float angle;
                 in vec2 mesh_pos;
                 in vec2 vertex_texture_coords;
                 
                 out vec2 fragment_texture_coords;

                 void main()
                 {   // Forward the texture coords
                     fragment_texture_coords = vertex_texture_coords;
                                          
                     // Scale
                     vec2 temp = mesh_pos * scale;
                     
                     // Rotate
                     float sint = sin(angle);
                     float cost = cos(angle);
                     float rotated_x = temp.x * cost - temp.y * sint;
                     float rotated_y = temp.x * sint + temp.y * cost;                     
                     temp = vec2(rotated_x, rotated_y);                     
                     
                     // Translate
                     temp += world_pos;
                     temp -= camera_pos;
                     
                     gl_Position = rot_mat * zoom_mat * vec4(temp, depth, 1.0);
                 }`

   fragSrc := `#version 150
   
               uniform sampler2D theTexture;

               in vec2 fragment_texture_coords;
               out vec4 outputColor;

               void main()
               {
                   outputColor = texture(theTexture, fragment_texture_coords);
               }`

   // vertex
   vertexShader, err := compileShader(vertexSrc, gl.VERTEX_SHADER)
   if err != nil {
      fmt.Printf("Failed to compile Vertex Shader!")
   }

   // fragment
   fragmentShader, err :=  compileShader(fragSrc, gl.FRAGMENT_SHADER)
   if err != nil {
      fmt.Printf("Failed to compile Fragment Shader!")
   }
   
   // program linking
   gl.AttachShader(theProgram, vertexShader)
   gl.AttachShader(theProgram, fragmentShader)
   gl.LinkProgram(theProgram)
   gl.DeleteShader(vertexShader)
   gl.DeleteShader(fragmentShader)
   
   var status int32
	gl.GetProgramiv(theProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
	   fmt.Printf("Failed to link the program!")
	}
   
   return theProgram
}

// Compiles a GLSL shader of some type
func compileShader(src string, shaderType uint32) (uint32, error) {
   shader := gl.CreateShader(shaderType)
   
   glsrc := gl.Str(src + "\x00")
   gl.ShaderSource(shader, 1, &glsrc, nil)
   gl.CompileShader(shader)
   
   var status int32
   gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
   
   if status != gl.TRUE {
      return 0, fmt.Errorf("Failed to compile %s", src)
   }
   
   return shader, nil
}
