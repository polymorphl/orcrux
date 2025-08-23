import { useState } from "react";
import { Split as SplitFn } from "../../wailsjs/go/main/App";

import SplitResults from "./SplitResults";
import SplitForm from "./SplitForm";

export default function Split() {
  const [result, setResult] = useState<string>('')

  const handleSplit = async (secret: string, shards: number, shardsNeeded: number, output: string) => {
    setResult('')
    console.log('#handleSplit', { secret, shards, shardsNeeded, output })
    const result = await SplitFn(secret, shards, shardsNeeded, output)
    setResult(result)
  }

  return (
    <div className="flex flex-col gap-4 flex items-center justify-between gap-3 p-4 bg-gradient-to-br from-crystal-700/40 to-crystal-600/30 rounded-sm border border-crystal-500/20 shadow-lg backdrop-blur-sm">
      <SplitForm onSplit={handleSplit} />
      <SplitResults results={result} />
    </div>
  )
}