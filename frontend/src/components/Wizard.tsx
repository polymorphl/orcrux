import { useState } from "react";

import Bind from "./Bind";
import Split from "./Split";
import { Button } from "./ui/button";
import { Icon } from "./Icon";

export default function Wizard() {
  const [workflow, setWorkflow] = useState<"bind" | "split" | null>(null);

  return <div>
    <div>
      {workflow === null ? <div className="flex flex-col gap-4 justify-center items-center">
        <Button onClick={() => setWorkflow("bind")}><Icon icon="Bind" />Bind</Button>
        <Button onClick={() => setWorkflow("split")}><Icon icon="Split" />Split</Button>
      </div> : <div>
        <Button onClick={() => setWorkflow(null)}><Icon icon="Back" />Back</Button>
      </div>}
    </div>
    <div className="flex flex-col gap-4 items-center justify-center">
      {workflow === "bind" && <Bind />}
      {workflow === "split" && <Split />}
    </div>
  </div>;
}