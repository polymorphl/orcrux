import { motion } from "framer-motion";

import { Button } from "./ui/button";

type SplitResultsProps = {
  results: string;
}

const resultVariants = {
  hidden: { opacity: 0, y: 20, scale: 0.95 },
  visible: {
    opacity: 1,
    y: 0,
    scale: 1,
    transition: { duration: 0.4, ease: "easeOut" as const }
  }
};

const cardVariants = {
  hidden: { opacity: 0, y: 15 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.3, ease: "easeOut" as const }
  }
};

export default function SplitResults({ results }: SplitResultsProps) {
  if (!results) return null;

  return (
    <motion.div
      variants={resultVariants}
      initial="hidden"
      animate="visible"
      className="w-full max-w-2xl mx-auto"
    >
      <div className="bg-gradient-to-br from-slate-800/80 to-slate-700/60 backdrop-blur-sm rounded-lg border border-slate-600/50 shadow-xl p-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h3 className="text-lg font-semibold text-slate-100 mb-1">Split Results</h3>
          </div>
        </div>
        <div className="grid grid-cols-1 gap-3">
          {results.split('\n')
            .filter(line => line.trim() !== '')
            .map((line, index) => (
              <motion.div
                key={index}
                variants={cardVariants}
                initial="hidden"
                animate="visible"
                transition={{ delay: index * 0.1 }}
                className="relative bg-crystal-700/50 rounded-md border border-crystal-600/50 p-4 hover:bg-crystal-700/70 hover:border-crystal-500/50 transition-all duration-200"
              >
                <div className="flex items-center justify-between gap-3">
                  <div className="flex-1 min-w-0">
                    <pre className="text-sm text-crystal-200 font-mono break-all leading-relaxed">
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
      </div>
    </motion.div>
  );
}