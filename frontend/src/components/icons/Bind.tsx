import * as React from "react";
import type { SVGProps } from "react";
const SvgBind = (props: SVGProps<SVGSVGElement>) => (
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
      d="m21 12-3-3m3 3-3 3m3-3h-5.929a5 5 0 0 0-3.535 1.464l-1.072 1.072A5 5 0 0 1 6.93 16H3m0-8h4.343a4 4 0 0 1 2.829 1.172l1.656 1.656A4 4 0 0 0 14.657 12H18"
    />
  </svg>
);
export default SvgBind;
