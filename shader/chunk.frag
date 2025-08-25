#version 330 core

layout(location = 0) out vec4 fragColor;

const vec3 gamma = vec3(2.2);
const vec3 inv_gamma = 1 / gamma;

uniform sampler2D texture0;

in vec3 voxel_color;
in vec2 uv;
in float shading;
flat in int blockId;
// variant can be 0,1,2
// 0 for top, 1 for bottom, 2 for side
flat in int variant;

void main() {
    vec2 faceUv = uv;
    ivec2 texSize = textureSize(texture0, 0);

    int atlasWidthInCells = 3;
    // number of block types
    int atlasHeightInCells = texSize.y / 16;

    float cellWidth = 1.0 / atlasWidthInCells;
    float cellHeight = 1.0 / atlasHeightInCells;

    int column = variant;
    int row = blockId;

    // base UV (bottom-left of tile)
    vec2 tileOrigin = vec2(column * cellWidth, row * cellHeight);

    // final UV = origin + local uv inside the cell
    faceUv = tileOrigin + uv * vec2(cellWidth, cellHeight);

    vec3 tex_col = texture(texture0, faceUv).rgb;
    tex_col = pow(tex_col, gamma);

    tex_col *= shading;

    tex_col = pow(tex_col, inv_gamma);
    fragColor = vec4(tex_col, 1);
}
