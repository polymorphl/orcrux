import { useState } from "react";
import { Split as SplitFn, SaveFileDialog as SaveFileDialogFn } from "../../wailsjs/go/main/App";

import SplitResults from "./SplitResults";
import SplitForm from "./SplitForm";
import { SplitResult } from "../types/core";
import { splitActiveColors, splitIdleColors } from "@/lib/colors";

export default function Split() {
  const [step, setStep] = useState<number>(0)
  const [result, setResult] = useState<SplitResult>({ error: null, data: null } as SplitResult)

  const handleSplit = async (secret: string, shards: number, shardsNeeded: number, output: string) => {
    setResult({ error: null, data: null })
    const result = await SplitFn(secret, shards, shardsNeeded, output)
    const parsedResult = JSON.parse(result) as SplitResult
    setResult(parsedResult)
    setStep(1)
    window.parent.postMessage({ type: 'color-change', color1: splitActiveColors[0], color2: splitActiveColors[1] }, '*')
  }

  const handleDownload = async (data: string) => {
    const blob = new Blob([data], { type: 'text/plain' })
    await SaveFileDialogFn(Array.from(new Uint8Array(await blob.arrayBuffer())), "shards.txt")
  }

  const handleBack = () => {
    setStep(0)
    window.parent.postMessage({ type: 'color-change', color1: splitIdleColors[0], color2: splitIdleColors[1] }, '*')
  }

  return (
    <div className="flex flex-col flex items-center justify-between gap-3 p-4">
      {step === 0 && <SplitForm onSplit={handleSplit} />}
      {step === 1 && <SplitResults results={result} onBack={handleBack} onDownload={handleDownload} />}
    </div>
  )
}