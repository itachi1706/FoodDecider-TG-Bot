"use server";
import { v4 as uuidv4 } from 'uuid';
import {TelegramAuthenticationData} from "@/types/tgauthdata";
import db from "@/utils/database";
import {RowDataPacket} from "mysql2";
import {TelegramUserData} from "@telegram-auth/server";

export async function createLoginRecord(user: TelegramUserData, adminChk: boolean, dataMap: Map<string, string | number>) : Promise<null|string> {

    const authDataStore: TelegramAuthenticationData = {
        id: user.id,
        first_name: user.first_name,
        last_name: user.last_name,
        photo_url: user.photo_url,
        username: user.username,
        is_bot: user.is_bot,
        language_code: user.language_code,
        is_premium: user.is_premium,
        is_admin: adminChk,
        data_map: JSON.stringify(Array.from(dataMap.entries())),
    }

    const base64Data = Buffer.from(JSON.stringify(authDataStore)).toString("base64");
    console.log(base64Data);
    const uuid: string = uuidv4().toString();

    const sql = "INSERT INTO web_auth_session (id, telegram_id, data, created_at, status) VALUES (?, ?, ?, NOW(), \"A\")";
    const [rows] = await db.execute<RowDataPacket[]>(sql, [uuid, user.id, base64Data]);

    if (rows.length == 0) {
        console.error("Failed to create login record");
        return null;
    }

    return uuid;
}

export async function logoutUser(uuid: string) : Promise<null|string> {

    const sql = "UPDATE web_auth_session SET status = \"D\" WHERE id = ?";
    const [rows] = await db.execute<RowDataPacket[]>(sql, [uuid]);

    if (rows.length == 0) {
        console.error("Failed to delete login record");
        return null;
    }

    return uuid;
}

export async function getUserSession(uuid: string): Promise<null|TelegramAuthenticationData> {
    console.log("Getting User Session for:", uuid);
    const sql = "SELECT * FROM web_auth_session WHERE id = ? AND status = \"A\"";
    const [rows] = await db.execute<RowDataPacket[]>(sql, [uuid]);

    if (rows.length == 0) {
        console.error("Failed to get user data");
        return null;
    }

    const data = rows[0].data as string;
    return JSON.parse(Buffer.from(data, "base64").toString("utf-8")) as TelegramAuthenticationData;
}