import {AuthDataValidator, objectToAuthDataMap} from "@telegram-auth/server";
import {cookies} from "next/headers";

const validator = new AuthDataValidator({ botToken: process.env.BOT_TOKEN });

export async function POST(req: Request) {
    console.log("IM CALLED");

    const rawData = await req.json();
    const dataMap = objectToAuthDataMap(rawData);
    const cookieStore = await cookies();

    try {
        const user = await validator.validate(dataMap);

        const base64Data = Buffer.from(JSON.stringify(user)).toString("base64");

        if (!cookieStore.has("auth")) {
            cookieStore.set("auth", base64Data, {
                httpOnly: true,
                secure: true,
                sameSite: "strict",
            });
        }
        return Response.json(user);
    } catch (error) {
        console.error("FAILED")
        if (cookieStore.has("auth")) {
            cookieStore.delete("auth");
        }
        return Response.json({ error: error }, { status: 400 });
    }
}