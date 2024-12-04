import {AuthDataValidator, objectToAuthDataMap} from "@telegram-auth/server";

const validator = new AuthDataValidator({ botToken: process.env.BOT_TOKEN });

export async function POST(req: Request) {
    console.log("IM CALLED");

    const rawData = await req.json();
    const dataMap = objectToAuthDataMap(rawData);

    console.log("HIHI", dataMap.get("auth_date")) ;
    try {
        const user = await validator.validate(dataMap);
        return Response.json(user);
    } catch (error) {
        console.error(error);
        return Response.json({ error: error }, { status: 400 });
    }
}