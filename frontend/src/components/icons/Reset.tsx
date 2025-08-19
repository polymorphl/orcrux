import * as React from "react";
import type { SVGProps } from "react";
const SvgReset = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={800}
    height={800}
    fill="none"
    viewBox="0 0 32 32"
    {...props}
  >
    <path
      fill="currentColor"
      d="M27 8H6.83l3.58-3.59L9 3 3 9l6 6 1.41-1.41L6.83 10H27v16H7v-7H5v7a2 2 0 0 0 2 2h20a2 2 0 0 0 2-2V10a2 2 0 0 0-2-2"
    />
  </svg>
);
export default SvgReset;
