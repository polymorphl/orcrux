import { useState } from "react";
import { Split as SplitFn } from "../../wailsjs/go/main/App";

import SplitResults from "./SplitResults";
import SplitForm from "./SplitForm";
import { SplitResult } from "../types/core";

export default function Split() {
  const [step, setStep] = useState<number>(0)
  const [result, setResult] = useState<SplitResult>({ error: null, data: null } as SplitResult)

  const handleSplit = async (secret: string, shards: number, shardsNeeded: number, output: string) => {
    setResult({ error: null, data: null })
    const result = await SplitFn(secret, shards, shardsNeeded, output)
    const parsedResult = JSON.parse(result) as SplitResult
    setResult(parsedResult)
    setStep(1)
  }

  const handleDownload = (data: string) => {
    const blob = new Blob([data], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'shards.txt'
    a.click()
  }

  return (
    <div className="flex flex-col flex items-center justify-between gap-3 p-4">
      {step === 0 && <SplitForm onSplit={handleSplit} />}
      {step === 1 && <SplitResults results={result} onBack={() => setStep(0)} onDownload={handleDownload} />}
    </div>
  )
}