#version 330 core

layout(location = 0) in vec3 InPosition;
layout(location = 1) in float BlockId;
layout(location = 2) in float faceDirection;

uniform mat4 matProjection;
uniform mat4 matView;
uniform mat4 matModel;

void main() {
    gl_Position = matProjection * matView * matModel * vec4(InPosition, 1.0);
}
