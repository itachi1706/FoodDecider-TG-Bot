import {cookies} from "next/headers";
import {logoutUser} from "@/utils/session";
import {TelegramCookieData} from "@/types/tgcookiedata";

export async function GET() {
    const cookieStore = await cookies();
    if (cookieStore.has("auth")) {
        const cookieData = cookieStore.get("auth");
        const authData = JSON.parse(Buffer.from(cookieData!.value, 'base64').toString('utf-8')) as TelegramCookieData;
        const uuid = authData.uuid;
        await logoutUser(uuid);

        // TODO: Delete from DB
        console.log("Deleting cookie");
        cookieStore.delete("auth");
    }


    return Response.json({message: "Logged out"});
}