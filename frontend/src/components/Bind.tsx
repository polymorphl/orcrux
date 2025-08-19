import { useState } from "react";
import { Recompose as RecomposeFn } from "../../wailsjs/go/main/App";

import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { Icon } from "./Icon";

export default function Bind() {
  const [shards, setShards] = useState(["", ""])
  const [result, setResult] = useState<string>("")

  const onReset = () => {
    setShards(["", ""])
  }

  const onRecompose = async () => {
    if (shards.length < 2 || shards.some(shard => shard === "")) {
      return
    }
    const result = await RecomposeFn(shards)
    setResult(result)
  }

  return <div className="mt-4">
    <Button onClick={() => setShards(shards.slice(0, -1))} disabled={shards.length <= 2}><Icon icon="Remove" />Remove shard</Button>
    <Button onClick={onReset}><Icon icon="Reset" />Reset</Button>
    <Button onClick={() => setShards([...shards, ""])} disabled={shards.length >= 255}><Icon icon="Add" />Add shard</Button>
    <div className="grid grid-cols-2 gap-2 mt-4">
      {shards.map((shard, i) => (
        <div key={i}>
          <Textarea placeholder={`Shard ${i + 1}`} value={shard} onChange={(e) => setShards(shards.map((s, j) => j === i ? e.target.value : s))} />
        </div>
      ))}
    </div>
    <Button onClick={onRecompose} className="mt-4" disabled={shards.length < 2}><Icon icon="Bind" />Recompose</Button>
    {result && <div className="mt-4">
      <Textarea value={result} readOnly />
    </div>}
  </div>;
}