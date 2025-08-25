import { Label } from "./ui/label";

type ShardsSliderProps = {
  label: string;
  value: number;
  min: number;
  max: number;
  onChange: (value: number) => void;
}

export default function ShardsSlider({ label, value, min, max, onChange }: ShardsSliderProps) {
  return (
    <div>
      <Label htmlFor="shards">{label}: {value}</Label>
      <input
        type="range"
        min={min}
        max={max}
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        className="w-full h-1 rounded-lg appearance-none cursor-pointer slider"
        style={{
          backgroundColor: '#e5e7eb',
        }}
      />
      <div className="flex justify-between text-xs text-crystal-500 mt-1">
        <span>{min}</span>
        <span>{max}</span>
      </div>
    </div>
  )
}