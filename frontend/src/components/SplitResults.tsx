import { motion } from "framer-motion";

import { Button } from "./ui/button";
import { SplitResultsProps } from "../types/core";
import { Icon } from "./Icon";
import { splitResultVariants } from "../lib/motions";

export default function SplitResults({ results, onBack, onDownload }: SplitResultsProps) {
  if (!results.data || results.error) return null;

  return (
    <motion.div
      variants={splitResultVariants.container}
      initial="hidden"
      animate="visible"
      className="w-full max-w-4xl mx-auto"
    >
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-semibold text-slate-100 mb-1">Shards</h3>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm" onClick={onBack}>
            <Icon icon="Back" className="w-4 h-4" />&nbsp;Back
          </Button>
        </div>
      </div>
      <div className="my-4 flex justify-start items-center gap-2">
        <Button disabled={!results.data} size="sm" onClick={() => onDownload(results.data!)}>
          <Icon icon="Download" className="w-4 h-4" />&nbsp;Save As...
        </Button>
        <p className="text-sm text-crystal-200">This will open a save dialog to choose where to save the shards as a txt file.</p>
      </div>
      <hr className="my-4 border-crystal-500/20" />
      <div className="grid grid-cols-1 gap-3 md:grid-cols-2 max-h-[200px] overflow-y-auto">
        {results.data?.split('\n')
          .filter(line => line.trim() !== '')
          .map((line, index) => (
            <motion.div
              key={index}
              variants={splitResultVariants.cardVariants}
              initial="hidden"
              animate="visible"
              transition={{ delay: index * 0.1 }}
              className="relative bg-crystal-700/50 rounded-md border border-crystal-600/50 p-2 hover:bg-crystal-700/70 hover:border-crystal-500/50 transition-all duration-200"
            >
              <div className="flex items-center justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <pre className="text-sm text-crystal-200 font-mono leading-relaxed select-none truncate" title={line}>
                    {line}
                  </pre>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => navigator.clipboard.writeText(line)}
                  className="duration-200 h-8 px-2 text-crystal-200 bg-crystal-600/50"
                >
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                    />
                  </svg>
                </Button>
              </div>
            </motion.div>
          ))}
      </div>
    </motion.div >
  );
}