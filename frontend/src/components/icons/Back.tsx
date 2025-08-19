import * as React from "react";
import type { SVGProps } from "react";
const SvgBack = (props: SVGProps<SVGSVGElement>) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={800}
    height={800}
    fill="none"
    viewBox="0 0 24 24"
    {...props}
  >
    <path
      stroke="#000"
      strokeLinecap="round"
      strokeLinejoin="round"
      strokeWidth={2}
      d="M11 18h3.75a5.25 5.25 0 1 0 0-10.5H5M7.5 4 4 7.5 7.5 11"
    />
  </svg>
);
export default SvgBack;
