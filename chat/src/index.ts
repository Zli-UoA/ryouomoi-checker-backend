import * as config from './config'
import { UserJWTService } from './jwt';
import { ChatRepository, UserRepository } from './repository';

import { Server, Socket } from "socket.io";

async function main() {
    await new Promise(s => setTimeout(s, 3000));
    const io = new Server(3001, {
        cors: {
            origin: '*',
        }
    });
    const userRepository = new UserRepository(config.LocalEnv.db);
    const chatRepository = new ChatRepository(config.LocalEnv.db);
    const userJWTService = new UserJWTService(config.LocalEnv.jwt);

    io.on("connection", async (socket: Socket) => {
        try {
            const token = socket.handshake.auth.token as string;
            console.log(token);
            const userId = userJWTService.getUserIdFromJWT(token);
            console.log(userId);
            const user = await userRepository.getUser(userId);
            const chatRoom = await chatRepository.getChatRoom(userId);
            if (chatRoom === null || user === null) {
                socket.disconnect();
                return;
            }
            console.log(chatRoom);
            socket.join(chatRoom.id.toString());
            const chatHistory = await chatRepository.getAllChatHistory(chatRoom.id);
            socket.emit("chatHistory", chatHistory);

            socket.on("chat", async (message: string) => {
                console.log(message);
                const chat = {
                    id: 0,
                    room: chatRoom,
                    user: user,
                    message: message,
                    createdAt: new Date(),
                };
                const chatId = await chatRepository.newMessage(chat);
                chat.id = chatId;
                console.log(chat);
                socket.emit("chat", chat);
                socket.to(chatRoom.id.toString()).emit("chat", chat);
            });
        } catch (e) {
            console.log(e);
        }
    });

}

main().catch(err => {
    console.log(err);
})
