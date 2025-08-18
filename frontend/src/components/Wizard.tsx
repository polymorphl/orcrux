import { useState } from "react";

import Bind from "./Bind";
import Split from "./Split";
import { Button } from "./ui/button";

export default function Wizard() {
  const [workflow, setWorkflow] = useState<"bind" | "split" | null>(null);

  return <div>
    <div>
      {workflow === null ? <div>
        <Button onClick={() => setWorkflow("bind")}>Bind</Button>
        <Button onClick={() => setWorkflow("split")}>Split</Button>
      </div> : <div>
        <Button onClick={() => setWorkflow(null)}>Back</Button>
      </div>}
    </div>
    <div className="flex flex-col gap-4 items-center justify-center">
      {workflow === "bind" && <Bind />}
      {workflow === "split" && <Split />}
    </div>
  </div>;
}