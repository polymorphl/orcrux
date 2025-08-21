import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-sm text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive cursor-pointer",
  {
    variants: {
      variant: {
        default:
          "border border-crystal-500 bg-transparent text-crystal-200 shadow-xs hover:bg-crystal-700 hover:text-crystal-100",
        destructive:
          "bg-red-600 text-white shadow-xs hover:bg-red-700 focus-visible:ring-red-600/20 dark:focus-visible:ring-red-600/40 dark:bg-red-700/60",
        outline:
          "border border-crystal-500 bg-transparent text-crystal-200 shadow-xs hover:bg-crystal-700 hover:text-crystal-100",
        secondary:
          "bg-crystal-700 text-crystal-200 shadow-xs hover:bg-crystal-600",
        ghost:
          "text-crystal-200 hover:bg-crystal-700 hover:text-crystal-100",
        link: "text-crystal-200 underline-offset-4 hover:underline hover:text-crystal-100",
        gradient: "bg-gradient-to-r from-crystal-600 to-crystal-500 hover:from-crystal-500 hover:to-crystal-400 text-crystal-100 font-medium shadow-md",
      },
      size: {
        default: "h-9 px-4 py-2 has-[>svg]:px-3",
        sm: "h-8 rounded-sm gap-1.5 px-3 has-[>svg]:px-2.5",
        lg: "h-10 rounded-sm px-6 has-[>svg]:px-4",
        icon: "size-9",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

function Button({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean
  }) {
  const Comp = asChild ? Slot : "button"

  return (
    <Comp
      data-slot="button"
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  )
}

export { Button, buttonVariants }
