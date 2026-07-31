export type MeterGaugeSegment = {
    label: string;
    value: number;
    color: string;
    percentage?: number;
}

type MeterGaugeProps = React.PropsWithChildren<{
    segmentData: MeterGaugeSegment[];
    total: number;
    children: React.ReactNode | undefined;
}>;

export function MeterGauge({ segmentData, total, children }: MeterGaugeProps) {
    const processedSegments = segmentData.map(seg => ({
        ...seg,
        percentage: (seg.value / total) * 100,
    }))
    return (
        <div className={"relative w-full h-[8px] rounded-full bg-neutral-600 overflow-hidden"}>
            <div className={"absolute inset-0 flex"}>
                {processedSegments.map((seg, i) => (
                    <Segment key={i} {...seg} />
                ))}
            </div>
        </div>
    );
}

function Segment({ label, color, percentage }: MeterGaugeSegment) {
    return (
        <div
            className={`mr-[2px] h-full ${color}`}
            style={{ width: `${percentage}%` }}
            title={label}
        />
    );
}

type BackgroundProps = React.PropsWithChildren<{
    children: React.ReactNode | undefined;
}>;

function Background({ children }: BackgroundProps) {
    return (
        <div className={"rounded-b-full h[5px] w-full"}>{children}</div>
    );
}

