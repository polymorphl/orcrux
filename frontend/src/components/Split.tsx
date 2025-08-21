import { useState } from "react";
import { Split as SplitFn } from "../../wailsjs/go/main/App";
import { motion } from "framer-motion";

import { Label } from "./ui/label";
import { RadioGroup, RadioGroupItem } from "./ui/radio-group";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { Input } from "./ui/input";

const MAX_SHARDS = 255

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

export default function Split() {
  const [secret, setSecret] = useState<string>('test')
  const [shards, setShards] = useState<number>(2)
  const [shardsNeeded, setShardsNeeded] = useState<number>(2)
  const [output, setOutput] = useState<'base64' | 'hex'>('base64')
  const [result, setResult] = useState<string>('')

  const handleSplit = async () => {
    console.log({ secret, shards, shardsNeeded, output })
    const result = await SplitFn(secret, shards, shardsNeeded, output)
    setResult(result)
  }

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="w-full"
    >
      <motion.div variants={itemVariants} className="grid w-full max-w-sm items-center gap-3">
        <Label htmlFor="secret">Secret</Label>
        <Textarea id="secret" value={secret} onChange={(e) => setSecret(e.target.value)} placeholder="Enter your secret here..." />
      </motion.div>

      <motion.div variants={itemVariants} className="grid grid-cols-2 gap-3 mt-8">
        <Label htmlFor="shards">Shards</Label>
        <Input type="number" min={2} max={MAX_SHARDS} value={shards} onChange={(e) => setShards(Number(e.target.value))} />
        <Label htmlFor="shardsNeeded">Shards Needed</Label>
        <Input type="number" min={2} max={shards} value={shardsNeeded} onChange={(e) => setShardsNeeded(Number(e.target.value))} />
      </motion.div>

      <motion.div variants={itemVariants} className="grid grid-cols-2 gap-3 mt-3">
        <Label htmlFor="output">Output</Label>
        <RadioGroup defaultValue="base64" onValueChange={(value) => setOutput(value as 'base64' | 'hex')}>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="base64" id="base64" />
            <Label htmlFor="base64" className="cursor-pointer">Base64</Label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="hex" id="hex" />
            <Label htmlFor="hex" className="cursor-pointer">Hex</Label>
          </div>
        </RadioGroup>
      </motion.div>

      <motion.div variants={itemVariants} className="mt-4">
        <motion.div
          variants={buttonVariants}
          whileHover="hover"
          whileTap="tap"
        >
          <Button onClick={handleSplit} disabled={!secret || !shards || !shardsNeeded}>
            Split
          </Button>
        </motion.div>
      </motion.div>

      {result && (
        <motion.div
          variants={resultVariants}
          initial="hidden"
          animate="visible"
          className="grid w-full max-w-sm items-center gap-3 mt-6"
        >
          <Label htmlFor="result">Result</Label>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 max-w-xs">
            {result.split('\n')
              .filter(line => line.trim() !== '')
              .map((line, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1, duration: 0.3 }}
                  className="flex items-center justify-between gap-2 p-3 bg-crystal-700/30 rounded-sm border border-crystal-500/30 hover:bg-crystal-700/40 transition-colors min-w-0"
                >
                  <pre className="text-xs text-crystal-200 font-mono truncate flex-1 min-w-0 overflow-hidden">
                    {line}
                  </pre>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => navigator.clipboard.writeText(line)}
                    className="h-6 px-2 text-crystal-300 hover:text-crystal-100 hover:bg-crystal-600/50 shrink-0 ml-2"
                  >
                    <svg
                      className="w-3 h-3"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                      />
                    </svg>
                  </Button>
                </motion.div>
              ))}
          </div>
        </motion.div>
      )}
    </motion.div>
  )
}