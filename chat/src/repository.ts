import * as mysql from 'mysql2/promise';
import { DBConfig } from './config';
import { Chat, ChatRoom, User } from './model';
import { ResultSetHeader, RowDataPacket } from 'mysql2';

export class UserRepository {
    pool: mysql.Pool;

    constructor(dbConfig: DBConfig) {
        this.pool = mysql.createPool(dbConfig);
    }

    async getUser(userId: string): Promise<User | null> {
        try {
            const sql = `
                SELECT * from twitter_users
                WHERE twitter_id = ?
            `;
            const params = [userId];
            const [rows] = await this.pool.execute(sql, params);
            if (Array.isArray(rows) && rows.length > 0) {
                const row = rows[0] as RowDataPacket
                return {
                    id: row.twitter_id as number,
                    screenName: row.screen_name,
                    displayName: row.display_name,
                    profileImageUrl: row.profile_image_url,
                    biography: row.biography,
                };
            } else {
                return null;
            }
        } catch (e) {
            console.log(e);
            return null;
        }
    }
}

export class ChatRepository {
    pool: mysql.Pool;

    constructor(dbConfig: DBConfig) {
        this.pool = mysql.createPool(dbConfig)
    }

    async getChatRoom(userId: string): Promise<ChatRoom | null> {
        try {
            const sql = `
            SELECT * FROM chat_rooms cr
            JOIN couples c ON cr.couple_id = c.id
            WHERE (c.user_id_1 = ? OR c.user_id_2 = ?) AND c.broken_at IS NULL
        `;
            const params = [userId, userId];
            const [rows] = await this.pool.query(sql, params) as any;
            if (Array.isArray(rows) && rows.length > 0) {
                const row = rows[0]
                return {
                    id: row.id as number,
                    coupleId: row.couple_id,
                    createdAt: row.created_at,
                };
            } else {
                return null;
            }
        } catch (e) {
            console.log(e);
            return null;
        }
    }

    async getAllChatHistory(chatRoomId: number): Promise<Chat[]> {
        try {
            const sql = `
            SELECT
                c.id cid,
                c.chat_room_id,
                c.user_id,
                c.message,
                c.created_at,
                u.screen_name,
                u.display_name,
                u.profile_image_url,
                u.biography,
                cr.id,
                cr.couple_id,
                cr.created_at
            FROM chats c
            JOIN twitter_users u ON c.user_id = u.twitter_id
            JOIN chat_rooms cr ON c.chat_room_id = cr.id
            WHERE cr.id = ?
            ORDER BY c.created_at ASC
        `;
            const params = [chatRoomId];
            const [rows] = await this.pool.query(sql, params);
            if (Array.isArray(rows)) {
                return rows.map((r: any) => ({
                    id: r.cid,
                    room: {
                        id: r.chat_room_id,
                        coupleId: r.couple_id,
                        createdAt: r.created_at,
                    },
                    user: {
                        id: r.user_id,
                        screenName: r.screen_name,
                        displayName: r.display_name,
                        profileImageUrl: r.profile_image_url,
                        biography: r.biography,
                    },
                    message: r.message,
                    createdAt: r.created_at,
                }));
            } else {
                return [];
            }
        } catch (e) {
            console.log(e);
            return [];
        }
    }

    async newMessage(chat: Chat): Promise<number> {
        try {
            const sql = `
            INSERT INTO chats (chat_room_id, user_id, message, created_at)
            VALUES (?, ?, ?, ?)
        `;
            const params = [chat.room.id, chat.user.id, chat.message, chat.createdAt];
            const [result] = await this.pool.execute(sql, params);
            const header = result as ResultSetHeader
            return header.insertId;
        } catch (e) {
            console.log(e)
            return 0;
        }
    }
}
