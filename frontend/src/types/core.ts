export type SplitResult = { error: string | null, data: string | null }
export type SplitFormProps = {
  onSplit: (secret: string, shards: number, shardsNeeded: number, output: string) => void;
}

export type RecomposeResult = { error: string | null, data: string | null }
export type SplitResultsProps = {
  results: {
    error: string | null;
    data: string | null;
  };
}