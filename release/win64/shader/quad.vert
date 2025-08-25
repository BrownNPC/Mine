#version 330 core
#extension GL_ARB_separate_shader_objects : enable

layout(location = 0) in vec3 in_position;
layout(location = 1) in vec3 in_color;

uniform mat4 mvp;

out vec3 color;

void main() {
    color = in_color;
    // gl_Position = matProj * matView * matModel * vec4(in_position, 1.0);
    gl_Position = mvp * vec4(in_position, 1.0);
}
