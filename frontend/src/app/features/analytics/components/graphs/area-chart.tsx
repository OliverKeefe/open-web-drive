import * as React from "react"
import { Area, AreaChart, CartesianGrid, XAxis } from "recharts"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    ChartContainer,
    ChartLegend,
    ChartLegendContent,
    ChartTooltip,
    ChartTooltipContent,
    type ChartConfig,
} from "@/components/ui/chart"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
export const description = "An interactive area chart"
const chartData = [
    { date: "2024-04-01", download: 222, upload: 150 },
    { date: "2024-04-02", download: 97, upload: 180 },
    { date: "2024-04-03", download: 167, upload: 120 },
    { date: "2024-04-04", download: 242, upload: 260 },
    { date: "2024-04-05", download: 373, upload: 290 },
    { date: "2024-04-06", download: 301, upload: 340 },
    { date: "2024-04-07", download: 245, upload: 180 },
    { date: "2024-04-08", download: 409, upload: 320 },
    { date: "2024-04-09", download: 59, upload: 110 },
    { date: "2024-04-10", download: 261, upload: 190 },
    { date: "2024-04-11", download: 327, upload: 350 },
    { date: "2024-04-12", download: 292, upload: 210 },
    { date: "2024-04-13", download: 342, upload: 380 },
    { date: "2024-04-14", download: 137, upload: 220 },
    { date: "2024-04-15", download: 120, upload: 170 },
    { date: "2024-04-16", download: 138, upload: 190 },
    { date: "2024-04-17", download: 446, upload: 360 },
    { date: "2024-04-18", download: 364, upload: 410 },
    { date: "2024-04-19", download: 243, upload: 180 },
    { date: "2024-04-20", download: 89, upload: 150 },
    { date: "2024-04-21", download: 137, upload: 200 },
    { date: "2024-04-22", download: 224, upload: 170 },
    { date: "2024-04-23", download: 138, upload: 230 },
    { date: "2024-04-24", download: 387, upload: 290 },
    { date: "2024-04-25", download: 215, upload: 250 },
    { date: "2024-04-26", download: 75, upload: 130 },
    { date: "2024-04-27", download: 383, upload: 420 },
    { date: "2024-04-28", download: 122, upload: 180 },
    { date: "2024-04-29", download: 315, upload: 240 },
    { date: "2024-04-30", download: 454, upload: 380 },
    { date: "2024-05-01", download: 165, upload: 220 },
    { date: "2024-05-02", download: 293, upload: 310 },
    { date: "2024-05-03", download: 247, upload: 190 },
    { date: "2024-05-04", download: 385, upload: 420 },
    { date: "2024-05-05", download: 481, upload: 390 },
    { date: "2024-05-06", download: 498, upload: 520 },
    { date: "2024-05-07", download: 388, upload: 300 },
    { date: "2024-05-08", download: 149, upload: 210 },
    { date: "2024-05-09", download: 227, upload: 180 },
    { date: "2024-05-10", download: 293, upload: 330 },
    { date: "2024-05-11", download: 335, upload: 270 },
    { date: "2024-05-12", download: 197, upload: 240 },
    { date: "2024-05-13", download: 197, upload: 160 },
    { date: "2024-05-14", download: 448, upload: 490 },
    { date: "2024-05-15", download: 473, upload: 380 },
    { date: "2024-05-16", download: 338, upload: 400 },
    { date: "2024-05-17", download: 499, upload: 420 },
    { date: "2024-05-18", download: 315, upload: 350 },
    { date: "2024-05-19", download: 235, upload: 180 },
    { date: "2024-05-20", download: 177, upload: 230 },
    { date: "2024-05-21", download: 82, upload: 140 },
    { date: "2024-05-22", download: 81, upload: 120 },
    { date: "2024-05-23", download: 252, upload: 290 },
    { date: "2024-05-24", download: 294, upload: 220 },
    { date: "2024-05-25", download: 201, upload: 250 },
    { date: "2024-05-26", download: 213, upload: 170 },
    { date: "2024-05-27", download: 420, upload: 460 },
    { date: "2024-05-28", download: 233, upload: 190 },
    { date: "2024-05-29", download: 78, upload: 130 },
    { date: "2024-05-30", download: 340, upload: 280 },
    { date: "2024-05-31", download: 178, upload: 230 },
    { date: "2024-06-01", download: 178, upload: 200 },
    { date: "2024-06-02", download: 470, upload: 410 },
    { date: "2024-06-03", download: 103, upload: 160 },
    { date: "2024-06-04", download: 439, upload: 380 },
    { date: "2024-06-05", download: 88, upload: 140 },
    { date: "2024-06-06", download: 294, upload: 250 },
    { date: "2024-06-07", download: 323, upload: 370 },
    { date: "2024-06-08", download: 385, upload: 320 },
    { date: "2024-06-09", download: 438, upload: 480 },
    { date: "2024-06-10", download: 155, upload: 200 },
    { date: "2024-06-11", download: 92, upload: 150 },
    { date: "2024-06-12", download: 492, upload: 420 },
    { date: "2024-06-13", download: 81, upload: 130 },
    { date: "2024-06-14", download: 426, upload: 380 },
    { date: "2024-06-15", download: 307, upload: 350 },
    { date: "2024-06-16", download: 371, upload: 310 },
    { date: "2024-06-17", download: 475, upload: 520 },
    { date: "2024-06-18", download: 107, upload: 170 },
    { date: "2024-06-19", download: 341, upload: 290 },
    { date: "2024-06-20", download: 408, upload: 450 },
    { date: "2024-06-21", download: 169, upload: 210 },
    { date: "2024-06-22", download: 317, upload: 270 },
    { date: "2024-06-23", download: 480, upload: 530 },
    { date: "2024-06-24", download: 132, upload: 180 },
    { date: "2024-06-25", download: 141, upload: 190 },
    { date: "2024-06-26", download: 434, upload: 380 },
    { date: "2024-06-27", download: 448, upload: 490 },
    { date: "2024-06-28", download: 149, upload: 200 },
    { date: "2024-06-29", download: 103, upload: 160 },
    { date: "2024-06-30", download: 446, upload: 400 },
]
const chartConfig = {
    download: {
        label: "Downloads",
        color: "var(--chart-1)",
    },
    upload: {
        label: "Uploads",
        color: "var(--chart-2)",
    },
} satisfies ChartConfig
export function ChartAreaInteractive() {
    const [timeRange, setTimeRange] = React.useState("90d")
    const filteredData = chartData.filter((item) => {
        const date = new Date(item.date)
        const referenceDate = new Date("2024-06-30")
        let daysToSubtract = 90
        if (timeRange === "30d") {
            daysToSubtract = 30
        } else if (timeRange === "7d") {
            daysToSubtract = 7
        }
        const startDate = new Date(referenceDate)
        startDate.setDate(startDate.getDate() - daysToSubtract)
        return date >= startDate
    })
    return (
        <Card className="pt-0">
            <CardHeader className="flex items-center gap-2 space-y-0 border-b py-5 sm:flex-row">
                <div className="grid flex-1 gap-1">
                    <CardTitle>Uploads and Downloads</CardTitle>
                    <CardDescription>
                        Total uploads and downloads by members of your organisation.
                    </CardDescription>
                </div>
                <Select value={timeRange} onValueChange={setTimeRange}>
                    <SelectTrigger
                        className="hidden w-[160px] rounded-lg sm:ml-auto sm:flex"
                        aria-label="Select a value"
                    >
                        <SelectValue placeholder="Last 3 months" />
                    </SelectTrigger>
                    <SelectContent className="rounded-xl">
                        <SelectItem value="90d" className="rounded-lg">
                            Last 3 months
                        </SelectItem>
                        <SelectItem value="30d" className="rounded-lg">
                            Last 30 days
                        </SelectItem>
                        <SelectItem value="7d" className="rounded-lg">
                            Last 7 days
                        </SelectItem>
                    </SelectContent>
                </Select>
            </CardHeader>
            <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
                <ChartContainer
                    config={chartConfig}
                    className="aspect-auto h-[250px] w-full"
                >
                    <AreaChart data={filteredData}>
                        <defs>
                            <linearGradient id="fillDownload" x1="0" y1="0" x2="0" y2="1">
                                <stop
                                    offset="5%"
                                    stopColor="var(--color-download)"
                                    stopOpacity={0.8}
                                />
                                <stop
                                    offset="95%"
                                    stopColor="var(--color-download)"
                                    stopOpacity={0.1}
                                />
                            </linearGradient>
                            <linearGradient id="fillUpload" x1="0" y1="0" x2="0" y2="1">
                                <stop
                                    offset="5%"
                                    stopColor="var(--color-upload)"
                                    stopOpacity={0.8}
                                />
                                <stop
                                    offset="95%"
                                    stopColor="var(--color-upload)"
                                    stopOpacity={0.1}
                                />
                            </linearGradient>
                        </defs>
                        <CartesianGrid vertical={false} />
                        <XAxis
                            dataKey="date"
                            tickLine={false}
                            axisLine={false}
                            tickMargin={8}
                            minTickGap={32}
                            tickFormatter={(value) => {
                                const date = new Date(value)
                                return date.toLocaleDateString("en-US", {
                                    month: "short",
                                    day: "numeric",
                                })
                            }}
                        />
                        <ChartTooltip
                            cursor={false}
                            content={
                                <ChartTooltipContent
                                    labelFormatter={(value) => {
                                        return new Date(value).toLocaleDateString("en-US", {
                                            month: "short",
                                            day: "numeric",
                                        })
                                    }}
                                    indicator="dot"
                                />
                            }
                        />
                        <Area
                            dataKey="upload"
                            type="natural"
                            fill="url(#fillUpload)"
                            stroke="var(--color-upload)"
                            stackId="a"
                        />
                        <Area
                            dataKey="download"
                            type="natural"
                            fill="url(#fillDownload)"
                            stroke="var(--color-download)"
                            stackId="a"
                        />
                        <ChartLegend content={<ChartLegendContent />} />
                    </AreaChart>
                </ChartContainer>
            </CardContent>
        </Card>
    )
}