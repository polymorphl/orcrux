import Bind from "./Bind";
import Split from "./Split";
import TabsSharp from "./customized/tabs/tabs-10";

const tabs = [
  { name: "Split" as const, value: "split", content: <Split /> },
  { name: "Bind" as const, value: "bind", content: <Bind /> },
]

export default function Wizard() {
  return <TabsSharp tabs={tabs} initialTab={"split"} />
}