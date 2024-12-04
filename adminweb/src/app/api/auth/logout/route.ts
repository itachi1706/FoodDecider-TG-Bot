import {cookies} from "next/headers";

export async function GET() {
    const cookieStore = await cookies();
    if (cookieStore.has("auth")) {
        console.log("Deleting cookie");
        cookieStore.delete("auth");
    }
    return Response.json({message: "Logged out"});
}