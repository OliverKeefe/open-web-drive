import {Card, CardContent, CardDescription, CardTitle} from "@/components/ui/card.tsx";

interface ActivityInfoCardProps {
    title: string;
    total: number;
    description: string;
    time: string;
    icon: React.JSX.Element;
}

export function ActivityInfoCard({title, total, time, icon, description}: ActivityInfoCardProps) {
    return (
        <Card className={"max-w-[190px] drop-shadow-xl"}>
            <CardContent>
                <p className="font-black text-7xl text-center">{total}</p>
                <div className={"flex items-center"}>
                    <div className={"pr-1"}>
                        {icon}
                    </div>
                    <CardTitle title={title}>{title}</CardTitle>
                </div>
                <CardDescription>{description}</CardDescription>

                <p>{time}</p>
            </CardContent>
        </Card>
    );
}