import { motion } from "framer-motion";
import { Label } from "./ui/label";
import { Textarea } from "./ui/textarea";
import { Input } from "./ui/input";
import { RadioGroup, RadioGroupItem } from "./ui/radio-group";
import { Button } from "./ui/button";
import { useState } from "react";

type SplitFormProps = {
  onSplit: (secret: string, shards: number, shardsNeeded: number, output: string) => void;
}

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

const MAX_SHARDS = 255

export default function SplitForm({ onSplit }: SplitFormProps) {
  const [secret, setSecret] = useState<string>('test')
  const [shards, setShards] = useState<number>(2)
  const [shardsNeeded, setShardsNeeded] = useState<number>(2)
  const [output, setOutput] = useState<'base64' | 'hex'>('base64')

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
          <Button onClick={() => onSplit(secret, shards, shardsNeeded, output)} disabled={!secret || !shards || !shardsNeeded}>
            Split
          </Button>
        </motion.div>
      </motion.div>
    </motion.div>
  )
}