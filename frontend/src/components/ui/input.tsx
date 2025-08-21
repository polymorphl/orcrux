import * as React from "react"

import { cn } from "@/lib/utils"

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        "file:text-crystal-200 placeholder:text-crystal-400 selection:bg-crystal-600 selection:text-crystal-100 bg-crystal-800/30 border-crystal-500 flex h-9 w-full min-w-0 rounded-sm border bg-transparent px-3 py-1 text-base shadow-sm transition-all duration-200 outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm text-crystal-200",
        "focus-visible:border-crystal-400 focus-visible:ring-crystal-400/30 focus-visible:ring-[3px]",
        "aria-invalid:ring-red-500/20 dark:aria-invalid:ring-red-500/40 aria-invalid:border-red-500",
        "hover:border-crystal-400 hover:bg-crystal-800/50",
        className
      )}
      {...props}
    />
  )
}

export { Input }
