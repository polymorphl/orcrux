import { Icon } from "@/components/Icon";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { motion, AnimatePresence } from "framer-motion";
import { useState, useEffect } from "react";

type TabSharpProps = {
  initialTab: string
  tabs: { name: "Bind" | "Split", value: string, content: React.ReactNode }[]
}

export default function TabsSharp({ tabs, initialTab }: TabSharpProps) {
  const [activeTab, setActiveTab] = useState(initialTab);
  const [indicatorPosition, setIndicatorPosition] = useState(0);
  const [indicatorWidth, setIndicatorWidth] = useState(0);

  useEffect(() => {
    const activeTabElement = document.querySelector(`[data-value="${activeTab}"]`) as HTMLElement;
    if (activeTabElement) {
      const rect = activeTabElement.getBoundingClientRect();
      const parentRect = activeTabElement.parentElement?.getBoundingClientRect();
      if (parentRect) {
        setIndicatorPosition(rect.left - parentRect.left);
        setIndicatorWidth(rect.width);
      }
    }
  }, [activeTab]);

  return (
    <Tabs defaultValue={initialTab} onValueChange={setActiveTab} className="max-w-xs w-full">
      <TabsList className="w-full p-0 justify-start border-b border-crystal-500/50 rounded-none bg-crystal-700/70 backdrop-blur-xl relative">
        {tabs.map((tab) => (
          <TabsTrigger
            key={tab.value}
            value={tab.value}
            data-value={tab.value}
            className="rounded-none h-full data-[state=active]:shadow-none border border-b-[3px] border-transparent data-[state=active]:border-crystal-300 bg-transparent text-crystal-200 hover:text-crystal-100 hover:bg-crystal-600/50 data-[state=active]:bg-crystal-600/70 data-[state=active]:text-crystal-100 transition-colors cursor-pointer relative"
          >
            <Icon icon={tab.name} />
            <code className="text-[13px]">{tab.name}</code>
          </TabsTrigger>
        ))}

        {/* Sliding indicator */}
        <motion.div
          className="absolute bottom-0 h-0.5 bg-crystal-300"
          initial={false}
          animate={{
            x: indicatorPosition,
            width: indicatorWidth,
          }}
          transition={{ duration: 0.3, ease: "easeInOut" }}
        />
      </TabsList>

      <AnimatePresence mode="wait">
        {tabs.find(tab => tab.value === activeTab) && (
          <TabsContent key={activeTab} value={activeTab} className="mt-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.15 }}
              className="w-full"
            >
              {tabs.find(tab => tab.value === activeTab)?.content}
            </motion.div>
          </TabsContent>
        )}
      </AnimatePresence>
    </Tabs>
  );
}
