import { Renderer, Program, Mesh, Triangle } from "ogl";
import { useEffect, useRef } from "react";

interface BalatroProps {
  spinRotation?: number;
  spinSpeed?: number;
  offset?: [number, number];
  color1?: string;
  color2?: string;
  color3?: string;
  contrast?: number;
  lighting?: number;
  spinAmount?: number;
  pixelFilter?: number;
  spinEase?: number;
  isRotate?: boolean;
  mouseInteraction?: boolean;
  transitionDuration?: number; // Duration in seconds for color transitions
}

function hexToVec4(hex: string): [number, number, number, number] {
  let hexStr = hex.replace("#", "");
  let r = 0,
    g = 0,
    b = 0,
    a = 1;
  if (hexStr.length === 6) {
    r = parseInt(hexStr.slice(0, 2), 16) / 255;
    g = parseInt(hexStr.slice(2, 4), 16) / 255;
    b = parseInt(hexStr.slice(4, 6), 16) / 255;
  } else if (hexStr.length === 8) {
    r = parseInt(hexStr.slice(0, 2), 16) / 255;
    g = parseInt(hexStr.slice(2, 4), 16) / 255;
    b = parseInt(hexStr.slice(4, 6), 16) / 255;
    a = parseInt(hexStr.slice(6, 8), 16) / 255;
  }
  return [r, g, b, a];
}

const vertexShader = `
attribute vec2 uv;
attribute vec2 position;
varying vec2 vUv;
void main() {
  vUv = uv;
  gl_Position = vec4(position, 0, 1);
}
`;

const fragmentShader = `
precision highp float;

#define PI 3.14159265359

uniform float iTime;
uniform vec3 iResolution;
uniform float uSpinRotation;
uniform float uSpinSpeed;
uniform vec2 uOffset;
uniform vec4 uColor1;
uniform vec4 uColor2;
uniform vec4 uColor3;
uniform float uContrast;
uniform float uLighting;
uniform float uSpinAmount;
uniform float uPixelFilter;
uniform float uSpinEase;
uniform bool uIsRotate;
uniform vec2 uMouse;

varying vec2 vUv;

vec4 effect(vec2 screenSize, vec2 screen_coords) {
    float pixel_size = length(screenSize.xy) / uPixelFilter;
    vec2 uv = (floor(screen_coords.xy * (1.0 / pixel_size)) * pixel_size - 0.5 * screenSize.xy) / length(screenSize.xy) - uOffset;
    float uv_len = length(uv);
    
    float speed = (uSpinRotation * uSpinEase * 0.2);
    if(uIsRotate){
       speed = iTime * speed;
    }
    speed += 302.2;
    
    float mouseInfluence = (uMouse.x * 2.0 - 1.0);
    speed += mouseInfluence * 0.1;
    
    float new_pixel_angle = atan(uv.y, uv.x) + speed - uSpinEase * 20.0 * (uSpinAmount * uv_len + (1.0 - uSpinAmount));
    vec2 mid = (screenSize.xy / length(screenSize.xy)) / 2.0;
    uv = (vec2(uv_len * cos(new_pixel_angle) + mid.x, uv_len * sin(new_pixel_angle) + mid.y) - mid);
    
    uv *= 30.0;
    float baseSpeed = iTime * uSpinSpeed;
    speed = baseSpeed + mouseInfluence * 2.0;
    
    vec2 uv2 = vec2(uv.x + uv.y);
    
    for(int i = 0; i < 5; i++) {
        uv2 += sin(max(uv.x, uv.y)) + uv;
        uv += 0.5 * vec2(
            cos(5.1123314 + 0.353 * uv2.y + speed * 0.131121),
            sin(uv2.x - 0.113 * speed)
        );
        uv -= cos(uv.x + uv.y) - sin(uv.x * 0.711 - uv.y);
    }
    
    float contrast_mod = (0.25 * uContrast + 0.5 * uSpinAmount + 1.2);
    float paint_res = min(2.0, max(0.0, length(uv) * 0.035 * contrast_mod));
    float c1p = max(0.0, 1.0 - contrast_mod * abs(1.0 - paint_res));
    float c2p = max(0.0, 1.0 - contrast_mod * abs(paint_res));
    float c3p = 1.0 - min(1.0, c1p + c2p);
    float light = (uLighting - 0.2) * max(c1p * 5.0 - 4.0, 0.0) + uLighting * max(c2p * 5.0 - 4.0, 0.0);
    
    return (0.3 / uContrast) * uColor1 + (1.0 - 0.3 / uContrast) * (uColor1 * c1p + uColor2 * c2p + vec4(c3p * uColor3.rgb, c3p * uColor1.a)) + light;
}

void main() {
    vec2 uv = vUv * iResolution.xy;
    gl_FragColor = effect(iResolution.xy, uv);
}
`;

