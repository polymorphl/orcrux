import { motion } from "framer-motion";

import { bindVariants } from "@/lib/motions";
import { Button } from "./ui/button";
import { Icon } from "./Icon";

type BindManualControllerProps = {
  shards: string[]
  onAdd: () => void
  onRemove: () => void
  onReset: () => void
}

export default function BindManualController({ shards, onAdd, onRemove, onReset }: BindManualControllerProps) {
  return (
    <motion.div variants={bindVariants.item} className="flex items-center justify-between gap-4 p-4 bg-crystal-700/20 rounded-sm border border-crystal-500/20">
      <div className="flex items-center gap-2">
        <motion.div variants={bindVariants.button} whileHover="hover" whileTap="tap">
          <Button
            variant="outline"
            size="sm"
            onClick={onRemove}
            disabled={shards.length <= 2}
            className="h-7 px-2"
          >
            <Icon icon="Remove" />
          </Button>
        </motion.div>
        <motion.div variants={bindVariants.button} whileHover="hover" whileTap="tap">
          <Button
            variant="outline"
            size="sm"
            onClick={onAdd}
            disabled={shards.length >= 255}
            className="h-7 px-2"
          >
            <Icon icon="Add" />
          </Button>
        </motion.div>
        <span className="text-crystal-300 text-sm font-medium">
          {shards.length} shard{shards.length !== 1 ? 's' : ''}
        </span>
      </div>

      <motion.div variants={bindVariants.button} whileHover="hover" whileTap="tap">
        <Button
          variant="default"
          size="sm"
          onClick={onReset}
          className="h-7 px-3 text-crystal-300 hover:text-crystal-100"
        >
          Reset
        </Button>
      </motion.div>
    </motion.div>
  )
}