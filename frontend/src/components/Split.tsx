import { useState } from "react";
import { Split as SplitFn } from "../../wailsjs/go/main/App";

import SplitResults from "./SplitResults";
import SplitForm from "./SplitForm";
import { SplitResult } from "../types/core";

export default function Split() {
  const [result, setResult] = useState<SplitResult>({ error: null, data: null } as SplitResult)

  const handleSplit = async (secret: string, shards: number, shardsNeeded: number, output: string) => {
    setResult({ error: null, data: null })
    const result = await SplitFn(secret, shards, shardsNeeded, output)
    const parsedResult = JSON.parse(result) as SplitResult
    setResult(parsedResult)
  }

  return (
    <div className="flex flex-col gap-4 flex items-center justify-between gap-3 p-4">
      <SplitForm onSplit={handleSplit} />
      <SplitResults results={result} />
    </div>
  )
}