export default function Balatro({
  spinRotation = -2.0,
  spinSpeed = 7.0,
  offset = [0.0, 0.0],
  color1 = "#DE443B",
  color2 = "#006BB4",
  color3 = "#162325",
  contrast = 3.5,
  lighting = 0.4,
  spinAmount = 0.25,
  pixelFilter = 745.0,
  spinEase = 1.0,
  isRotate = false,
  mouseInteraction = true,
  transitionDuration = 0.8,
}: BalatroProps) {
  const containerRef = useRef<HTMLDivElement>(null);

  // Color transition state
  const prevColorsRef = useRef({ color1, color2, color3 });
  const currentColorsRef = useRef({ color1, color2, color3 });
  const transitionStartTimeRef = useRef(0);
  const isTransitioningRef = useRef(false);

  // Effect to handle color changes and trigger transitions
  useEffect(() => {
    const currentColors = { color1, color2, color3 };
    const prevColors = prevColorsRef.current;

    // Check if any colors have changed
    if (prevColors.color1 !== color1 || prevColors.color2 !== color2 || prevColors.color3 !== color3) {
      // Start transition
      prevColorsRef.current = { ...prevColors };
      currentColorsRef.current = { ...currentColors };
      transitionStartTimeRef.current = performance.now();
      isTransitioningRef.current = true;
    }
  }, [color1, color2, color3]);

  useEffect(() => {
    if (!containerRef.current) return;
    const container = containerRef.current;
    const renderer = new Renderer();
    const gl = renderer.gl;
    gl.clearColor(0, 0, 0, 1);

    let program: Program;

    function resize() {
      renderer.setSize(container.offsetWidth, container.offsetHeight);
      if (program) {
        program.uniforms.iResolution.value = [
          gl.canvas.width,
          gl.canvas.height,
          gl.canvas.width / gl.canvas.height,
        ];
      }
    }
    window.addEventListener("resize", resize);
    resize();

    const geometry = new Triangle(gl);
    program = new Program(gl, {
      vertex: vertexShader,
      fragment: fragmentShader,
      uniforms: {
        iTime: { value: 0 },
        iResolution: {
          value: [
            gl.canvas.width,
            gl.canvas.height,
            gl.canvas.width / gl.canvas.height,
          ],
        },
        uSpinRotation: { value: spinRotation },
        uSpinSpeed: { value: spinSpeed },
        uOffset: { value: offset },
        uColor1: { value: hexToVec4(color1) },
        uColor2: { value: hexToVec4(color2) },
        uColor3: { value: hexToVec4(color3) },
        uContrast: { value: contrast },
        uLighting: { value: lighting },
        uSpinAmount: { value: spinAmount },
        uPixelFilter: { value: pixelFilter },
        uSpinEase: { value: spinEase },
        uIsRotate: { value: isRotate },
        uMouse: { value: [0.5, 0.5] },
      },
    });

    const mesh = new Mesh(gl, { geometry, program });
    let animationFrameId: number;

    // Color interpolation function
    function interpolateColors(prev: string, current: string, progress: number): [number, number, number, number] {
      const prevVec = hexToVec4(prev);
      const currentVec = hexToVec4(current);

      return [
        prevVec[0] + (currentVec[0] - prevVec[0]) * progress,
        prevVec[1] + (currentVec[1] - prevVec[1]) * progress,
        prevVec[2] + (currentVec[2] - prevVec[2]) * progress,
        prevVec[3] + (currentVec[3] - prevVec[3]) * progress,
      ];
    }

    function update(time: number) {
      animationFrameId = requestAnimationFrame(update);

      // Handle color transitions
      if (isTransitioningRef.current) {
        const elapsed = (time - transitionStartTimeRef.current) * 0.001;
        const progress = Math.min(elapsed / transitionDuration, 1.0);

        if (progress >= 1.0) {
          // Transition complete
          isTransitioningRef.current = false;
          prevColorsRef.current = { ...currentColorsRef.current };
        } else {
          // Interpolate colors
          const easedProgress = 1 - Math.pow(1 - progress, 3); // Ease-out cubic
          program.uniforms.uColor1.value = interpolateColors(
            prevColorsRef.current.color1,
            currentColorsRef.current.color1,
            easedProgress
          );
          program.uniforms.uColor2.value = interpolateColors(
            prevColorsRef.current.color2,
            currentColorsRef.current.color2,
            easedProgress
          );
          program.uniforms.uColor3.value = interpolateColors(
            prevColorsRef.current.color3,
            currentColorsRef.current.color3,
            easedProgress
          );
        }
      }

      program.uniforms.iTime.value = time * 0.001;
      renderer.render({ scene: mesh });
    }
    animationFrameId = requestAnimationFrame(update);
    container.appendChild(gl.canvas);

    function handleMouseMove(e: MouseEvent) {
      if (!mouseInteraction) return;
      const rect = container.getBoundingClientRect();
      const x = (e.clientX - rect.left) / rect.width;
      const y = 1.0 - (e.clientY - rect.top) / rect.height;
      program.uniforms.uMouse.value = [x, y];
    }
    container.addEventListener("mousemove", handleMouseMove);

    return () => {
      cancelAnimationFrame(animationFrameId);
      window.removeEventListener("resize", resize);
      container.removeEventListener("mousemove", handleMouseMove);
      container.removeChild(gl.canvas);
      gl.getExtension("WEBGL_lose_context")?.loseContext();
    };
  }, [
    spinRotation,
    spinSpeed,
    offset,
    color1,
    color2,
    color3,
    contrast,
    lighting,
    spinAmount,
    pixelFilter,
    spinEase,
    isRotate,
    mouseInteraction,
    transitionDuration
  ]);

  return <div ref={containerRef} className="w-full h-full" />;
}
