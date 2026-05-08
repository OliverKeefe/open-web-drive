import {Container} from "@/app/features/shared/components/layout/container.tsx";
import {ActivityInfoCard} from "@/app/features/git/components/activity-info-card.tsx";
import {ClockAlert} from "lucide-react";


export default function GitHome() {
    return (
        <Container>
            <div className="grid grid-cols-2 gap-5">
                <div className="grid grid-cols-4 gap-4">
                    <ActivityInfoCard title={"Issues"} total={9} time={"12 mins ago."} description={"Assigned to You"} icon={<ClockAlert />} />
                    <ActivityInfoCard title={"Pull Requests"} total={3} time={"53 mins ago."} description={"Assigned to You"} icon={<ClockAlert />} />
                    <ActivityInfoCard title={"Merges Requests"} total={24} time={"2hrs ago."} description={"Assigned to You"} icon={<ClockAlert />} />
                    <ActivityInfoCard title={"Tasks"} total={0} time={"26 secs ago."} description={"Assigned to You"} icon={<ClockAlert />} />
                </div>
            </div>
        </Container>
    );
}