import { useState } from "react";
import { Split as SplitFn } from "../../wailsjs/go/main/App";

import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { RadioGroup, RadioGroupItem } from "./ui/radio-group";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { Icon } from "./Icon";

export default function Split() {
  const [secret, setSecret] = useState<string>('test')
  const [shards, setShards] = useState<number>(2)
  const [shardsNeeded, setShardsNeeded] = useState<number>(2)
  const [output, setOutput] = useState<'base64' | 'hex'>('base64')
  const [result, setResult] = useState<string>('')

  const handleSplit = async () => {
    const result = await SplitFn(secret, shards, shardsNeeded, output)
    setResult(result)
  }

  return (
    <>
      <div className="grid w-full max-w-sm items-center gap-3">
        <Label htmlFor="secret">Secret</Label>
        <Textarea id="secret" value={secret} onChange={(e) => setSecret(e.target.value)} />
        <Label htmlFor="shards">Shards</Label>
        <Input id="shards" type="number" value={shards} onChange={(e) => setShards(Number(e.target.value))} />
        <Label htmlFor="shardsNeeded">Shards Needed</Label>
        <Input id="shardsNeeded" type="number" value={shardsNeeded} onChange={(e) => setShardsNeeded(Number(e.target.value))} />
        <Label htmlFor="output">Output</Label>
        <RadioGroup defaultValue="base64" onValueChange={(value) => setOutput(value as 'base64' | 'hex')}>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="base64" id="base64" />
            <Label htmlFor="base64">Base64</Label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="hex" id="hex" />
            <Label htmlFor="hex">Hex</Label>
          </div>
        </RadioGroup>
        <Button onClick={handleSplit} disabled={!secret || !shards || !shardsNeeded}><Icon icon="Split" />Split</Button>
      </div>
      {result && <div className="grid w-full max-w-sm items-center gap-3">
        <Label htmlFor="result">Result</Label>
        <pre id="result" className="text-sm grid grid-cols-2 gap-2">{result.split('\n').map((line, index) => <div key={index} className="text-xs">{line}</div>)}</pre>
      </div>}
    </>
  )
}