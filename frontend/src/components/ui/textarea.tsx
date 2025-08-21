import * as React from "react"

import { cn } from "@/lib/utils"

function Textarea({ className, ...props }: React.ComponentProps<"textarea">) {
  return (
    <textarea
      data-slot="textarea"
      className={cn(
        "border-crystal-500 placeholder:text-crystal-400 focus-visible:border-crystal-400 focus-visible:ring-crystal-400/30 aria-invalid:ring-red-500/20 dark:aria-invalid:ring-red-500/40 aria-invalid:border-red-500 bg-crystal-800/30 flex field-sizing-content min-h-16 w-full rounded-sm border bg-transparent px-3 py-2 text-base shadow-sm transition-all duration-200 outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm text-crystal-200 hover:border-crystal-400 hover:bg-crystal-800/50",
        className
      )}
      {...props}
    />
  )
}

export { Textarea }
