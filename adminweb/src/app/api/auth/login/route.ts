import {AuthDataValidator, objectToAuthDataMap} from "@telegram-auth/server";
import {cookies} from "next/headers";
import {checkIsAdmin} from "@/utils/users";
import {TelegramCookieData} from "@/types/tgcookiedata";
import {createLoginRecord} from "@/utils/session";

const validator = new AuthDataValidator({ botToken: process.env.BOT_TOKEN });

export async function POST(req: Request) {
    const rawData = await req.json();
    const dataMap = objectToAuthDataMap(rawData);
    const cookieStore = await cookies();

    try {
        const user = await validator.validate(dataMap);

        // Check if user is admin (fail otherwise)
        const adminChk = await checkIsAdmin(user.id);
        if (!await checkIsAdmin(user.id)) {
            return Response.json({ error: "User is not an admin" }, { status: 403 });
        }

        const uuid = await createLoginRecord(user, adminChk, dataMap);
        if (uuid == null) {
            return Response.json({ error: "Failed to create login record" }, { status: 500 });
        }

        // Save base64 data to DB
        const cookieData: TelegramCookieData = {
            id: user.id,
            is_admin: adminChk,
            uuid: uuid,
        }

        const base64CookieData = Buffer.from(JSON.stringify(cookieData)).toString("base64");

        cookieStore.set("auth", base64CookieData, {
            httpOnly: true,
            secure: true,
            sameSite: "strict",
        });
        return Response.json({ message: "Logged in" });
    } catch (error) {
        console.error("Failed Login")
        if (cookieStore.has("auth")) {
            cookieStore.delete("auth");
        }
        return Response.json({ error: error }, { status: 400 });
    }
}