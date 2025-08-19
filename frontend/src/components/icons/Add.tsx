import * as React from "react";
import type { SVGProps } from "react";
const SvgAdd = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={800}
    height={800}
    viewBox="0 0 24 24"
    {...props}
  >
    <g
      fill="none"
      stroke="#000"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      data-name="add"
    >
      <path d="M12 19V5M5 12h14" />
    </g>
  </svg>
);
export default SvgAdd;
