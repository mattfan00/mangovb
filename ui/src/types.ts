export interface Event {
    id: string;
    name: string;
    location: string;
    startTime: Date;
    endTime: Date;
    skillLevel: number;
    price: number;
    isAvailable: boolean;
    spotsLeft: number;
    url: string;
    sourceId: number;
    updatedOn: Date;
}
