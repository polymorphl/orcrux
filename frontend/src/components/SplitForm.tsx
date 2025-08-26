import { motion } from "framer-motion";
import { useState } from "react";

import { Label } from "./ui/label";
import { Textarea } from "./ui/textarea";
import { RadioGroup, RadioGroupItem } from "./ui/radio-group";
import { Button } from "./ui/button";
import { SplitFormProps } from "../types/core";
import ShardsSlider from "./ShardsSlider";
import { splitFormVariants } from "../lib/motions";

const MIN_SHARDS = 2
const MAX_SHARDS = 255

export default function SplitForm({ onSplit }: SplitFormProps) {
  const [secret, setSecret] = useState<string>('')
  const [shards, setShards] = useState<number>(MIN_SHARDS)
  const [shardsNeeded, setShardsNeeded] = useState<number>(MIN_SHARDS)
  const [output, setOutput] = useState<'base64' | 'hex'>('base64')

  return (
    <motion.div
      variants={splitFormVariants.container}
      initial="hidden"
      animate="visible"
      className="w-full"
    >
      <motion.div variants={splitFormVariants.item} className="grid w-full items-center gap-3">
        <Label htmlFor="secret">Secret</Label>
        <Textarea id="secret" value={secret} onChange={(e) => setSecret(e.target.value)} placeholder="Enter your secret here..." className="max-h-[120px] w-full" />
      </motion.div>

      <motion.div variants={splitFormVariants.item} className="grid grid-cols-3 gap-6 mt-4">
        <ShardsSlider label="Total Shards" value={shards} min={MIN_SHARDS} max={MAX_SHARDS} onChange={(value) => setShards(value)} />
        <ShardsSlider label="Shards Needed" value={shardsNeeded} min={MIN_SHARDS} max={shards} onChange={(value) => setShardsNeeded(value)} />
        <div className="flex flex-col gap-2">
          <Label htmlFor="output">Output</Label>
          <RadioGroup defaultValue="base64" onValueChange={(value) => setOutput(value as 'base64' | 'hex')}>
            <div className="flex items-center space-x-6">
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="base64" id="base64" />
                <Label htmlFor="base64" className="cursor-pointer">Base64</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="hex" id="hex" />
                <Label htmlFor="hex" className="cursor-pointer">Hex</Label>
              </div>
            </div>
          </RadioGroup>
        </div>
      </motion.div>

      <motion.div variants={splitFormVariants.item} className="mt-4">
        <motion.div
          variants={splitFormVariants.button}
          whileHover="hover"
          whileTap="tap"
        >
          <Button onClick={() => onSplit(secret, shards, shardsNeeded, output)} disabled={!secret || !shards || !shardsNeeded}>
            Split
          </Button>
        </motion.div>
      </motion.div>
    </motion.div>
  )
}