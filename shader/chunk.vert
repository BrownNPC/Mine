#version 330 core

layout(location = 0) in ivec3 InPosition;
layout(location = 1) in int BlockId;
layout(location = 2) in int faceDirection;

uniform mat4 matModel;
uniform mat4 matView;
uniform mat4 matProjection;

out vec3 voxel_color;
out vec2 uv;
out float shading;

const float face_shading[6] = float[6](
        1.0, 0.5, // top bottom
        0.5, 0.8, // right left
        0.5, 0.8 // front back
    );

const vec2 uv_coords[4] = vec2[4](
        vec2(0, 0), vec2(0, 1),
        vec2(1, 0), vec2(1, 1)
    );

const int uv_indices[12] = int[12](
        1, 0, 2, 1, 2, 3, // tex coords indices for vertices of an even face
        3, 0, 2, 3, 1, 0 // odd face
    );

vec3 hash31(float p) {
    vec3 p3 = fract(vec3(p * 21.2) * vec3(0.1031, 0.1030, 0.0973));
    p3 += dot(p3, p3.yzx + 33.33);
    return fract((p3.xxy + p3.yzz) * p3.zyx) + 0.05;
}

void main() {
    int uv_index = gl_VertexID % 6 + (faceDirection & 1) * 6;
    uv = uv_coords[uv_indices[uv_index]];
    voxel_color = hash31(BlockId);

    shading = face_shading[faceDirection];
    mat4 mvp = matProjection * matView * matModel;
    gl_Position = mvp * vec4(InPosition, 1.0);
}
