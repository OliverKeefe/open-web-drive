import {Container} from "@/app/features/shared/components/layout/container.tsx";
import {ChartAreaInteractive} from "@/app/features/analytics/components/graphs/area-chart.tsx";
import {ChartBarMixed} from "@/app/features/analytics/components/graphs/bar-mixed.tsx";


export function Analytics() {
    return (
        <Container>
            <h1 className="text-2xl font-semibold pb-4 pt-4 m-1">Your Analytics</h1>
            <ChartBarMixed />
            <ChartAreaInteractive />
        </Container>
    );
}

export default Analytics