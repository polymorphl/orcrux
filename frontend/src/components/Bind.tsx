import { useState } from "react";
import { Recompose as RecomposeFn, UploadFile as UploadFileFn } from "../../wailsjs/go/main/App";
import { motion } from "framer-motion";

import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { Label } from "./ui/label";
import { RecomposeResult } from "../types/core";
import { bindVariants } from "../lib/motions";
import { Input } from "./ui/input";
import BindManualController from "./BindManualController";

export default function Bind() {
  const [shards, setShards] = useState(["", ""])
  const [result, setResult] = useState<RecomposeResult>({ error: null, data: null })

  const onReset = () => {
    setResult({ error: null, data: null })
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

  const onUpload = async () => {
    const fileContent = await UploadFileFn()
    if (!fileContent) return
    const shards = fileContent.split('\n').filter(line => line.trim() !== '')
    setShards(shards)
    await onRecompose()
  }

  return (
    <motion.div
      variants={bindVariants.container}
      initial="hidden"
      animate="visible"
      className="flex flex-col items-center justify-between gap-3 p-4"
    >
      <div className="flex items-center gap-2">
        <Button variant="outline" size="sm" onClick={onUpload}>
          Upload
        </Button>
        <p className="text-sm text-crystal-200">This will upload the shards from your local machine.</p>
      </div>
      <div className="grid grid-cols-2 gap-6">
        {/* Left Column - Controls and Shards */}
        <div className="flex flex-col">
          <BindManualController shards={shards} onAdd={() => setShards([...shards, ""])} onRemove={() => setShards(shards.slice(0, -1))} onReset={onReset} />

          <motion.div variants={bindVariants.item} className="grid grid-cols-1 gap-3 mt-4 overflow-y-scroll max-h-[135px] w-full">
            {shards.map((shard, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.3, delay: i * 0.1 }}
                className="relative w-full"
              >
                <div className="flex items-start gap-3">
                  <Label className="text-sm font-medium mt-2 flex-shrink-0" htmlFor={`shard-${i}`}>
                    Shard {i + 1}
                  </Label>
                  <Input
                    id={`shard-${i}`}
                    placeholder={`Paste shard ${i + 1} here...`}
                    value={shard}
                    onChange={(e) => setShards(shards.map((s, j) => j === i ? e.target.value : s))}
                    className="flex-1"
                  />
                </div>
              </motion.div>
            ))}
          </motion.div>

          <motion.div variants={bindVariants.item} className="mt-4">
            <motion.div variants={bindVariants.button} whileHover="hover" whileTap="tap">
              <Button onClick={onRecompose} disabled={shards.length < 2}>
                Recompose
              </Button>
            </motion.div>
          </motion.div>
        </div>

        {/* Right Column - Results */}
        <div className="flex flex-col">
          {result.data || result.error ? (
            <motion.div
              variants={bindVariants.result}
              initial="hidden"
              animate="visible"
              className="h-full"
            >
              {result.data && <Textarea value={result.data} readOnly className="h-full min-h-[200px] resize-none" />}
              {result.error && <p className="text-red-500">{result.error}</p>}
            </motion.div>
          ) : (
            <div className="flex items-center justify-center h-full text-crystal-400 text-sm">
              Results will appear here after recomposing
            </div>
          )}
        </div>
      </div>
    </motion.div>
  );
}