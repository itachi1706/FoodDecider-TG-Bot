import {AuthDataMap, AuthDataValidator} from "@telegram-auth/server";
import {cookies} from "next/headers";
import {TelegramCookieData} from "@/types/tgcookiedata";
import {getUserSession} from "@/utils/session";

const validator = new AuthDataValidator({ botToken: process.env.BOT_TOKEN });

export async function POST() {
    const cookieStore = await cookies();

    // Match with cookie data
    const cookieData = cookieStore.get("auth");
    console.log("Cookie data", cookieData);
    if (cookieData == null) {
        return Response.json({ error: "No auth cookie" }, { status: 401 });
    }

    const authData = JSON.parse(Buffer.from(cookieData.value, 'base64').toString('utf-8')) as TelegramCookieData;
    console.log("Auth data", authData);
    const userObject = await getUserSession(authData.uuid);
    if (userObject == null) {
        return Response.json({ error: "Failed to get user" }, { status: 403 });
    }

    console.log(userObject);

    if (userObject.id != authData.id) {
        return Response.json({ error: "User ID mismatch" }, { status: 403 });
    }
    const dataMap = new Map(JSON.parse(userObject.data_map)) as AuthDataMap;

    try {
        const user = await validator.validate(dataMap);
        return Response.json(user);
    } catch (error) {
        console.error("FAILED")
        if (cookieStore.has("auth")) {
            cookieStore.delete("auth");
        }
        return Response.json({ error: error }, { status: 400 });
    }
}