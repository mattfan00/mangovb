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

export interface Filter {
    source: FilterEntry[];
    skillLevel: FilterEntry[];
    spots: FilterEntry[];
}

export interface FilterEntry {
    value: number;
    text: string;
}
