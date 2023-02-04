export interface Event {
    id: string;
    name: string;
    location: string;
    startTime: string;
    endTime: string;
    skillLevel: number;
    price: number;
    isAvailable: boolean;
    spotsLeft: number;
    url: string;
    sourceId: number;
    updatedOn: string;
}
