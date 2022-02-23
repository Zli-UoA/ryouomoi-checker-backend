export type User = {
    id: number;
    screenName: string;
    displayName: string;
    profileImageUrl: string;
    biography: string;
};

export type ChatRoom = {
    id: number;
    coupleId: number;
    createdAt: Date;
};

export type Chat = {
    id: number;
    room: ChatRoom;
    user: User;
    message: string;
    createdAt: Date;
};
