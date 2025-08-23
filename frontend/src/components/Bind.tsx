import { useState } from "react";
import { Recompose as RecomposeFn } from "../../wailsjs/go/main/App";
import { motion } from "framer-motion";

import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { Icon } from "./Icon";
import { Label } from "./ui/label";
import { RecomposeResult } from "../types/core";

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.2
    }
  }
};

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.4, ease: "easeOut" as const }
  }
};

const buttonVariants = {
  hover: {
    scale: 1.02,
    transition: { duration: 0.2, ease: "easeInOut" as const }
  },
  tap: {
    scale: 0.98,
    transition: { duration: 0.1 }
  }
};

const resultVariants = {
  hidden: { opacity: 0, height: 0 },
  visible: {
    opacity: 1,
    height: "auto",
    transition: { duration: 0.5, ease: "easeOut" as const }
  }
};

export default function Bind() {
  const [shards, setShards] = useState(["", ""])
  const [result, setResult] = useState<RecomposeResult>({ error: null, data: null })

  const onReset = () => {
    setShards(["", ""])
  }

  const onRecompose = async () => {
    setResult({ error: null, data: null })
    if (shards.length < 2 || shards.some(shard => shard === "")) {
      return
    }
    const result = await RecomposeFn(shards)
    const parsedResult = JSON.parse(result) as RecomposeResult
    setResult(parsedResult)
  }

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="flex flex-col items-center justify-between gap-3 p-4"
    >
      <motion.div variants={itemVariants} className="flex items-center justify-between gap-4 p-4 bg-crystal-700/20 rounded-sm border border-crystal-500/20">
        <div className="flex items-center gap-2">
          <motion.div variants={buttonVariants} whileHover="hover" whileTap="tap">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShards(shards.slice(0, -1))}
              disabled={shards.length <= 2}
              className="h-8 px-3"
            >
              <Icon icon="Remove" />
            </Button>
          </motion.div>
          <motion.div variants={buttonVariants} whileHover="hover" whileTap="tap">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShards([...shards, ""])}
              disabled={shards.length >= 255}
              className="h-8 px-3"
            >
              <Icon icon="Add" />
            </Button>
          </motion.div>
          <span className="text-crystal-300 text-sm font-medium">
            {shards.length} shard{shards.length !== 1 ? 's' : ''}
          </span>
        </div>

        <motion.div variants={buttonVariants} whileHover="hover" whileTap="tap">
          <Button
            variant="default"
            size="sm"
            onClick={onReset}
            className="h-8 px-3 text-crystal-300 hover:text-crystal-100"
          >
            Reset
          </Button>
        </motion.div>
      </motion.div>

      <motion.div variants={itemVariants} className="grid grid-cols-1 gap-3 mt-4 overflow-y-scroll max-h-[200px] w-full">
        {shards.map((shard, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.3, delay: i * 0.1 }}
            className="relative w-1/2 mx-auto"
          >
            <div className="flex items-center justify-between gap-2 mb-2">
              <Label className="text-sm font-medium">
                Shard {i + 1}
              </Label>
              <span className="text-xs text-muted-foreground">
                {shard.length} chars
              </span>
            </div>
            <Textarea
              placeholder={`Paste shard ${i + 1} here...`}
              value={shard}
              onChange={(e) => setShards(shards.map((s, j) => j === i ? e.target.value : s))}
              className="min-h-[80px] max-h-[120px] resize-none"
            />
          </motion.div>
        ))}
      </motion.div>

      <motion.div variants={itemVariants} className="mt-4">
        <motion.div variants={buttonVariants} whileHover="hover" whileTap="tap">
          <Button onClick={onRecompose} disabled={shards.length < 2}>
            Recompose
          </Button>
        </motion.div>
      </motion.div>

      {result.data || result.error && (
        <motion.div
          variants={resultVariants}
          initial="hidden"
          animate="visible"
          className="mt-4"
        >
          {result.data && <Textarea value={result.data} readOnly />}
          {result.error && <p className="text-red-500">{result.error}</p>}
        </motion.div>
      )}
    </motion.div>
  );
}