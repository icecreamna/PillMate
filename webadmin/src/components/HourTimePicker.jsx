// 24HourTimePicker.jsx
import { useMemo } from 'react';

export default function  TwentyFourHourTimePicker({
  value = '09:00',               // "HH:mm"
  onChange,                      // (str) => void
  minuteStep = 5,                // 1/5/10/15...
  disabled = false,
  className = '',
}) {
  const [hh, mm] = (value || '00:00').split(':');

  const hours = useMemo(
    () => Array.from({ length: 24 }, (_, i) => String(i).padStart(2, '0')),
    []
  );
  const minutes = useMemo(() => {
    const arr = [];
    for (let m = 0; m < 60; m += minuteStep) arr.push(String(m).padStart(2, '0'));
    return arr;
  }, [minuteStep]);

  const update = (newH, newM) => {
    const next = `${newH ?? hh}:${newM ?? mm}`;
    onChange?.(next);
  };

  return (
    <div className={className} style={{ display: 'inline-flex', gap: 8 }}>
      <select
        value={hh}
        disabled={disabled}
        onChange={(e) => update(e.target.value, undefined)}
        aria-label="ชั่วโมง (00-23)"
      >
        {hours.map(h => <option key={h} value={h}>{h}</option>)}
      </select>
      :
      <select
        value={mm}
        disabled={disabled}
        onChange={(e) => update(undefined, e.target.value)}
        aria-label="นาที"
      >
        {minutes.map(m => <option key={m} value={m}>{m}</option>)}
      </select>
    </div>
  );
}
